package ui

import (
	"fmt"
	"manga/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// showAddContentDialog shows a dialog to add a new entry
func (a *App) showAddContentDialog(contentType models.ContentType, onSuccess func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter name/title...")

	altNameEntry := widget.NewEntry()
	altNameEntry.SetPlaceHolder("Enter alternative name...")

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://...")

	statusEntry := widget.NewEntry()
	statusEntry.SetPlaceHolder("ongoing, completed, etc.")
	statusEntry.SetText("unknown")

	var formItems []*widget.FormItem

	if contentType == models.TypeManga {
		authorEntry := widget.NewEntry()
		authorEntry.SetPlaceHolder("Author name...")

		descEntry := widget.NewMultiLineEntry()
		descEntry.SetPlaceHolder("Description...")
		descEntry.SetMinRowsVisible(3)

		coverURLEntry := widget.NewEntry()
		coverURLEntry.SetPlaceHolder("https://...")

		mangadexIDEntry := widget.NewEntry()
		mangadexIDEntry.SetPlaceHolder("MangaDex ID (optional)")

		formItems = []*widget.FormItem{
			{Text: "Name/Title", Widget: nameEntry},
			{Text: "Alt Name", Widget: altNameEntry},
			{Text: "URL", Widget: urlEntry},
			{Text: "Author", Widget: authorEntry},
			{Text: "Description", Widget: descEntry},
			{Text: "Cover URL", Widget: coverURLEntry},
			{Text: "MangaDex ID", Widget: mangadexIDEntry},
			{Text: "Status", Widget: statusEntry},
		}

		form := &widget.Form{Items: formItems}

		d := dialog.NewCustomConfirm(
			"Add New "+string(contentType),
			"Add",
			"Cancel",
			form,
			func(confirmed bool) {
				if !confirmed {
					return
				}

				if nameEntry.Text == "" && altNameEntry.Text == "" {
					dialog.ShowError(
						fmt.Errorf("Name or Alt Name is required"),
						a.mainWindow,
					)
					return
				}

				content := &models.Content{
					Name:        nameEntry.Text,
					AltName:     altNameEntry.Text,
					URL:         urlEntry.Text,
					Author:      authorEntry.Text,
					Description: descEntry.Text,
					CoverURL:    coverURLEntry.Text,
					MangadexID:  mangadexIDEntry.Text,
					Status:      statusEntry.Text,
				}

				if err := a.contentService.Create(contentType, content); err != nil {
					dialog.ShowError(err, a.mainWindow)
					return
				}

				dialog.ShowInformation("Success", "Entry added successfully!", a.mainWindow)
				onSuccess()
			},
			a.mainWindow,
		)

		d.Resize(fyne.NewSize(600, 500))
		d.Show()

	} else {
		// For Anime, Light Novel, Webtoon, Web Novel - simpler form
		formItems = []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Alt Name", Widget: altNameEntry},
			{Text: "URL", Widget: urlEntry},
			{Text: "Status", Widget: statusEntry},
		}

		form := &widget.Form{Items: formItems}

		d := dialog.NewCustomConfirm(
			"Add New "+string(contentType),
			"Add",
			"Cancel",
			form,
			func(confirmed bool) {
				if !confirmed {
					return
				}

				if nameEntry.Text == "" && altNameEntry.Text == "" {
					dialog.ShowError(
						fmt.Errorf("Name or Alt Name is required"),
						a.mainWindow,
					)
					return
				}

				content := &models.Content{
					Name:    nameEntry.Text,
					AltName: altNameEntry.Text,
					URL:     urlEntry.Text,
					Status:  statusEntry.Text,
				}

				if err := a.contentService.Create(contentType, content); err != nil {
					dialog.ShowError(err, a.mainWindow)
					return
				}

				dialog.ShowInformation("Success", "Entry added successfully!", a.mainWindow)
				onSuccess()
			},
			a.mainWindow,
		)

		d.Resize(fyne.NewSize(500, 400))
		d.Show()
	}
}

// showEditContentDialog shows a dialog to edit an existing entry
func (a *App) showEditContentDialog(
	content *models.Content,
	contentType models.ContentType,
	onSuccess func(),
) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(content.Name)

	altNameEntry := widget.NewEntry()
	altNameEntry.SetText(content.AltName)

	urlEntry := widget.NewEntry()
	urlEntry.SetText(content.URL)

	statusEntry := widget.NewEntry()
	statusEntry.SetText(content.Status)

	var formItems []*widget.FormItem

	if contentType == models.TypeManga {
		authorEntry := widget.NewEntry()
		authorEntry.SetText(content.Author)

		descEntry := widget.NewMultiLineEntry()
		descEntry.SetText(content.Description)
		descEntry.SetMinRowsVisible(3)

		coverURLEntry := widget.NewEntry()
		coverURLEntry.SetText(content.CoverURL)

		mangadexIDEntry := widget.NewEntry()
		mangadexIDEntry.SetText(content.MangadexID)

		formItems = []*widget.FormItem{
			{Text: "Name/Title", Widget: nameEntry},
			{Text: "Alt Name", Widget: altNameEntry},
			{Text: "URL", Widget: urlEntry},
			{Text: "Author", Widget: authorEntry},
			{Text: "Description", Widget: descEntry},
			{Text: "Cover URL", Widget: coverURLEntry},
			{Text: "MangaDex ID", Widget: mangadexIDEntry},
			{Text: "Status", Widget: statusEntry},
		}

		form := &widget.Form{Items: formItems}

		d := dialog.NewCustomConfirm(
			"Edit "+string(contentType),
			"Save",
			"Cancel",
			form,
			func(confirmed bool) {
				if !confirmed {
					return
				}

				if nameEntry.Text == "" && altNameEntry.Text == "" {
					dialog.ShowError(
						fmt.Errorf("Name or Alt Name is required"),
						a.mainWindow,
					)
					return
				}

				content.Name = nameEntry.Text
				content.AltName = altNameEntry.Text
				content.URL = urlEntry.Text
				content.Author = authorEntry.Text
				content.Description = descEntry.Text
				content.CoverURL = coverURLEntry.Text
				content.MangadexID = mangadexIDEntry.Text
				content.Status = statusEntry.Text

				if err := a.contentService.Update(contentType, content); err != nil {
					dialog.ShowError(err, a.mainWindow)
					return
				}

				dialog.ShowInformation("Success", "Entry updated successfully!", a.mainWindow)
				onSuccess()
			},
			a.mainWindow,
		)

		d.Resize(fyne.NewSize(600, 500))
		d.Show()

	} else {
		formItems = []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Alt Name", Widget: altNameEntry},
			{Text: "URL", Widget: urlEntry},
			{Text: "Status", Widget: statusEntry},
		}

		form := &widget.Form{Items: formItems}

		d := dialog.NewCustomConfirm(
			"Edit "+string(contentType),
			"Save",
			"Cancel",
			form,
			func(confirmed bool) {
				if !confirmed {
					return
				}

				if nameEntry.Text == "" && altNameEntry.Text == "" {
					dialog.ShowError(
						fmt.Errorf("Name or Alt Name is required"),
						a.mainWindow,
					)
					return
				}

				content.Name = nameEntry.Text
				content.AltName = altNameEntry.Text
				content.URL = urlEntry.Text
				content.Status = statusEntry.Text

				if err := a.contentService.Update(contentType, content); err != nil {
					dialog.ShowError(err, a.mainWindow)
					return
				}

				dialog.ShowInformation("Success", "Entry updated successfully!", a.mainWindow)
				onSuccess()
			},
			a.mainWindow,
		)

		d.Resize(fyne.NewSize(500, 400))
		d.Show()
	}
}
