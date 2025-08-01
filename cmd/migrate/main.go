package main

import (
	"flag"
	"fmt"
	"log"

	"jcourse_go/internal/config"
	"jcourse_go/internal/infrastructure/database"
	"jcourse_go/internal/infrastructure/migrations"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/config.yaml", "Path to config file")
	flag.Parse()

	conf, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDatabase(conf.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations completed successfully!")
}
