package main

import (
	//"encoding/json"
	"fmt"
	"log"
	"main/auth"
	"main/mangadex"
	"strings"
	//"main/compare"
	//"main/mangadex"
	"flag"
	"main/actions"
	"main/parser"
	"main/postgresqldb"
	"main/webfrontend"
	"os"
	"path/filepath"
)

func init() {
	logFile, err := os.OpenFile("manga.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

func main() {

	startWeb := flag.Bool("w", false, "Start web server")
	flag.Parse()

	if *startWeb {
		webfrontend.StartServer("8080")
	} else {
		// placeholder for manga name cmoparisons
		DbNameCompare()
	}
}

func PgQueryByID(id string) {
	/*
		Dumps the postgresql db.
	*/

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
		log.Fatalf("Error opening database: %v", err)
	}
	defer pgDb.Close()

	// get all data in postgresql manga table
	data, err := postgresqldb.LookupByID(pgDb, "manga", id)
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
	}

	for key, value := range data {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func DownloadChapters(mangaName, mangadexId string) {
	chapters, err := mangadex.ChaptersWithDetails(mangadexId)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range chapters {
		fmt.Printf("Chapter: %v | ID: %v\n", c["chapter"], c["id"])
		id := c["id"].(string)
		chapter := c["chapter"].(string)

		chapterPages, _ := mangadex.ChapterPages(id)

		baseUrl := chapterPages["baseUrl"].(string)
		chapterData := chapterPages["chapter"].(map[string]any)
		hash := chapterData["hash"].(string)
		pages := chapterData["data"].([]any)

		tempDir, _ := os.MkdirTemp("", "mangadex_pages_*")

		for _, p := range pages {
			page := p.(string)
			err := mangadex.DownloadPage(baseUrl, hash, page, tempDir)
			if err != nil {
				log.Println("Failed to download page:", page, err)
			}
		}

		cbzPath, err := mangadex.CreateCBZ(tempDir, mangaName, "Ch"+chapter)
		if err != nil {
			log.Println("Failed to create CBZ:", err)
		} else {
			fmt.Println("Saved:", cbzPath)
		}

		os.RemoveAll(tempDir) // Clean up
	}
}

func ListManagdexMangaStatus(status string) {
	// Load the configuration
	config, err := auth.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to the PostgreSQL database
	pgDb, err := postgresqldb.OpenDatabase(
		config.PgServer,
		config.PgPort,
		config.PgUser,
		config.PgPassword,
		config.PgDbName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer pgDb.Close()

	// Query for completed manga
	completedManga, err := postgresqldb.LookupByStatus(pgDb, "mangadex", status)
	if err != nil {
		log.Fatalf("Error querying completed manga: %v", err)
	}

	for _, manga := range completedManga {
		fmt.Println(manga["name"])
	}
}

func copyDirs(status, srcDir, destDir string) {

	// Load the configuration
	config, err := auth.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to the PostgreSQL database
	pgDb, err := postgresqldb.OpenDatabase(
		config.PgServer,
		config.PgPort,
		config.PgUser,
		config.PgPassword,
		config.PgDbName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer pgDb.Close()

	// Query for completed manga
	completedManga, err := postgresqldb.LookupByStatus(pgDb, "mangadex", status)
	if err != nil {
		log.Fatalf("Error querying completed manga: %v", err)
	}

	// Copy directories from source to destination
	for _, dirEntry := range completedManga {
		// typecast to string
		name, ok := dirEntry["name"].(string)
		if !ok {
			log.Fatalf("Expected string for 'name', but got: %T", dirEntry["name"])
		}

		name = strings.TrimSpace(name)
		srcPath := filepath.Join(srcDir, name)
		dstPath := filepath.Join(destDir, name)

		// Check if destination exists and is up-to-date
		same, err := parser.DirsAreEqual(srcPath, dstPath)
		if err != nil {
			log.Printf("Error comparing directories %s and %s: %v", srcPath, dstPath, err)
		}
		if same {
			fmt.Printf("Skipping copy of %s; already exists and is identical\n", name)
			continue
		}

		if _, err := os.Stat(srcPath); err != nil {
			log.Printf("Directory not found: %s (from db: %q)", srcPath, name)
			continue
		}

		err = parser.CopyDir(srcPath, dstPath)
		if err != nil {
			log.Fatalf("Error copying directory %s:  %v", name, err)
		} else {
			fmt.Printf("Directory copied from %s to %s\n", srcPath, dstPath)
		}
	}
	// List the directories in the destination
	actions.DirList(destDir)
}

func DbNameCompare() {
	// list of all directories

	dirList, _ := actions.DirList("/mnt/manga/")

	//list of mangadex table entries
	mangadexList, _ := actions.MangadexNames()

	// list of manga table entries
	mangaList, _ := actions.MangaNames()

	// merge list of names frmo both tables
	allNames := parser.NormalizeAndDeduplicate(mangadexList, mangaList)

	//for _, name := range allNames {
	//	fmt.Println(name)
	//}

	// compare them both

	// manga names that ONLY appear in the directory list (not in database)
	dirCompare := parser.FindUniqueStrings(dirList, allNames)
	parser.WriteMissingTableEntriesWithSourceTags("MissingFromDB.txt", dirCompare, mangadexList, mangaList)
	fmt.Println("Names found in the driectory list BUT missing from the database:")
	fmt.Println("Output file written to: ./MissingFromDB.txt")

	// names that are not present in the directory list but ARE in the database tables
	tableNameCompare := parser.FindUniqueStrings(allNames, dirList)

	// save dir list to file and which table it appears in
	parser.WriteMissingDirsWithSourceTags("MissingDirs.txt", tableNameCompare, mangadexList, mangaList)
	fmt.Println("Names found in the database and missing from directory / disk:")
	fmt.Println("Output file written to: ./MissingDirs.txt")

	// name could be different - spelling
	// cant detect this just in code, it would be too much of a PITA

}
