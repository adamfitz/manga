package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	PgServer   string `json:"db_server"`
	PgPort     string `json:"db_port"`
	PgUser     string `json:"db_user"`
	PgPassword string `json:"db_user_pass"`
	PgDbName   string `json:"db_name"`
}

// load config
func LoadConfig() (Config, error) {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	//construct path to config gile
	configPath := fmt.Sprintf("%s/.config", homeDir)
	if err != nil {
		log.Println("Error getting home directory:", err)
		return Config{}, err
	}

	// Build the path to the config file
	configFile := filepath.Join(configPath, "manga.config")

	// Open the JSON file
	file, err := os.Open(configFile)
	if err != nil {
		log.Println("Error opening file:", err)
		return Config{}, err
	}
	defer file.Close()

	// Decode the JSON file
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Println("Error decoding JSON:", err)
		return Config{}, err
	}

	// Return the config and no error
	return config, nil
}
