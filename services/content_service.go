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

	query := fmt.Sprintf(`SELECT %s FROM %s ORDER BY id DESC`, selectCols, m.table)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", contentType, err)
	}
	defer rows.Close()

	var list []models.Content

	for rows.Next() {

		switch contentType {

		case models.TypeManga:
			var c models.Content
			if err := rows.Scan(
				&c.ID,
				&c.Name,
				&c.AltName,
				&c.Author,
				&c.Description,
				&c.CoverURL,
				&c.URL,
				&c.MangadexID,
				&c.Status,
				new(interface{}), // created_at ignored
				new(interface{}), // updated_at ignored
			); err != nil {
				return nil, err
			}
			list = append(list, c)

		case models.TypeLightNovel:
			var ln models.LightNovel
			if err := rows.Scan(
				&ln.ID,
				&ln.Name,
				&ln.AltName,
				&ln.URL,
				&ln.Volumes,
				&ln.Status,
			); err != nil {
				return nil, err
			}

			list = append(list, models.Content{
				ID:      ln.ID,
				Name:    ln.Name,
				AltName: ln.AltName,
				URL:     ln.URL,
				Status:  ln.Status,
				Volumes: ln.Volumes,
			})

		default: // Anime, WebNovel, Webtoons
			var id int
			var name, alt, url, status string

			if err := rows.Scan(&id, &name, &alt, &url, &status); err != nil {
				return nil, err
			}

			list = append(list, models.Content{
				ID:      id,
				Name:    name,
				AltName: alt,
				URL:     url,
				Status:  status,
			})
		}
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
		selectCols, m.table, m.name, m.altName, m.name,
	)

	rows, err := s.db.Query(query, term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Content

	for rows.Next() {

		switch contentType {

		case models.TypeManga:
			var c models.Content
			if err := rows.Scan(
				&c.ID,
				&c.Name,
				&c.AltName,
				&c.Author,
				&c.Description,
				&c.CoverURL,
				&c.URL,
				&c.MangadexID,
				&c.Status,
				new(interface{}),
				new(interface{}),
			); err != nil {
				return nil, err
			}
			list = append(list, c)

		case models.TypeLightNovel:
			var ln models.LightNovel
			if err := rows.Scan(
				&ln.ID,
				&ln.Name,
				&ln.AltName,
				&ln.URL,
				&ln.Volumes,
				&ln.Status,
			); err != nil {
				return nil, err
			}

			list = append(list, models.Content{
				ID:      ln.ID,
				Name:    ln.Name,
				AltName: ln.AltName,
				URL:     ln.URL,
				Status:  ln.Status,
				Volumes: ln.Volumes,
			})

		default:
			var id int
			var name, alt, url, status string

			if err := rows.Scan(&id, &name, &alt, &url, &status); err != nil {
				return nil, err
			}

			list = append(list, models.Content{
				ID:      id,
				Name:    name,
				AltName: alt,
				URL:     url,
				Status:  status,
			})
		}
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

	var id int
	_, scanErr := fmt.Sscanf(lookupValue, "%d", &id)

	var query string
	if scanErr == nil {
		query = fmt.Sprintf(`SELECT %s FROM %s WHERE id=$1`, selectCols, m.table)
	} else {
		query = fmt.Sprintf(
			`SELECT %s FROM %s WHERE %s=$1 OR %s=$1 LIMIT 1`,
			selectCols, m.table, m.name, m.altName,
		)
	}

	switch contentType {

	case models.TypeManga:
		var c models.Content
		err = s.db.QueryRow(query, lookupValue).Scan(
			&c.ID,
			&c.Name,
			&c.AltName,
			&c.Author,
			&c.Description,
			&c.CoverURL,
			&c.URL,
			&c.MangadexID,
			&c.Status,
			new(interface{}),
			new(interface{}),
		)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &c, err

	case models.TypeLightNovel:
		var ln models.LightNovel
		err = s.db.QueryRow(query, lookupValue).Scan(
			&ln.ID,
			&ln.Name,
			&ln.AltName,
			&ln.URL,
			&ln.Volumes,
			&ln.Status,
		)
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return &models.Content{
			ID:      ln.ID,
			Name:    ln.Name,
			AltName: ln.AltName,
			URL:     ln.URL,
			Status:  ln.Status,
			Volumes: ln.Volumes,
		}, err

	default:
		var id int
		var name, alt, url, status string

		err = s.db.QueryRow(query, lookupValue).Scan(&id, &name, &alt, &url, &status)
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return &models.Content{
			ID:      id,
			Name:    name,
			AltName: alt,
			URL:     url,
			Status:  status,
		}, err
	}
}

// Create adds a new entry
func (s *ContentService) Create(contentType models.ContentType, content *models.Content) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	// -------------------------
	// UNIVERSAL VALIDATION RULES
	// -------------------------

	// Rule 1: Must have Name OR AltName
	if content.Name == "" && content.AltName == "" {
		return fmt.Errorf("either name or alt_name must be provided")
	}

	// Rule 2: Must have URL
	if content.URL == "" {
		return fmt.Errorf("url must be provided")
	}

	// -------------------------
	// SQL INSERT LOGIC
	// -------------------------

	m, ok := contentColumnMaps[contentType]
	if !ok {
		return fmt.Errorf("unsupported content type: %v", contentType)
	}

	var query string
	var err error

	switch contentType {

	case models.TypeManga:
		query = fmt.Sprintf(`
            INSERT INTO %s (title, alt_title, author, description, cover_url, url, mangadex_id, status)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
            RETURNING id`,
			m.table,
		)
		err = s.db.QueryRow(
			query,
			content.Name,
			content.AltName,
			content.Author,
			content.Description,
			content.CoverURL,
			content.URL,
			content.MangadexID,
			content.Status,
		).Scan(&content.ID)

	case models.TypeLightNovel:
		query = fmt.Sprintf(`
            INSERT INTO %s (name, alt_name, url, volumes, status)
            VALUES ($1,$2,$3,$4,$5)
            RETURNING id`,
			m.table,
		)
		err = s.db.QueryRow(
			query,
			content.Name,
			content.AltName,
			content.URL,
			content.Volumes,
			content.Status,
		).Scan(&content.ID)

	default: // Anime, WebNovel, Webtoons
		query = fmt.Sprintf(`
            INSERT INTO %s (name, alt_name, url, status)
            VALUES ($1,$2,$3,$4)
            RETURNING id`,
			m.table,
		)
		err = s.db.QueryRow(
			query,
			content.Name,
			content.AltName,
			content.URL,
			content.Status,
		).Scan(&content.ID)
	}

	if err != nil {
		return fmt.Errorf("failed to create %s: %w", contentType, err)
	}

	return nil
}

// Update modifies an existing entry
func (s *ContentService) Update(contentType models.ContentType, content *models.Content) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	// -------------------------
	// UNIVERSAL VALIDATION RULES
	// -------------------------

	// Rule 1: Must have Name OR AltName
	if content.Name == "" && content.AltName == "" {
		return fmt.Errorf("either name or alt_name must be provided")
	}

	// Rule 2: Must have URL
	if content.URL == "" {
		return fmt.Errorf("url must be provided")
	}

	// -------------------------
	// SQL UPDATE LOGIC
	// -------------------------

	m, ok := contentColumnMaps[contentType]
	if !ok {
		return fmt.Errorf("unsupported content type: %v", contentType)
	}

	var query string
	var result sql.Result
	var err error

	switch contentType {

	case models.TypeManga:
		query = fmt.Sprintf(`
            UPDATE %s
            SET title=$1, alt_title=$2, author=$3, description=$4, cover_url=$5, url=$6, mangadex_id=$7, status=$8
            WHERE id=$9`,
			m.table,
		)
		result, err = s.db.Exec(
			query,
			content.Name,
			content.AltName,
			content.Author,
			content.Description,
			content.CoverURL,
			content.URL,
			content.MangadexID,
			content.Status,
			content.ID,
		)

	case models.TypeLightNovel:
		query = fmt.Sprintf(`
            UPDATE %s
            SET name=$1, alt_name=$2, url=$3, volumes=$4, status=$5
            WHERE id=$6`,
			m.table,
		)
		result, err = s.db.Exec(
			query,
			content.Name,
			content.AltName,
			content.URL,
			content.Volumes,
			content.Status,
			content.ID,
		)

	default: // Anime, WebNovel, Webtoons
		query = fmt.Sprintf(`
            UPDATE %s
            SET name=$1, alt_name=$2, url=$3, status=$4
            WHERE id=$5`,
			m.table,
		)
		result, err = s.db.Exec(
			query,
			content.Name,
			content.AltName,
			content.URL,
			content.Status,
			content.ID,
		)
	}

	if err != nil {
		return fmt.Errorf("failed to update %s: %w", contentType, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%s with id %d not found", contentType, content.ID)
	}

	return nil
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
