package services

import (
	"database/sql"
	"fmt"
	"mangadb/models"
)

type BookmarkService struct {
	db *sql.DB
}

func NewBookmarkService(db *sql.DB) *BookmarkService {
	return &BookmarkService{db: db}
}

func (s *BookmarkService) GetAll() ([]models.Bookmark, error) {
	// Nil DB check
	if s.db == nil {
		return []models.Bookmark{}, nil
	}

	query := `SELECT b.id, b.manga_id, b.chapter_id, b.note, b.created_at,
			  m.id, m.title, m.author, m.description, m.cover_url, m.mangadex_id, m.status, m.created_at, m.updated_at
			  FROM bookmarks b
			  JOIN manga m ON b.manga_id = m.id
			  ORDER BY b.created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookmarks: %w", err)
	}
	defer rows.Close()

	var bookmarks []models.Bookmark
	for rows.Next() {
		var b models.Bookmark
		var m models.Manga
		var chapterID sql.NullInt64

		err := rows.Scan(&b.ID, &b.MangaID, &chapterID, &b.Note, &b.CreatedAt,
			&m.ID, &m.Title, &m.Author, &m.Description, &m.CoverURL, &m.MangadexID,
			&m.Status, &m.CreatedAt, &m.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan bookmark: %w", err)
		}

		if chapterID.Valid {
			id := int(chapterID.Int64)
			b.ChapterID = &id
		}

		b.Manga = &m
		bookmarks = append(bookmarks, b)
	}

	return bookmarks, nil
}

func (s *BookmarkService) Create(b *models.Bookmark) error {
	// Nil DB check
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	query := `INSERT INTO bookmarks (manga_id, chapter_id, note) 
			  VALUES ($1, $2, $3) RETURNING id, created_at`

	var chapterID interface{}
	if b.ChapterID != nil {
		chapterID = *b.ChapterID
	}

	err := s.db.QueryRow(query, b.MangaID, chapterID, b.Note).
		Scan(&b.ID, &b.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create bookmark: %w", err)
	}

	return nil
}

func (s *BookmarkService) Delete(id int) error {
	// Nil DB check
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}

	query := `DELETE FROM bookmarks WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("bookmark not found")
	}

	return nil
}

func (s *BookmarkService) GetByMangaID(mangaID int) (*models.Bookmark, error) {
	// Nil DB check
	if s.db == nil {
		return nil, nil
	}

	query := `SELECT id, manga_id, chapter_id, note, created_at 
			  FROM bookmarks WHERE manga_id = $1`

	var b models.Bookmark
	var chapterID sql.NullInt64

	err := s.db.QueryRow(query, mangaID).Scan(&b.ID, &b.MangaID, &chapterID, &b.Note, &b.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmark: %w", err)
	}

	if chapterID.Valid {
		id := int(chapterID.Int64)
		b.ChapterID = &id
	}

	return &b, nil
}
