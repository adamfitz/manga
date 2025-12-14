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

// GetAll retrieves all entries for the given content type
func (s *ContentService) GetAll(contentType models.ContentType) ([]models.Content, error) {
	if s.db == nil {
		return []models.Content{}, nil
	}

	selectCols, m, err := normalizedSelect(contentType)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(
		`SELECT %s FROM %s ORDER BY id DESC`,
		selectCols,
		m.table,
	)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", contentType, err)
	}
	defer rows.Close()

	var list []models.Content
	for rows.Next() {
		var c models.Content
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.AltName,
			&c.URL,
			&c.MangadexID,
			&c.Author,
			&c.Description,
			&c.CoverURL,
			&c.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan %s: %w", contentType, err)
		}
		list = append(list, c)
	}

	return list, nil
}

// Search performs substring search on name and alt_name
func (s *ContentService) Search(contentType models.ContentType, searchTerm string) ([]models.Content, error) {
	if s.db == nil {
		return []models.Content{}, nil
	}

	selectCols, m, err := normalizedSelect(contentType)
	if err != nil {
		return nil, err
	}

	term := "%" + searchTerm + "%"

	query := fmt.Sprintf(
		`SELECT %s FROM %s
		 WHERE %s ILIKE $1 OR %s ILIKE $1
		 ORDER BY %s`,
		selectCols,
		m.table,
		m.name,
		m.altName,
		m.name,
	)

	rows, err := s.db.Query(query, term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Content
	for rows.Next() {
		var c models.Content
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.AltName,
			&c.URL,
			&c.MangadexID,
			&c.Author,
			&c.Description,
			&c.CoverURL,
			&c.Status,
		); err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	return list, nil
}

// Lookup performs exact match on name, alt_name, or ID
func (s *ContentService) Lookup(contentType models.ContentType, lookupValue string) (*models.Content, error) {
	if s.db == nil {
		return nil, nil
	}

	selectCols, m, err := normalizedSelect(contentType)
	if err != nil {
		return nil, err
	}

	var c models.Content
	var id int

	_, scanErr := fmt.Sscanf(lookupValue, "%d", &id)

	if scanErr == nil {
		// Lookup by ID
		query := fmt.Sprintf(
			`SELECT %s FROM %s WHERE id=$1`,
			selectCols,
			m.table,
		)

		err = s.db.QueryRow(query, id).Scan(
			&c.ID,
			&c.Name,
			&c.AltName,
			&c.URL,
			&c.MangadexID,
			&c.Author,
			&c.Description,
			&c.CoverURL,
			&c.Status,
		)

	} else {
		// Lookup by name / alt_name
		query := fmt.Sprintf(
			`SELECT %s FROM %s
			 WHERE %s=$1 OR %s=$1
			 LIMIT 1`,
			selectCols,
			m.table,
			m.name,
			m.altName,
		)

		err = s.db.QueryRow(query, lookupValue).Scan(
			&c.ID,
			&c.Name,
			&c.AltName,
			&c.URL,
			&c.MangadexID,
			&c.Author,
			&c.Description,
			&c.CoverURL,
			&c.Status,
		)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Delete removes an entry
func (s *ContentService) Delete(contentType models.ContentType, id int) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	m, ok := contentColumnMaps[contentType]
	if !ok {
		return fmt.Errorf("unsupported content type: %v", contentType)
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, m.table)
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", contentType, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%s with id %d not found", contentType, id)
	}

	return nil
}
