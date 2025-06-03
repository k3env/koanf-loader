package extensions

import (
	"fmt"
	"net/url"
)

type AMQPConfig struct {
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Vhost    string `koanf:"vhost"`
}

func (cfg AMQPConfig) Uri() string {
	vhost := cfg.Vhost
	if vhost == "" {
		vhost = "/"
	}
	port := cfg.Port
	if port == 0 {
		port = 5672
	}
	var user *url.Userinfo
	if (cfg.User != "") || (cfg.Password != "") {
		user = url.UserPassword(cfg.User, cfg.Password)
	}
	uri := url.URL{
		Scheme: "amqp",
		User:   user,
		Host:   fmt.Sprintf("%s:%d", vhost, port),
		Path:   cfg.Vhost,
	}
	return uri.String()
}
