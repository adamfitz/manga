package models

type WebNovel struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	AltName string `json:"alt_name"`
	URL     string `json:"url"`
	Status  string `json:"status"`
}
