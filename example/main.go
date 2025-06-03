package main

import (
	"log"
	"time"

	"github.com/k3env/koanf-loader/loader"
)

type Config struct {
	Name   string `koanf:"name"`
	Nested struct {
		Foo int           `koanf:"foo"`
		Bar bool          `koanf:"bar"`
		Baz time.Duration `koanf:"baz"`
	} `koanf:"nested"`
}

func main() {
	var cfg Config
	err := loader.Load(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)
}
