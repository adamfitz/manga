package ui

import (
	"mangadb/config"
	"mangadb/database"
	"mangadb/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) showConfigDialog(required bool) {
	serverEntry := widget.NewEntry()
	serverEntry.SetText(a.config.DBServer)
	serverEntry.SetPlaceHolder("localhost")

	portEntry := widget.NewEntry()
	portEntry.SetText(a.config.DBPort)
	portEntry.SetPlaceHolder("5432")

	userEntry := widget.NewEntry()
	userEntry.SetText(a.config.DBUser)
	userEntry.SetPlaceHolder("postgres")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetText(a.config.DBPassword)
	passwordEntry.SetPlaceHolder("password")

	dbnameEntry := widget.NewEntry()
	dbnameEntry.SetText(a.config.DBName)
	dbnameEntry.SetPlaceHolder("manga_db")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Server", Widget: serverEntry, HintText: "Database server address"},
			{Text: "Port", Widget: portEntry, HintText: "Database port (default: 5432)"},
			{Text: "Username", Widget: userEntry, HintText: "Database username"},
			{Text: "Password", Widget: passwordEntry, HintText: "Database password"},
			{Text: "Database", Widget: dbnameEntry, HintText: "Database name"},
		},
	}

	var d *dialog.CustomDialog

	testButton := widget.NewButton("Test Connection", func() {
		testConfig := &config.Config{
			DBServer:   serverEntry.Text,
			DBPort:     portEntry.Text,
			DBUser:     userEntry.Text,
			DBPassword: passwordEntry.Text,
			DBName:     dbnameEntry.Text,
		}

		db, err := database.InitDB(testConfig)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}
		db.Close()

		dialog.ShowInformation("Success", "Connection successful!", a.mainWindow)
	})

	saveButton := widget.NewButton("Save & Connect", func() {
		a.config.DBServer = serverEntry.Text
		a.config.DBPort = portEntry.Text
		a.config.DBUser = userEntry.Text
		a.config.DBPassword = passwordEntry.Text
		a.config.DBName = dbnameEntry.Text

		if err := a.config.Save(); err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}

		db, err := database.InitDB(a.config)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}

		if err := database.RunMigrations(db); err != nil {
			dialog.ShowError(err, a.mainWindow)
			db.Close()
			return
		}

		if a.db != nil {
			a.db.Close()
		}

		a.SetDatabase(db)
		d.Hide()

		dialog.ShowInformation(
			"Success",
			"Database configuration saved and connected successfully!",
			a.mainWindow,
		)
	})

	cancelButton := widget.NewButton("Cancel", func() {
		if required && a.db == nil {
			dialog.ShowInformation(
				"Info",
				"Database connection is required to use the application.",
				a.mainWindow,
			)
			return
		}
		d.Hide()
	})

	buttons := container.NewHBox(testButton, saveButton, cancelButton)
	content := container.NewBorder(nil, buttons, nil, nil, form)

	d = dialog.NewCustom(
		"Database Configuration",
		"",
		content,
		a.mainWindow,
	)

	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}

func (a *App) showAboutDialog() {
	title := widget.NewLabel("MangaDb")
	title.TextStyle = fyne.TextStyle{Bold: true}

	version := widget.NewLabel("Version: " + utils.Version)

	description := widget.NewLabel(
		"A cross-platform anime, manga, lightnovel, webnovel and webtoon database front end tracking application",
	)
	description.Wrapping = fyne.TextWrapWord

	features := widget.NewLabel(
		"Features:\n" +
			"• Manage your libraries\n" +
			"• Bookmark your favorite series\n" +
			"• Search and import from MangaDex\n" +
			"• Cross-platform support",
	)
	features.Wrapping = fyne.TextWrapWord

	// Centered bold title
	centeredTitle := container.NewCenter(title)

	// centered version
	centeredVersion := container.NewCenter(version)

	// Declare window first so the close button can reference it
	var aboutWin fyne.Window
	closeBtn := widget.NewButton("Close", func() {
		aboutWin.Close()
	})

	// Main content (scrollable)
	mainContent := container.NewVBox(
		centeredTitle,
		centeredVersion,
		widget.NewSeparator(),
		description,
		widget.NewSeparator(),
		features,
	)

	scroll := container.NewScroll(mainContent)

	// Bottom area: separator + centered Close button
	bottom := container.NewVBox(
		widget.NewSeparator(),
		container.NewCenter(closeBtn),
	)

	// Border layout: scroll in center, close button at bottom
	content := container.NewBorder(nil, bottom, nil, nil, scroll)

	// Create and show window
	aboutWin = a.FyneApp.NewWindow("About MangaDb")
	aboutWin.SetContent(content)
	aboutWin.Resize(fyne.NewSize(400, 350))
	aboutWin.SetFixedSize(true)
	aboutWin.Show()
}
