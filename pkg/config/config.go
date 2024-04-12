package config

// Config holds the configuration for our server.
type Config struct {
	Path string
	Port string
}

// LoadConfig loads configuration, overriding defaults with provided config values.
func LoadConfig(override Config) *Config {
	// Set default values and override them if provided
	cfg := Config{
		Path: ifEmpty(override.Path, "/"),    // Default path if override is empty
		Port: ifEmpty(override.Port, "8080"), // Default port if override is empty
	}

	return &cfg
}

// ifEmpty returns the override value if it is not empty, otherwise the default.
func ifEmpty(override, defaultVal string) string {
	if override != "" {
		return override
	}
	return defaultVal
}
