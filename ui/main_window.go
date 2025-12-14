package ui

import (
	"manga/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *App) setupMainWindow() {
	// Create menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Database Settings", func() {
			a.showConfigDialog(false)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			a.fyneApp.Quit()
		}),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			a.showAboutDialog()
		}),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)
	a.mainWindow.SetMainMenu(mainMenu)

	// Create tabs
	if a.db == nil {
		// Show placeholder when no database connection
		content := container.NewCenter(
			container.NewVBox(
				widget.NewLabel("No Database Connection"),
				widget.NewLabel("Please configure database settings from File > Database Settings"),
			),
		)
		a.mainWindow.SetContent(content)
	} else {
		tabs := container.NewAppTabs(
			// Use content_view style for Manga
			container.NewTabItem("Manga", a.createContentView(models.TypeManga)),
			container.NewTabItem("Bookmarks", a.createBookmarkView()),
			container.NewTabItem("MangaDex Search", a.createMangadexSearchView()),
			container.NewTabItem("Anime", a.createContentView(models.TypeAnime)),
			container.NewTabItem("Light Novel", a.createContentView(models.TypeLightNovel)),
		)
		a.mainWindow.SetContent(tabs)
	}
}
