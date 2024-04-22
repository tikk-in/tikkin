package config

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

type ConfigFlags struct {
	ConfigPath    string
	AdminPassword string
	SMTPPassword  string
	JWTSecret     string
}
