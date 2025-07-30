package config

type Config struct {
	DB    DBConfig    `yaml:"db"`
	Redis RedisConfig `yaml:"redis"`
	SMTP  SMTPConfig  `yaml:"smtp"`
}

type DBConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
}
