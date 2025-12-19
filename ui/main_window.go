package ui

import (
	"fmt"
	"manga/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func (a *App) setupMainWindow() {
	// Create menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Database Settings", func() {
			a.showConfigDialog(false)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Compact Database", func() {
			a.compactDatabase()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			a.FyneApp.Quit()
		}),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			a.showAboutDialog()
		}),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)
	a.mainWindow.SetMainMenu(mainMenu)

	// Add Ctrl+Q keyboard shortcut to quit
	ctrlQ := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}
	a.mainWindow.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
		a.FyneApp.Quit()
	})

	if a.db == nil {
		// Placeholder when no DB
		content := container.NewCenter(
			container.NewVBox(
				widget.NewLabel("No Database Connection"),
				widget.NewLabel("Please configure database settings from File > Database Settings"),
			),
		)
		a.mainWindow.SetContent(content)
	} else {
		tabs := container.NewAppTabs(
			container.NewTabItem("Manga", a.createContentView(models.TypeManga)),
			container.NewTabItem("Bookmarks", a.createBookmarkView()),
			container.NewTabItem("MangaDex Search", a.createMangadexSearchView()),
			container.NewTabItem("Anime", a.createContentView(models.TypeAnime)),
			container.NewTabItem("Light Novel", a.createContentView(models.TypeLightNovel)),
			container.NewTabItem("Webtoons", a.createContentView(models.TypeWebtoons)),
			container.NewTabItem("Web Novel", a.createContentView(models.TypeWebNovel)),
		)
		a.mainWindow.SetContent(tabs)
	}
}

// compactDatabase runs VACUUM on the database
func (a *App) compactDatabase() {
	if a.db == nil {
		dialog.ShowError(
			fmt.Errorf("No database connection"),
			a.mainWindow,
		)
		return
	}

	dialog.ShowConfirm(
		"Compact Database",
		"This will optimize and compact the database. Continue?",
		func(confirmed bool) {
			if !confirmed {
				return
			}

			// PostgreSQL uses VACUUM
			if _, err := a.db.Exec("VACUUM ANALYZE"); err != nil {
				dialog.ShowError(err, a.mainWindow)
				return
			}

			dialog.ShowInformation(
				"Success",
				"Database compacted and optimized successfully!",
				a.mainWindow,
			)
		},
		a.mainWindow,
	)
}
