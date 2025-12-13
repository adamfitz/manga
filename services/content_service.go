package services

import (
	"database/sql"
	"fmt"
	"manga/models"
)

type ContentService struct {
	db *sql.DB
}

func NewContentService(db *sql.DB) *ContentService {
	return &ContentService{db: db}
}

// GetAll retrieves all entries from the specified table
func (s *ContentService) GetAll(contentType models.ContentType) ([]models.Content, error) {
	var query string
	if contentType.HasMangadexID() {
		query = fmt.Sprintf(`SELECT id, name, COALESCE(alt_name, ''), url, COALESCE(mangadex_id, ''), status, created_at, updated_at 
							 FROM %s ORDER BY updated_at DESC`, contentType)
	} else {
		query = fmt.Sprintf(`SELECT id, name, COALESCE(alt_name, ''), url, status, created_at, updated_at 
							 FROM %s ORDER BY updated_at DESC`, contentType)
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", contentType, err)
	}
	defer rows.Close()

	var contents []models.Content
	for rows.Next() {
		var c models.Content
		if contentType.HasMangadexID() {
			err := rows.Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.MangadexID, &c.Status, &c.CreatedAt, &c.UpdatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to scan %s: %w", contentType, err)
			}
		} else {
			err := rows.Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.Status, &c.CreatedAt, &c.UpdatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to scan %s: %w", contentType, err)
			}
		}
		contents = append(contents, c)
	}

	return contents, nil
}

// Search performs substring search on name and alt_name
func (s *ContentService) Search(contentType models.ContentType, searchTerm string) ([]models.Content, error) {
	var query string
	if contentType.HasMangadexID() {
		query = fmt.Sprintf(`SELECT id, name, COALESCE(alt_name, ''), url, COALESCE(mangadex_id, ''), status, created_at, updated_at 
							 FROM %s 
							 WHERE name ILIKE $1 OR alt_name ILIKE $1 
							 ORDER BY name`, contentType)
	} else {
		query = fmt.Sprintf(`SELECT id, name, COALESCE(alt_name, ''), url, status, created_at, updated_at 
							 FROM %s 
							 WHERE name ILIKE $1 OR alt_name ILIKE $1 
							 ORDER BY name`, contentType)
	}

	rows, err := s.db.Query(query, "%"+searchTerm+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search %s: %w", contentType, err)
	}
	defer rows.Close()

	var contents []models.Content
	for rows.Next() {
		var c models.Content
		if contentType.HasMangadexID() {
			err := rows.Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.MangadexID, &c.Status, &c.CreatedAt, &c.UpdatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to scan %s: %w", contentType, err)
			}
		} else {
			err := rows.Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.Status, &c.CreatedAt, &c.UpdatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to scan %s: %w", contentType, err)
			}
		}
		contents = append(contents, c)
	}

	return contents, nil
}

// Lookup performs exact match on name, alt_name, or ID
func (s *ContentService) Lookup(contentType models.ContentType, lookupValue string) (*models.Content, error) {
	var query string
	var c models.Content

	// Try to parse as ID first
	var id int
	_, err := fmt.Sscanf(lookupValue, "%d", &id)

	if err == nil {
		// It's a number, do ID lookup
		if contentType.HasMangadexID() {
			query = fmt.Sprintf(`SELECT id, name, COALESCE(alt_name, ''), url, COALESCE(mangadex_id, ''), status, created_at, updated_at 
								 FROM %s WHERE id = $1`, contentType)
			err = s.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.MangadexID, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		} else {
			query = fmt.Sprintf(`SELECT id, name, COALESCE(alt_name, ''), url, status, created_at, updated_at 
								 FROM %s WHERE id = $1`, contentType)
			err = s.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		}
	} else {
		// It's a string, do exact name/alt_name lookup
		if contentType.HasMangadexID() {
			query = fmt.Sprintf(`SELECT id, name, COALESCE(alt_name, ''), url, COALESCE(mangadex_id, ''), status, created_at, updated_at 
								 FROM %s WHERE name = $1 OR alt_name = $1 LIMIT 1`, contentType)
			err = s.db.QueryRow(query, lookupValue).Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.MangadexID, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		} else {
			query = fmt.Sprintf(`SELECT id, name, COALESCE(alt_name, ''), url, status, created_at, updated_at 
								 FROM %s WHERE name = $1 OR alt_name = $1 LIMIT 1`, contentType)
			err = s.db.QueryRow(query, lookupValue).Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		}
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no %s found with that identifier", contentType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to lookup %s: %w", contentType, err)
	}

	return &c, nil
}

// Create adds a new entry to the specified table
func (s *ContentService) Create(contentType models.ContentType, c *models.Content) error {
	var query string
	var err error

	// Handle NULL for empty strings
	var altName interface{}
	if c.AltName == "" {
		altName = nil
	} else {
		altName = c.AltName
	}

	if contentType.HasMangadexID() {
		var mangadexID interface{}
		if c.MangadexID == "" {
			mangadexID = nil
		} else {
			mangadexID = c.MangadexID
		}

		query = fmt.Sprintf(`INSERT INTO %s (name, alt_name, url, mangadex_id, status) 
							 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`, contentType)
		err = s.db.QueryRow(query, c.Name, altName, c.URL, mangadexID, c.Status).
			Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	} else {
		query = fmt.Sprintf(`INSERT INTO %s (name, alt_name, url, status) 
							 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`, contentType)
		err = s.db.QueryRow(query, c.Name, altName, c.URL, c.Status).
			Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	}

	if err != nil {
		return fmt.Errorf("failed to create %s: %w", contentType, err)
	}

	return nil
}

// Delete removes an entry from the specified table
func (s *ContentService) Delete(contentType models.ContentType, id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, contentType)

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", contentType, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%s not found", contentType)
	}

	return nil
}
