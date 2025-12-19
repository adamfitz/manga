package services

import (
	"fmt"
	"manga/models"
)

/*
This file is the SQL normalization layer.

It guarantees that ALL content types return the same 9 columns
in the same order, matching models.Content and rows.Scan().
*/

const normalizedContentSelect = `
	id,
	%s AS name,
	COALESCE(%s, '') AS alt_name,
	COALESCE(url, '') AS url,
	%s AS mangadex_id,
	%s AS author,
	%s AS description,
	%s AS cover_url,
	status
`

type contentColumnMap struct {
	table       string
	name        string
	altName     string
	mangadexID  string
	author      string
	description string
	coverURL    string
}

var contentColumnMaps = map[models.ContentType]contentColumnMap{
	models.TypeManga: {
		table:       "manga",
		name:        "title",
		altName:     "alt_title",
		mangadexID:  "COALESCE(mangadex_id, '')",
		author:      "COALESCE(author, '')",
		description: "COALESCE(description, '')",
		coverURL:    "COALESCE(cover_url, '')",
	},
	models.TypeAnime: {
		table:       "anime",
		name:        "name",
		altName:     "alt_name",
		mangadexID:  "''",
		author:      "''",
		description: "''",
		coverURL:    "''",
	},
	models.TypeLightNovel: {
		table:       "lightnovel",
		name:        "name",
		altName:     "alt_name",
		mangadexID:  "''",
		author:      "''",
		description: "''",
		coverURL:    "''",
	},
	models.TypeWebtoons: {
		table:       "webtoon",
		name:        "name",
		altName:     "alt_name",
		mangadexID:  "''",
		author:      "''",
		description: "''",
		coverURL:    "''",
	},
	models.TypeWebNovel: {
		table:       "webnovel",
		name:        "name",
		altName:     "alt_name",
		mangadexID:  "''",
		author:      "''",
		description: "''",
		coverURL:    "''",
	},
}

// normalizedSelect builds the SELECT clause and returns table + column info
func normalizedSelect(contentType models.ContentType) (string, contentColumnMap, error) {
	m, ok := contentColumnMaps[contentType]
	if !ok {
		return "", contentColumnMap{}, fmt.Errorf("unsupported content type: %v", contentType)
	}

	selectClause := fmt.Sprintf(
		normalizedContentSelect,
		m.name,
		m.altName,
		m.mangadexID,
		m.author,
		m.description,
		m.coverURL,
	)

	return selectClause, m, nil
}
