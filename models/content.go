package models

// Content is a generic model for all content types
type Content struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	AltName    string `json:"alt_name"`
	URL        string `json:"url"`
	MangadexID string `json:"mangadex_id,omitempty"` // Only for manga
	Status     string `json:"status"`
}

// ContentType represents the type of content
type ContentType string

const (
	TypeManga      ContentType = "manga"
	TypeAnime      ContentType = "anime"
	TypeLightNovel ContentType = "lightnovel"
	TypeWebNovel   ContentType = "webnovel"
	TypeWebtoons   ContentType = "webtoons"
)

// String returns the string representation of ContentType
func (ct ContentType) String() string {
	return string(ct)
}

// GetAllContentTypes returns all available content types
func GetAllContentTypes() []ContentType {
	return []ContentType{
		TypeManga,
		TypeAnime,
		TypeLightNovel,
		TypeWebNovel,
		TypeWebtoons,
	}
}

// HasMangadexID returns true if this content type supports mangadex_id
func (ct ContentType) HasMangadexID() bool {
	return ct == TypeManga
}
