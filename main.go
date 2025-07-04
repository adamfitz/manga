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
	"main/postgresqldb"
	//"main/sqlitedb"
	"flag"
	"main/actions"
	"main/parser"
	"main/scraper"
	"main/webfrontend"
	"os"
	"path/filepath"
	"sort"
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
		webfrontend.StartServer("8080")
		//var exclusionList = []string{"Completed"}
		//copyDirs("completed", "/mnt/manga", "/mnt/manga/completed")
		//MangaStatusAttributes()
		//NewMangaDbUpdate()
		//CheckIfBookmarkInDb()
		//CompareNames()
		//BlanketUpdateDb()
		//ExtractMangasWithoutChapterList()
		//UpdateMangasWithoutChapterList()
		//actions.DumpPostgressTable("manga", []string{"name", "url"})
		//PgQueryByID("21")
		//DownloadChapters("At First Glance, Shinoda-san Seems Cool but Is Actually Adorable!", "5187376e-3b32-4c8c-9fff-e95aca386463")
		//actions.GetDirList("/mnt/manga", exclusionList...)
		//ListManagdexMangaStatus("completed")
		//KunManga("ugly-complex")
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
	actions.GetDirList(destDir)
}

// mangaName is the name of the manga from the url eg:
// From: https://kunmanga.com/manga/ugly-complex/
// the mangaNmae will be the string "ugly-complex"
func KunManga(mangaName string) {
	{
		// Chapter holds URL and parsed chapter number
		type Chapter struct {
			URL    string
			Slug   string
			Number int
		}

		chapterURLs := scraper.KunMangaChapterUrls(mangaName)
		outputDir := "downloads"

		startChapter := 1
		endChapter := 13

		var filtered []Chapter

		// Filter and parse chapters on the fly
		for _, chapterUrl := range chapterURLs {
			parts := strings.Split(strings.Trim(chapterUrl, "/"), "/")
			chapterSlug := parts[len(parts)-1]

			chNum := scraper.ParseChapterNumber(chapterSlug)
			if chNum == 0 {
				log.Printf("[MAIN] Skipping invalid chapter slug: %s", chapterSlug)
				continue
			}

			if chNum < startChapter || chNum > endChapter {
				log.Printf("[MAIN] Skipping chapter %s (number %d): outside range %d-%d", chapterSlug, chNum, startChapter, endChapter)
				continue
			}

			filtered = append(filtered, Chapter{
				URL:    chapterUrl,
				Slug:   chapterSlug,
				Number: chNum,
			})
		}

		// Sort filtered chapters ascending by chapter number
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Number < filtered[j].Number
		})

		// Download each filtered chapter
		for _, ch := range filtered {
			//log.Printf("[MAIN] Starting download for chapter %s (number %d)", ch.Slug, ch.Number)
			fmt.Printf("Downloading chapter %d (%s)...\n", ch.Number, ch.Slug)
			err := scraper.DownloadKunMangaChapters(ch.URL, outputDir)
			if err != nil {
				log.Printf("[MAIN] Error downloading chapter %s: %v", ch.Slug, err)
			}
		}

	}
}
