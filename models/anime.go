package models

type Anime struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AltName   string `json:"alt_name"`
	URL       string `json:"url"`
	Completed bool   `json:"completed"`
	Watched   bool   `json:"watched"`
	Status    string `json:"status"`
}
