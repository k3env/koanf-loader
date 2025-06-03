package extensions

import "net/url"

type PostgresConfig struct {
	Username string            `koanf:"username"`
	Password string            `koanf:"password"`
	Host     string            `koanf:"host"`
	Port     string            `koanf:"port"`
	Database string            `koanf:"database"`
	Options  map[string]string `koanf:"options"`
}

func (c PostgresConfig) Uri() string {
	vals := make(url.Values)
	for k, v := range c.Options {
		vals.Set(k, v)
	}

	port := c.Port
	if port == "" {
		port = "5432"
	}

	uri := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.Username, c.Password),
		Host:     c.Host + ":" + c.Port,
		Path:     c.Database,
		RawQuery: vals.Encode(),
	}
	return uri.String()
}
