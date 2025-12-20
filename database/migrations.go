package database

import (
	"database/sql"
	"fmt"
	"log"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		// -------------------- MANGA (UNIFIED) --------------------
		`CREATE TABLE IF NOT EXISTS manga (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			alt_title VARCHAR(255),
			author VARCHAR(255),
			description TEXT,
			cover_url TEXT,
			url TEXT NOT NULL,
			mangadex_id VARCHAR(255),
			status VARCHAR(50) NOT NULL CHECK (status IN ('ongoing', 'completed', 'hiatus', 'cancelled')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// -------------------- OTHER MEDIA --------------------
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

		`CREATE TABLE IF NOT EXISTS webnovel (
		id SERIAL PRIMARY KEY,
		name VARCHAR(500) NOT NULL,
		alt_name VARCHAR(500),
		url TEXT,
		status VARCHAR(50) DEFAULT 'unknown',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// -------------------- INDEXES --------------------
		`CREATE INDEX IF NOT EXISTS idx_manga_title ON manga(title)`,
		`CREATE INDEX IF NOT EXISTS idx_manga_alt_title ON manga(alt_title)`,

		`CREATE INDEX IF NOT EXISTS idx_anime_name ON anime(name)`,
		`CREATE INDEX IF NOT EXISTS idx_lightnovel_name ON lightnovel(name)`,
		`CREATE INDEX IF NOT EXISTS idx_webtoons_name ON webtoons(name)`,
		`CREATE INDEX IF NOT EXISTS idx_webtoons_alt_name ON webtoons(alt_name)`,
		`CREATE INDEX IF NOT EXISTS idx_webnovel_name ON webnovel(name)`,
		`CREATE INDEX IF NOT EXISTS idx_webnovel_alt_name ON webnovel(alt_name)`,
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
