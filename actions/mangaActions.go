package actions

import (
	"fmt"
	"log"
	"main/auth"
	"main/postgresqldb"
)

// Return all entries in name column
func MangaNames() ([]string, error) {
	//load db connection config
	config, _ := auth.LoadConfig()

	// Connect to postgresql db
	pgDb, err := postgresqldb.OpenDatabase(
		config.PgServer,
		config.PgPort,
		config.PgUser,
		config.PgPassword,
		config.PgDbName)
	if err != nil {
		log.Fatalf("MangaNames() - Error opening database: %v", err)
	}
	defer pgDb.Close()

	// all entry names
	mangaNames, err := postgresqldb.LookupColumnValues(pgDb, "manga", "name")
	if err != nil {
		log.Println("error retrieving \"manga\" table column n\"name\"", err)
		return nil, fmt.Errorf("error retrieving \"manga\" table column n\"name\" %v", err)
	}

	return mangaNames, nil
}
