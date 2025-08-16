package config

type Config struct {
	DB    DBConfig    `yaml:"db"`
	SMTP  SMTPConfig  `yaml:"smtp"`
	Event EventConfig `yaml:"event"`
}

type DBConfig struct {
	DSN string `yaml:"dsn"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
}

type EventConfig struct {
	Enabled bool `yaml:"enabled"`
}
