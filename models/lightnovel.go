package models

type LightNovel struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	AltName string `json:"alt_name"`
	URL     string `json:"url"`
	Volumes int    `json:"volumes"`
	Status  string `json:"status"`
}
