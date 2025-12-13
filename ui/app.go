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
	contentService  *services.ContentService
	mangadexService *services.MangadexService
	mangaService    *services.MangaService    // added
	bookmarkService *services.BookmarkService // added
}

func NewApp(db *sql.DB, cfg *config.Config) *App {
	a := &App{
		fyneApp:         app.New(),
		db:              db,
		config:          cfg,
		contentService:  services.NewContentService(db),
		mangadexService: services.NewMangadexService(), // adjust constructor if needed
		mangaService:    services.NewMangaService(db),
		bookmarkService: services.NewBookmarkService(db),
	}

	if db != nil {
		a.contentService = services.NewContentService(db)
	}

	a.mainWindow = a.fyneApp.NewWindow("Content Database Manager")
	a.mainWindow.Resize(fyne.NewSize(1400, 900))

	a.setupMainWindow()

	return a
}

func (a *App) SetDatabase(db *sql.DB) {
	a.db = db
	a.contentService = services.NewContentService(db)

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
