// Package config holds the configuration struct.
package config

type Config struct {
	ServeAddress string
	DatabaseFile string
	InsecureCORS bool
}

// Default returns the default config.
func Default() *Config {
	conf := &Config{
		ServeAddress: "127.0.0.1:8118",
		DatabaseFile: "/path/to/database_file.bolt",
		InsecureCORS: false,
	}
	return conf
}
