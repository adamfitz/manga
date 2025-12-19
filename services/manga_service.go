package services

import (
	"database/sql"
	"fmt"
	"mangadb/models"
)

type MangaService struct {
	db *sql.DB
}

func NewMangaService(db *sql.DB) *MangaService {
	return &MangaService{db: db}
}

func (s *MangaService) GetAll() ([]models.Manga, error) {
	if s.db == nil {
		return []models.Manga{}, nil
	}

	query := `SELECT id, title, COALESCE(alt_title, ''), COALESCE(author, ''), 
			  COALESCE(description, ''), COALESCE(cover_url, ''), COALESCE(url, ''), 
			  COALESCE(mangadex_id, ''), status, created_at, updated_at 
			  FROM manga ORDER BY updated_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query manga: %w", err)
	}
	defer rows.Close()

	var mangas []models.Manga
	for rows.Next() {
		var m models.Manga
		if err := rows.Scan(&m.ID, &m.Title, &m.AltTitle, &m.Author, &m.Description,
			&m.CoverURL, &m.URL, &m.MangadexID, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan manga: %w", err)
		}
		mangas = append(mangas, m)
	}

	return mangas, nil
}

func (s *MangaService) GetByID(id int) (*models.Manga, error) {
	if s.db == nil {
		return nil, nil
	}

	query := `SELECT id, title, COALESCE(alt_title, ''), COALESCE(author, ''), 
			  COALESCE(description, ''), COALESCE(cover_url, ''), COALESCE(url, ''), 
			  COALESCE(mangadex_id, ''), status, created_at, updated_at 
			  FROM manga WHERE id = $1`

	var m models.Manga
	err := s.db.QueryRow(query, id).Scan(&m.ID, &m.Title, &m.AltTitle, &m.Author,
		&m.Description, &m.CoverURL, &m.URL, &m.MangadexID, &m.Status,
		&m.CreatedAt, &m.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get manga: %w", err)
	}

	return &m, nil
}

func (s *MangaService) GetByMangadexID(mangadexID string) (*models.Manga, error) {
	if s.db == nil {
		return nil, nil
	}

	query := `SELECT id, title, COALESCE(alt_title, ''), COALESCE(author, ''), 
			  COALESCE(description, ''), COALESCE(cover_url, ''), COALESCE(url, ''), 
			  COALESCE(mangadex_id, ''), status, created_at, updated_at 
			  FROM manga WHERE mangadex_id = $1`

	var m models.Manga
	err := s.db.QueryRow(query, mangadexID).Scan(&m.ID, &m.Title, &m.AltTitle,
		&m.Author, &m.Description, &m.CoverURL, &m.URL, &m.MangadexID, &m.Status,
		&m.CreatedAt, &m.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get manga by mangadex_id: %w", err)
	}

	return &m, nil
}

func (s *MangaService) Create(m *models.Manga) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	query := `INSERT INTO manga (title, alt_title, author, description, cover_url, 
			  url, mangadex_id, status) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			  RETURNING id, created_at, updated_at`

	err := s.db.QueryRow(query, m.Title, m.AltTitle, m.Author, m.Description,
		m.CoverURL, m.URL, m.MangadexID, m.Status).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create manga: %w", err)
	}

	return nil
}

func (s *MangaService) Update(m *models.Manga) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	query := `UPDATE manga SET title = $1, alt_title = $2, author = $3, 
			  description = $4, cover_url = $5, url = $6, mangadex_id = $7, 
			  status = $8, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $9`

	result, err := s.db.Exec(query, m.Title, m.AltTitle, m.Author, m.Description,
		m.CoverURL, m.URL, m.MangadexID, m.Status, m.ID)

	if err != nil {
		return fmt.Errorf("failed to update manga: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("manga not found")
	}

	return nil
}

func (s *MangaService) Delete(id int) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	query := `DELETE FROM manga WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete manga: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("manga not found")
	}

	return nil
}

func (s *MangaService) Search(searchTerm string) ([]models.Manga, error) {
	if s.db == nil {
		return []models.Manga{}, nil
	}

	query := `SELECT id, title, COALESCE(alt_title, ''), COALESCE(author, ''), 
			  COALESCE(description, ''), COALESCE(cover_url, ''), COALESCE(url, ''), 
			  COALESCE(mangadex_id, ''), status, created_at, updated_at 
			  FROM manga 
			  WHERE title ILIKE $1 OR author ILIKE $1 OR alt_title ILIKE $1
			  ORDER BY title`

	rows, err := s.db.Query(query, "%"+searchTerm+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search manga: %w", err)
	}
	defer rows.Close()

	var mangas []models.Manga
	for rows.Next() {
		var m models.Manga
		if err := rows.Scan(&m.ID, &m.Title, &m.AltTitle, &m.Author, &m.Description,
			&m.CoverURL, &m.URL, &m.MangadexID, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan manga: %w", err)
		}
		mangas = append(mangas, m)
	}

	return mangas, nil
}
