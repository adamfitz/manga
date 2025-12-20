package models

import "time"

type Bookmark struct {
	ID        int       `json:"id"`
	MangaID   int       `json:"manga_id"`
	ChapterID *int      `json:"chapter_id"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	Manga     *Manga    `json:"manga,omitempty"`
}
