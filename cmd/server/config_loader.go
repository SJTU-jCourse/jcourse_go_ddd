package main

import (
	"flag"
	"log"

	"jcourse_go/internal/config"
)

func LoadConfiguration() *config.Config {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return cfg
}