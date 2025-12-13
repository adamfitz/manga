package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// check config directory exists or create it
func verifyConfigDirectory() (string, error) {
	configDirectory, expandError := expandPath("~/.config/mangaDbApp")
	if expandError != nil {
		log.Fatalf("cannot verify local configuration directory: %v", expandError)
	}

	// Check if the directory exists
	_, err := os.Stat(configDirectory)

	if os.IsNotExist(err) {
		// Create the directory with read/write/execute permissions for owner, and read/execute for others
		err := os.MkdirAll(configDirectory, 0755)
		if err != nil {
			log.Fatalf("error creating directory %s: %v", configDirectory, err)
		}
		log.Printf("Directory %s created successfully.\n", configDirectory)
	} else if err != nil {
		log.Fatalf("error checking directory %s: %v", configDirectory, err)
	}

	return configDirectory, nil
}

// logging configuration
func Logger() error {
	// local config dir
	configDir, configDirErr := verifyConfigDirectory()
	if configDirErr != nil {
		return configDirErr
	}
	// open log file or creat it if it does not exist
	logFilePath := fmt.Sprintf("%s/mangaDbApp.log", configDir)
	logFile, logFileErr := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if logFileErr != nil {
		log.Fatalf("Failed to open log file: %v", logFileErr)
	}

	log.SetOutput(logFile)

	return nil
}

// expands ~ to the user's home directory, or returns the path as-is
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		// Path starts with ~/ so expand it
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	// Path doesn't start with ~/ so return it unchanged
	return path, nil
}
