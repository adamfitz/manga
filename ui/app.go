package ui

import (
	"database/sql"
	"manga/config"
	"manga/services"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type App struct {
	fyneApp         fyne.App
	mainWindow      fyne.Window
	db              *sql.DB
	config          *config.Config
	mangaService    *services.MangaService
	bookmarkService *services.BookmarkService
	mangadexService *services.MangadexService
}

func NewApp(db *sql.DB, cfg *config.Config) *App {
	a := &App{
		fyneApp:         app.New(),
		db:              db,
		config:          cfg,
		mangadexService: services.NewMangadexService(),
	}

	if db != nil {
		a.mangaService = services.NewMangaService(db)
		a.bookmarkService = services.NewBookmarkService(db)
	}

	a.mainWindow = a.fyneApp.NewWindow("Manga Tracker")
	a.mainWindow.Resize(fyne.NewSize(1200, 800))

	a.setupMainWindow()

	return a
}

func (a *App) SetDatabase(db *sql.DB) {
	a.db = db
	a.mangaService = services.NewMangaService(db)
	a.bookmarkService = services.NewBookmarkService(db)

	// Refresh the UI
	a.setupMainWindow()
}

func (a *App) Run() {
	// Check if database is connected, if not show config dialog
	if a.db == nil {
		a.showConfigDialog(true)
	}

	a.mainWindow.ShowAndRun()
}
