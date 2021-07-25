// Package config holds the configuration struct.
package config

type Config struct {
	ServeAddress string `json:"serveAddress"`
	DatabaseFile string `json:"databaseFile"`
}

// Default returns the default config.
func Default() *Config {
	conf := &Config{
		ServeAddress: "127.0.0.1:8118",
		DatabaseFile: "/path/to/database_file.bolt",
	}
	return conf
}
