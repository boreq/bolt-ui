// Package config holds the configuration struct.
package config

type Config struct {
	ServeAddress  string `json:"serveAddress"`
	DataDirectory string `json:"dataDirectory"`
}

// Default returns the default config.
func Default() *Config {
	conf := &Config{
		ServeAddress:  "127.0.0.1:8118",
		DataDirectory: "/path/to/data/directory",
	}
	return conf
}
