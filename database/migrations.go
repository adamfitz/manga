package database

import (
	"database/sql"
	"fmt"
	"log"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS manga (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			alt_name VARCHAR(255),
			url TEXT NOT NULL,
			mangadex_id VARCHAR(255),
			status VARCHAR(50) NOT NULL CHECK (status IN ('ongoing', 'completed', 'hiatus', 'cancelled')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS mangadex (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			alt_name VARCHAR(255),
			url TEXT NOT NULL,
			mangadex_id VARCHAR(255),
			status VARCHAR(50) NOT NULL CHECK (status IN ('ongoing', 'completed', 'hiatus', 'cancelled')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS anime (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			alt_name VARCHAR(255),
			url TEXT NOT NULL,
			status VARCHAR(50) NOT NULL CHECK (status IN ('ongoing', 'completed', 'hiatus', 'cancelled')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS lightnovel (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			alt_name VARCHAR(255),
			url TEXT NOT NULL,
			status VARCHAR(50) NOT NULL CHECK (status IN ('ongoing', 'completed', 'hiatus', 'cancelled')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS webnovel (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			alt_name VARCHAR(255),
			url TEXT NOT NULL,
			status VARCHAR(50) NOT NULL CHECK (status IN ('ongoing', 'completed', 'hiatus', 'cancelled')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS webtoons (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			alt_name VARCHAR(255),
			url TEXT NOT NULL,
			status VARCHAR(50) NOT NULL CHECK (status IN ('ongoing', 'completed', 'hiatus', 'cancelled')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_manga_name ON manga(name)`,
		`CREATE INDEX IF NOT EXISTS idx_manga_alt_name ON manga(alt_name)`,
		`CREATE INDEX IF NOT EXISTS idx_mangadex_name ON mangadex(name)`,
		`CREATE INDEX IF NOT EXISTS idx_mangadex_alt_name ON mangadex(alt_name)`,
		`CREATE INDEX IF NOT EXISTS idx_anime_name ON anime(name)`,
		`CREATE INDEX IF NOT EXISTS idx_lightnovel_name ON lightnovel(name)`,
		`CREATE INDEX IF NOT EXISTS idx_webnovel_name ON webnovel(name)`,
		`CREATE INDEX IF NOT EXISTS idx_webtoons_name ON webtoons(name)`,
	}

	for i, migration := range migrations {
		log.Printf("Running migration %d:\n%s\n", i+1, migration)

		if _, err := db.Exec(migration); err != nil {
			log.Printf("Migration %d failed!\nSQL: %s\nError: %v\n", i+1, migration, err)
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
		log.Printf("Migration %d succeeded.\n", i+1)
	}

	return nil
}
