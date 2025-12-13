package main

import (
	"log"
	"manga/config"
	"manga/database"
	"manga/ui"
)

func main() {
	// Check if config exists, if not show config dialog first
	if !config.ConfigExists() {
		log.Println("No configuration found. Please configure database settings.")
	}

	// Try to load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// If config doesn't exist or is invalid, create empty config
		cfg = &config.Config{
			DBServer: "localhost",
			DBPort:   "5432",
		}
	}

	// Create the UI (will show config dialog if needed)
	app := ui.NewApp(nil, cfg)

	// If we have a valid config, try to connect to database
	if cfg.DBServer != "" && cfg.DBName != "" {
		db, err := database.InitDB(cfg)
		if err != nil {
			log.Printf("Failed to initialize database: %v", err)
			// App will run but show config dialog
		} else {
			// Run migrations
			if err := database.RunMigrations(db); err != nil {
				log.Printf("Failed to run migrations: %v", err)
			} else {
				app.SetDatabase(db)
			}
		}
	}

	app.Run()
}
