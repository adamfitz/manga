package main

import (
	"log"
	"mangadb/assets/fonts" // imports follow folder structure
	"mangadb/config"
	"mangadb/database"
	"mangadb/ui"

	"fyne.io/fyne/v2"
)

func main() {
	// start logging
	config.Logger()

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

	// IMPORTANT: Use STATIC font files, NOT variable fonts (-VF)
	// Download from: https://github.com/notofonts/noto-cjk/releases
	// Use NotoSansCJKjp-Regular.otf (NOT NotoSansCJKjp-VF.otf)

	// Load fonts as resources
	// These fonts include Latin characters, so English will work!
	jpFont := fyne.NewStaticResource("NotoSansCJKjp-Regular.otf", fonts.JP)
	krFont := fyne.NewStaticResource("NotoSansCJKkr-Regular.otf", fonts.KR)
	scFont := fyne.NewStaticResource("NotoSansCJKsc-Regular.otf", fonts.SC)

	ui.InitCJKTheme(app.FyneApp, jpFont, krFont, scFont)

	app.Run()
}
