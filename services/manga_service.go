package services

import (
	"database/sql"
	"fmt"
	"manga/models"
)

type MangaService struct {
	db *sql.DB
}

func NewMangaService(db *sql.DB) *MangaService {
	return &MangaService{db: db}
}

func (s *MangaService) GetAll() ([]models.Manga, error) {
	query := `SELECT id, title, author, description, cover_url, mangadex_id, status, created_at, updated_at 
			  FROM manga ORDER BY updated_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query manga: %w", err)
	}
	defer rows.Close()

	var mangas []models.Manga
	for rows.Next() {
		var m models.Manga
		err := rows.Scan(&m.ID, &m.Title, &m.Author, &m.Description, &m.CoverURL,
			&m.MangadexID, &m.Status, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan manga: %w", err)
		}
		mangas = append(mangas, m)
	}

	return mangas, nil
}

func (s *MangaService) GetByID(id int) (*models.Manga, error) {
	query := `SELECT id, title, author, description, cover_url, mangadex_id, status, created_at, updated_at 
			  FROM manga WHERE id = $1`

	var m models.Manga
	err := s.db.QueryRow(query, id).Scan(&m.ID, &m.Title, &m.Author, &m.Description,
		&m.CoverURL, &m.MangadexID, &m.Status, &m.CreatedAt, &m.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("manga not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get manga: %w", err)
	}

	return &m, nil
}

func (s *MangaService) Create(m *models.Manga) error {
	query := `INSERT INTO manga (title, author, description, cover_url, mangadex_id, status) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	err := s.db.QueryRow(query, m.Title, m.Author, m.Description, m.CoverURL, m.MangadexID, m.Status).
		Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create manga: %w", err)
	}

	return nil
}

func (s *MangaService) Update(m *models.Manga) error {
	query := `UPDATE manga SET title = $1, author = $2, description = $3, cover_url = $4, 
			  mangadex_id = $5, status = $6, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $7`

	result, err := s.db.Exec(query, m.Title, m.Author, m.Description, m.CoverURL,
		m.MangadexID, m.Status, m.ID)

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
	query := `SELECT id, title, author, description, cover_url, mangadex_id, status, created_at, updated_at 
			  FROM manga WHERE title ILIKE $1 OR author ILIKE $1 ORDER BY title`

	rows, err := s.db.Query(query, "%"+searchTerm+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search manga: %w", err)
	}
	defer rows.Close()

	var mangas []models.Manga
	for rows.Next() {
		var m models.Manga
		err := rows.Scan(&m.ID, &m.Title, &m.Author, &m.Description, &m.CoverURL,
			&m.MangadexID, &m.Status, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan manga: %w", err)
		}
		mangas = append(mangas, m)
	}

	return mangas, nil
}
