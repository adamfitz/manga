package services

import (
	"encoding/json"
	"fmt"
	"io"
	"manga/models"
	"net/http"
	"net/url"
)

type MangadexService struct {
	baseURL string
	client  *http.Client
}

func NewMangadexService() *MangadexService {
	return &MangadexService{
		baseURL: "https://api.mangadex.org",
		client:  &http.Client{},
	}
}

type MangadexResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Title       map[string]string `json:"title"`
			Description map[string]string `json:"description"`
			Status      string            `json:"status"`
		} `json:"attributes"`
		Relationships []struct {
			Type       string `json:"type"`
			Attributes struct {
				Name string `json:"name"`
			} `json:"attributes"`
		} `json:"relationships"`
	} `json:"data"`
}

func (s *MangadexService) Search(title string) ([]models.Manga, error) {
	searchURL := fmt.Sprintf("%s/manga?title=%s&limit=20", s.baseURL, url.QueryEscape(title))

	resp, err := s.client.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search mangadex: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mangadex API error: %s", string(body))
	}

	var mdResp MangadexResponse
	if err := json.NewDecoder(resp.Body).Decode(&mdResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var mangas []models.Manga
	for _, item := range mdResp.Data {
		manga := models.Manga{
			MangadexID: item.ID,
			Status:     item.Attributes.Status,
		}

		// Get English title or first available
		if t, ok := item.Attributes.Title["en"]; ok {
			manga.Title = t
		} else {
			for _, t := range item.Attributes.Title {
				manga.Title = t
				break
			}
		}

		// Get English description or first available
		if d, ok := item.Attributes.Description["en"]; ok {
			manga.Description = d
		} else {
			for _, d := range item.Attributes.Description {
				manga.Description = d
				break
			}
		}

		// Get author
		for _, rel := range item.Relationships {
			if rel.Type == "author" {
				manga.Author = rel.Attributes.Name
				break
			}
		}

		manga.CoverURL = fmt.Sprintf("https://mangadex.org/title/%s", item.ID)
		mangas = append(mangas, manga)
	}

	return mangas, nil
}
