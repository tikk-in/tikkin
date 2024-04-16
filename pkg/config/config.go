package config

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"os"
)

type JwtConfig struct {
	Secret string
}

type ServerConfig struct {
	Port int       `yaml:"port"`
	Jwt  JwtConfig `yaml:"jwt"`
}

type AdminConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type DatabaseConfig struct {
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Database    string `yaml:"database"`
	Connections int32  `yaml:"connections"`
}

type LinksConfig struct {
	Length int `yaml:"length"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type EmailConfig struct {
	Enabled bool       `yaml:"enabled"`
	SMTP    SMTPConfig `yaml:"smtp"`
}

type SiteConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Site     SiteConfig     `yaml:"site"`
	Email    EmailConfig    `yaml:"email"`
	Admin    AdminConfig    `yaml:"admin"`
	Database DatabaseConfig `yaml:"db"`
	Links    LinksConfig    `yaml:"links"`
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

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// Open config file
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
	log.Info().
		Msgf("Loaded config: server.port=%d, server.jwt.secret=<redacted>, db.host=%s,  db.port=%d, db.user=%s, db.password=<redacted>, db.database=%s, db.connections=%d, links.length=%d",
			config.Server.Port, config.Database.Host, config.Database.Port, config.Database.User, config.Database.Database, config.Database.Connections, config.Links.Length)
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
		config.Links.Length = 7
	}
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
