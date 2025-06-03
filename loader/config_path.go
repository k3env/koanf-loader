package loader

import (
	"github.com/hashicorp/consul/api"
	"github.com/knadh/koanf/providers/consul"
)

type ConfigLoader struct {
	Files  []string      `koanf:"file"`
	Consul *ConsulConfig `koanf:"consul"`
}

type ConsulConfig struct {
	Addr       string `koanf:"addr"`
	Scheme     string `koanf:"scheme"`
	PathPrefix string `koanf:"prefix"`
	Token      string `koanf:"token"`
	Namespace  string `koanf:"namespace"`
	Key        string `koanf:"key"`
	Recursive  bool   `koanf:"recursive"`
}

func (c *ConsulConfig) Config() *api.Config {
	config := api.DefaultConfig()
	if c.Addr != "" {
		config.Address = c.Addr
	}
	if c.Scheme != "" {
		config.Scheme = c.Scheme
	}
	if c.Token != "" {
		config.Token = c.Token
	}
	if c.Namespace != "" {
		config.Namespace = c.Namespace
	}

	return config
}

func (c *ConsulConfig) ProviderConfig() consul.Config {
	cfg := c.Config()
	return consul.Config{
		Key:      c.Key,
		Recurse:  c.Recursive,
		Detailed: false,
		Cfg:      cfg,
	}
}
