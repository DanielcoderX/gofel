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
	// Verbose enables verbose logging.
	Verbose bool
	// Serialization format: "json" or "msgpack"
	Format string
	// ConnectionPoolCapacity is the maximum number of concurrent connections our server can handle.
	ConnectionPoolCapacity int
}

// GlobalConfig holds the application's configuration.
var GlobalConfig *Config


// LoadConfig loads configuration, overriding defaults with provided config values.
func LoadConfig(override Config) (*Config, error) {
	cfg := Config{
		Path:                   "/ws",    // Default path if override is empty
		Port:                   "8080", // Default port if override is empty
		Verbose:                true,  // Default verbose if override is empty
		Format:                 "json", // Default format if override is empty
		ConnectionPoolCapacity: 100,    // Default connection pool capacity if override is empty
	}

	if override.Path != "" {
		cfg.Path = override.Path
	}

	if override.Port != "" {
		if _, err := strconv.Atoi(override.Port); err != nil {
			return nil, fmt.Errorf("invalid port '%s': %w", override.Port, err)
		}

		if err := IsConfigValid(override); !err {
			return nil, fmt.Errorf("invalid configuration: %+v", cfg)
		}

		cfg.Port = override.Port
	}
	if override.Format != "" {
		if override.Format != "json" && override.Format != "msgpack" {
			return nil, fmt.Errorf("invalid format '%s': must be 'json' or 'msgpack'", override.Format)
		}
		cfg.Format = override.Format
	}
	cfg.ConnectionPoolCapacity = override.ConnectionPoolCapacity
	cfg.Verbose = override.Verbose
	return &cfg, nil
}

// IsPortAvailable checks if a given TCP port is available for listening.
func IsPortAvailable(port string) bool {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

// IsValidPort returns whether the given string is a valid TCP port number.
func IsValidPort(port string) bool {
	_, err := strconv.Atoi(port)
	return err == nil
}

// IsConfigValid returns whether the given configuration is valid.
//
// This will check that the port is a valid TCP port number, and that the
// server is not already listening on the given host and port.
func IsConfigValid(cfg Config) bool {
	return IsValidPort(cfg.Port) && IsPortAvailable(cfg.Port)
}

// String returns a string representation of the configuration.
func (c Config) String() string {
	return fmt.Sprintf("Config(Path=%s, Port=%s)", c.Path, c.Port)
}
