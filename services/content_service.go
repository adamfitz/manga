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

	var query string
	switch contentType {
	case models.TypeManga:
		query = `SELECT id, title as name, COALESCE(alt_title,'') as alt_name, 
			 COALESCE(url,'') as url, COALESCE(mangadex_id,'') as mangadex_id, 
			 COALESCE(author,'') as author, COALESCE(description,'') as description,
			 COALESCE(cover_url,'') as cover_url, status
			 FROM manga ORDER BY id DESC`

	case models.TypeAnime:
		query = `SELECT id, name, COALESCE(alt_name,''), url, '' as mangadex_id, status
				 FROM anime ORDER BY id DESC`
	case models.TypeLightNovel:
		query = `SELECT id, name, COALESCE(alt_name,''), url, '' as mangadex_id, status
				 FROM lightnovel ORDER BY id DESC`
	default:
		return nil, fmt.Errorf("unsupported content type: %v", contentType)
	}

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", contentType, err)
	}
	defer rows.Close()

	var list []models.Content
	for rows.Next() {
		var c models.Content
		if err := rows.Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.MangadexID,
			&c.Author, &c.Description, &c.CoverURL, &c.Status); err != nil {
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

	term := "%" + searchTerm + "%"
	var query string

	switch contentType {
	case models.TypeManga:
		query = `SELECT id, title as name, COALESCE(alt_title,'') as alt_name, 
			 COALESCE(url,'') as url, COALESCE(mangadex_id,'') as mangadex_id, 
			 COALESCE(author,'') as author, COALESCE(description,'') as description,
			 COALESCE(cover_url,'') as cover_url, status
			 FROM manga 
			 WHERE title ILIKE $1 OR alt_title ILIKE $1 OR author ILIKE $1
			 ORDER BY title`

	case models.TypeAnime:
		query = `SELECT id, name, COALESCE(alt_name,''), url, '' as mangadex_id, status
				 FROM anime WHERE name ILIKE $1 OR alt_name ILIKE $1 ORDER BY name`
	case models.TypeLightNovel:
		query = `SELECT id, name, COALESCE(alt_name,''), url, '' as mangadex_id, status
				 FROM lightnovel WHERE name ILIKE $1 OR alt_name ILIKE $1 ORDER BY name`
	default:
		return nil, fmt.Errorf("unsupported content type: %v", contentType)
	}

	rows, err := s.db.Query(query, term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Content
	for rows.Next() {
		var c models.Content
		if err := rows.Scan(&c.ID, &c.Name, &c.AltName, &c.URL, &c.MangadexID,
			&c.Author, &c.Description, &c.CoverURL, &c.Status); err != nil {
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

	var id int
	var query string
	var c models.Content

	_, err := fmt.Sscanf(lookupValue, "%d", &id)
	if err == nil {
		// Lookup by ID
		switch contentType {
		case models.TypeManga:
			query = `SELECT id, title as name, COALESCE(alt_title,'') as alt_name, 
			 COALESCE(url,'') as url, COALESCE(mangadex_id,'') as mangadex_id, 
			 COALESCE(author,'') as author, COALESCE(description,'') as description,
			 COALESCE(cover_url,'') as cover_url, status
			 FROM manga WHERE id=$1`

		case models.TypeAnime:
			query = `SELECT id, name, COALESCE(alt_name,''), url, '' as mangadex_id, status
					 FROM anime WHERE id=$1`
		case models.TypeLightNovel:
			query = `SELECT id, name, COALESCE(alt_name,''), url, '' as mangadex_id, status
					 FROM lightnovel WHERE id=$1`
		default:
			return nil, fmt.Errorf("unsupported content type: %v", contentType)
		}

		err = s.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.AltName, &c.URL,
			&c.MangadexID, &c.Author, &c.Description, &c.CoverURL, &c.Status)

	} else {
		// Lookup by name/alt_name
		switch contentType {
		case models.TypeManga:
			query = `SELECT id, title as name, COALESCE(alt_title,'') as alt_name, 
			 COALESCE(url,'') as url, COALESCE(mangadex_id,'') as mangadex_id, 
			 COALESCE(author,'') as author, COALESCE(description,'') as description,
			 COALESCE(cover_url,'') as cover_url, status
			 FROM manga WHERE title=$1 OR alt_title=$1 LIMIT 1`
		case models.TypeAnime:
			query = `SELECT id, name, COALESCE(alt_name,''), url, '' as mangadex_id, status
					 FROM anime WHERE name=$1 OR alt_name=$1 LIMIT 1`
		case models.TypeLightNovel:
			query = `SELECT id, name, COALESCE(alt_name,''), url, '' as mangadex_id, status
					 FROM lightnovel WHERE name=$1 OR alt_name=$1 LIMIT 1`
		default:
			return nil, fmt.Errorf("unsupported content type: %v", contentType)
		}

		err = s.db.QueryRow(query, lookupValue).Scan(&c.ID, &c.Name, &c.AltName, &c.URL,
			&c.MangadexID, &c.Author, &c.Description, &c.CoverURL, &c.Status)
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

	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, contentType)
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
