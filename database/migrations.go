package database

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS manga (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255),
			description TEXT,
			cover_url TEXT,
			mangadex_id VARCHAR(255) UNIQUE,
			status VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS chapters (
			id SERIAL PRIMARY KEY,
			manga_id INTEGER REFERENCES manga(id) ON DELETE CASCADE,
			chapter_number VARCHAR(50) NOT NULL,
			title VARCHAR(255),
			url TEXT,
			read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(manga_id, chapter_number)
		)`,
		`CREATE TABLE IF NOT EXISTS bookmarks (
			id SERIAL PRIMARY KEY,
			manga_id INTEGER REFERENCES manga(id) ON DELETE CASCADE,
			chapter_id INTEGER REFERENCES chapters(id) ON DELETE SET NULL,
			note TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(manga_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_manga_title ON manga(title)`,
		`CREATE INDEX IF NOT EXISTS idx_chapters_manga_id ON chapters(manga_id)`,
		`CREATE INDEX IF NOT EXISTS idx_bookmarks_manga_id ON bookmarks(manga_id)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	return nil
}
