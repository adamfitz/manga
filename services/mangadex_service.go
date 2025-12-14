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

// Search searches for manga on MangaDex API and returns API models
func (s *MangadexService) Search(title string) ([]models.MangadexManga, error) {
	searchURL := fmt.Sprintf("%s/manga?title=%s&limit=20&includes[]=author&includes[]=cover_art",
		s.baseURL, url.QueryEscape(title))

	resp, err := s.client.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search mangadex: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mangadex API error (status %d): %s", resp.StatusCode, string(body))
	}

	var apiResp models.MangadexAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.convertAPIResponseToMangas(apiResp.Data), nil
}

// GetMangaByID retrieves a specific manga by its MangaDex ID
func (s *MangadexService) GetMangaByID(mangadexID string) (*models.MangadexManga, error) {
	apiURL := fmt.Sprintf("%s/manga/%s?includes[]=author&includes[]=cover_art", s.baseURL, mangadexID)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get manga: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mangadex API error (status %d): %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Data models.MangadexAPIData `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	mangas := s.convertAPIResponseToMangas([]models.MangadexAPIData{apiResp.Data})
	if len(mangas) == 0 {
		return nil, fmt.Errorf("no manga found")
	}

	return &mangas[0], nil
}

// convertAPIResponseToMangas converts raw API data to MangadexManga models
func (s *MangadexService) convertAPIResponseToMangas(apiData []models.MangadexAPIData) []models.MangadexManga {
	var mangas []models.MangadexManga

	for _, item := range apiData {
		manga := models.MangadexManga{
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

		// Get author name from relationships
		for _, rel := range item.Relationships {
			if rel.Type == "author" && rel.Attributes != nil {
				manga.Author = rel.Attributes.Name
				break
			}
		}

		// Get cover art URL from relationships
		var coverFileName string
		for _, rel := range item.Relationships {
			if rel.Type == "cover_art" && rel.Attributes != nil {
				coverFileName = rel.Attributes.FileName
				break
			}
		}

		// Construct cover URL if we have the filename
		if coverFileName != "" {
			manga.CoverURL = fmt.Sprintf("https://uploads.mangadex.org/covers/%s/%s.256.jpg", item.ID, coverFileName)
		} else {
			// Fallback to manga page URL
			manga.CoverURL = fmt.Sprintf("https://mangadex.org/title/%s", item.ID)
		}

		mangas = append(mangas, manga)
	}

	return mangas
}
