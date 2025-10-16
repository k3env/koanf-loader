package loader

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/k3env/koanf-loader/providers/args"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/consul"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

func Load(out interface{}, envPrefix string) error {
	var err error
	var pc struct {
		Config ConfigLoader `koanf:"config"`
	}

	ap := args.Provider()

	pathLoader := koanf.New(".")
	err = pathLoader.Load(ap, nil)
	if err != nil {
		return err
	}
	err = pathLoader.Unmarshal("", &pc)
	if err != nil {
		return err
	}

	realLoader := koanf.New(".")
	for _, f := range pc.Config.Files {
		parser, err := selectParser(filepath.Ext(f))
		if err != nil {
			continue
		}
		err = realLoader.Load(file.Provider(f), parser)
	}
	if pc.Config.Consul != nil {
		err := realLoader.Load(consul.Provider(pc.Config.Consul.ProviderConfig()), nil)
		if err != nil {
			return err
		}
	}

	err = realLoader.Load(env.Provider("APP", "_", envparser(envPrefix)), nil)
	if err != nil {
		return err
	}
	err = realLoader.Load(ap, nil)
	if err != nil {
		return err
	}
	err = realLoader.Unmarshal("", &out)
	return err
}

func selectParser(ext string) (koanf.Parser, error) {
	switch ext {
	case ".json":
		return json.Parser(), nil
	case ".yaml":
		return yaml.Parser(), nil
	case ".yml":
		return yaml.Parser(), nil
	case ".toml":
		return toml.Parser(), nil
	case ".hcl":
		// recomendation: https://github.com/knadh/koanf?tab=readme-ov-file#api
		return hcl.Parser(true), nil
	default:
		return nil, fmt.Errorf("unknown parser extension: %s", ext)
	}
}

func envparser(prefix string) func(s string) string {
	return func(s string) string {
		realPrefix := fmt.Sprintf("%s_", strings.ToUpper(prefix))
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, realPrefix)), "_", ".", -1)
	}
}
