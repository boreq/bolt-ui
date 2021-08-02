// Package config holds the configuration struct.
package config

import "crypto/tls"

type Config struct {
	ServeAddress  string
	DatabaseFile  string
	Token         string
	Certificate   tls.Certificate
	InsecureCORS  bool
	InsecureToken bool
	InsecureTLS   bool
}
