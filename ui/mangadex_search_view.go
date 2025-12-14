package ui

import (
	"fmt"
	"manga/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) createMangadexSearchView() *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search for manga on MangaDex API...")

	var resultsList *widget.List
	var searchResults []models.MangadexManga

	searchButton := widget.NewButton("Search MangaDex", func() {
		if searchEntry.Text == "" {
			dialog.ShowInformation("Info", "Please enter a search term", a.mainWindow)
			return
		}

		results, err := a.mangadexService.Search(searchEntry.Text)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}

		searchResults = results
		resultsList.Refresh()
		dialog.ShowInformation("Results", fmt.Sprintf("Found %d manga on MangaDex", len(results)), a.mainWindow)
	})

	resultsList = widget.NewList(
		func() int {
			return len(searchResults)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Template Title"),
				widget.NewLabel("Template Author"),
				widget.NewLabel("Template Status"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(searchResults) {
				box := obj.(*fyne.Container)
				box.Objects[0].(*widget.Label).SetText(searchResults[id].Title)
				box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("Author: %s", searchResults[id].Author))
				box.Objects[2].(*widget.Label).SetText(fmt.Sprintf("Status: %s", searchResults[id].Status))
			}
		},
	)

	detailContainer := container.NewVBox(
		widget.NewLabel("Select a manga to view details"),
	)

	resultsList.OnSelected = func(id widget.ListItemID) {
		if id < len(searchResults) {
			detailContainer.Objects = a.createMangadexSearchResultDetail(&searchResults[id])
			detailContainer.Refresh()
		}
	}

	searchBox := container.NewBorder(nil, nil, nil, searchButton, searchEntry)
	resultsContainer := container.NewHSplit(resultsList, detailContainer)

	return container.NewBorder(searchBox, nil, nil, nil, resultsContainer)
}

// createMangadexSearchResultDetail creates the detail view for a selected manga
func (a *App) createMangadexSearchResultDetail(manga *models.MangadexManga) []fyne.CanvasObject {
	titleLabel := widget.NewLabel(manga.Title)
	titleLabel.Wrapping = fyne.TextWrapWord
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	authorLabel := widget.NewLabel(fmt.Sprintf("Author: %s", manga.Author))
	statusLabel := widget.NewLabel(fmt.Sprintf("Status: %s", manga.Status))
	mangadexIDLabel := widget.NewLabel(fmt.Sprintf("MangaDex ID: %s", manga.MangadexID))

	descriptionLabel := widget.NewLabel(manga.Description)
	descriptionLabel.Wrapping = fyne.TextWrapWord

	// Button to add to library
	addToLibraryBtn := widget.NewButton("Add to Library", func() {
		// Check if already in library
		existing, err := a.mangaService.GetByMangadexID(manga.MangadexID)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to check library: %w", err), a.mainWindow)
			return
		}

		if existing != nil {
			dialog.ShowInformation("Already in Library",
				fmt.Sprintf("'%s' is already in your library.", manga.Title),
				a.mainWindow)
			return
		}

		// Convert MangadexManga to Manga model
		newManga := &models.Manga{
			Title:       manga.Title,
			Author:      manga.Author,
			Description: manga.Description,
			CoverURL:    manga.CoverURL,
			MangadexID:  manga.MangadexID,
			Status:      manga.Status,
		}

		if err := a.mangaService.Create(newManga); err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}

		dialog.ShowInformation("Success",
			fmt.Sprintf("'%s' added to your library!", manga.Title),
			a.mainWindow)
	})

	// Button to view on MangaDex
	viewOnMangadexBtn := widget.NewButton("View on MangaDex", func() {
		// Open in browser
		dialog.ShowInformation("MangaDex Link",
			fmt.Sprintf("https://mangadex.org/title/%s", manga.MangadexID),
			a.mainWindow)
	})

	return []fyne.CanvasObject{
		titleLabel,
		authorLabel,
		statusLabel,
		mangadexIDLabel,
		widget.NewSeparator(),
		widget.NewLabel("Description:"),
		descriptionLabel,
		widget.NewSeparator(),
		container.NewHBox(addToLibraryBtn, viewOnMangadexBtn),
	}
}
