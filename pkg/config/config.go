// Package config contains functions for loading and validating the server configuration.
package config

import (
	"fmt"
	"net"
	"strconv"
)

// Config holds the configuration for our server.
type Config struct {
	// Path is the URL path where the server will listen for incoming connections.
	Path string
	// Port is the TCP port on which the server will listen for incoming connections.
	Port string
}

// LoadConfig loads configuration, overriding defaults with provided config values.
//
// The config parameter can be used to override default configuration values. It
// will be merged with the default configuration, with provided values taking
// precedence over the default values.
//
// If the Path or Port values are not provided, the default values will be used.
// If the Port value is not a valid TCP port, an error will be returned.
func LoadConfig(override Config) (*Config, error) {
	cfg := Config{
		Path: "/",    // Default path if override is empty
		Port: "8080", // Default port if override is empty
	}

	if override.Path != "" {
		cfg.Path = override.Path
	}

	if override.Port != "" {
		_, err := strconv.Atoi(override.Port)
		if err != nil {
			return nil, fmt.Errorf("invalid port '%s': %w", override.Port, err)
		}
		cfg.Port = override.Port
	}

	return &cfg, nil
}

// ifEmpty returns the override value if it is not empty, otherwise the default.
func ifEmpty(override, defaultVal string) string {
	if override != "" {
		return override
	}
	return defaultVal
}

// IsValidPort returns whether the given string is a valid TCP port number.
func IsValidPort(port string) bool {
	_, err := strconv.Atoi(port)
	return err == nil
}

// IsListening returns whether the server is listening on the given host and port.
func IsListening(host string, port string) bool {
	ln, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// IsConfigValid returns whether the given configuration is valid.
//
// This will check that the port is a valid TCP port number, and that the
// server is not already listening on the given host and port.
func IsConfigValid(cfg Config) bool {
	return IsValidPort(cfg.Port) && IsListening("", cfg.Port)
}

// String returns a string representation of the configuration.
func (c Config) String() string {
	return fmt.Sprintf("Config(Path=%s, Port=%s)", c.Path, c.Port)
}

