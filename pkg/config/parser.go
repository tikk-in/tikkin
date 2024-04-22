package config

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

func lookupEnvOrString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func lookupEnvOrInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (*ConfigFlags, error) {
	flags := ConfigFlags{}

	flag.StringVar(&flags.ConfigPath, "config", lookupEnvOrString("CONFIG_PATH", "./config.yml"), "config file path")
	flag.StringVar(&flags.AdminPassword, "admin-password", lookupEnvOrString("ADMIN_PASSWORD", ""), "Admin password")
	flag.StringVar(&flags.SMTPPassword, "smtp-password", lookupEnvOrString("SMTP_PASSWORD", ""), "SMTP password")
	flag.StringVar(&flags.JWTSecret, "jwt-secret", lookupEnvOrString("SERVER_JWT_SECRET", ""), "JWT secret")

	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(flags.ConfigPath); err != nil {
		return nil, err
	}

	// Return the configuration path
	return &flags, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func LoadConfig() (*Config, error) {
	flags, err := ParseFlags()
	if err != nil {
		return nil, err
	}
	return loadConfig(flags.ConfigPath, *flags)
}

func loadConfig(configPath string, configFlags ConfigFlags) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	ensureDefaultValues(config)

	if configFlags.SMTPPassword != "" {
		config.Email.SMTP.Password = configFlags.SMTPPassword
	}
	if configFlags.AdminPassword != "" {
		config.Admin.Password = configFlags.AdminPassword
	}
	if configFlags.JWTSecret != "" {
		config.Server.Jwt.Secret = configFlags.JWTSecret
	}

	log.Info().
		Msgf("Loaded config: server.port=%d, server.jwt.secret=<redacted>, db.host=%s,  db.port=%d, db.user=%s, db.password=<redacted>, db.database=%s, db.connections=%d, links.length=%d, email.enabled=%t, email.smtp.host=%s, email.smtp.port=%d, email.smtp.username=%s, email.smtp.from=%s, site.name=%s, site.url=%s",
			config.Server.Port, config.Database.Host, config.Database.Port, config.Database.User, config.Database.Database,
			config.Database.Connections, config.Links.Length, config.Email.Enabled, config.Email.SMTP.Host,
			config.Email.SMTP.Port, config.Email.SMTP.Username, config.Email.SMTP.From, config.Site.Name, config.Site.URL)

	return config, nil
}

func ensureDefaultValues(config *Config) {
	if config.Server.Port == 0 {
		config.Server.Port = 3000
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.Connections == 0 {
		config.Database.Connections = 10
	}
	if config.Links.Length == 0 {
		config.Links.Length = 6
	}
}
