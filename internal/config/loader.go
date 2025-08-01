package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

func LoadFromEnv() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	return Load(configPath)
}
