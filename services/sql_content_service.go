package services

import (
	"fmt"
	"mangadb/models"
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
	table   string
	columns string
	name    string
	altName string
}

var contentColumnMaps = map[models.ContentType]contentColumnMap{
	models.TypeManga: {
		table: "manga",
		columns: `
            id,
            COALESCE(title, '')       AS title,
            COALESCE(alt_title, '')   AS alt_title,
            COALESCE(author, '')      AS author,
            COALESCE(description, '') AS description,
            COALESCE(cover_url, '')   AS cover_url,
            COALESCE(url, '')         AS url,
            COALESCE(mangadex_id, '') AS mangadex_id,
            COALESCE(status, '')      AS status,
            created_at,
            updated_at`,
		name:    "title",
		altName: "alt_title",
	},

	models.TypeLightNovel: {
		table: "lightnovel",
		columns: `
            id,
            COALESCE(name, '')      AS name,
            COALESCE(alt_name, '')  AS alt_name,
            COALESCE(url, '')       AS url,
            volumes,
            COALESCE(status, '')    AS status`,
		name:    "name",
		altName: "alt_name",
	},

	models.TypeAnime: {
		table: "anime",
		columns: `
            id,
            COALESCE(name, '')      AS name,
            COALESCE(alt_name, '')  AS alt_name,
            COALESCE(url, '')       AS url,
            COALESCE(status, '')    AS status`,
		name:    "name",
		altName: "alt_name",
	},

	models.TypeWebNovel: {
		table: "webnovel",
		columns: `
            id,
            COALESCE(name, '')      AS name,
            COALESCE(alt_name, '')  AS alt_name,
            COALESCE(url, '')       AS url,
            COALESCE(status, '')    AS status`,
		name:    "name",
		altName: "alt_name",
	},

	models.TypeWebtoons: {
		table: "webtoons",
		columns: `
            id,
            COALESCE(name, '')      AS name,
            COALESCE(alt_name, '')  AS alt_name,
            COALESCE(url, '')       AS url,
            COALESCE(status, '')    AS status`,
		name:    "name",
		altName: "alt_name",
	},
}

// normalizedSelect builds the SELECT clause and returns table + column info
func normalizedSelect(contentType models.ContentType) (string, contentColumnMap, error) {
	m, ok := contentColumnMaps[contentType]
	if !ok {
		return "", contentColumnMap{}, fmt.Errorf("unsupported content type: %v", contentType)
	}

	return m.columns, m, nil
}
