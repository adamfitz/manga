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
	var searchResults []models.Manga

	// Table selection for import
	tableSelect := widget.NewSelect([]string{"manga", "mangadex"}, nil)
	tableSelect.SetSelected("manga")

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
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(searchResults) {
				box := obj.(*fyne.Container)
				box.Objects[0].(*widget.Label).SetText(searchResults[id].Title)
				box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("Author: %s", searchResults[id].Author))
			}
		},
	)

	detailContainer := container.NewVBox(
		widget.NewLabel("Select a manga to view details"),
	)

	resultsList.OnSelected = func(id widget.ListItemID) {
		if id < len(searchResults) {
			detailContainer.Objects = a.createSearchResultDetail(&searchResults[id])
			detailContainer.Refresh()
		}
	}

	searchBox := container.NewBorder(nil, nil, tableSelect, searchButton, searchEntry)
	resultsContainer := container.NewHSplit(resultsList, detailContainer)

	return container.NewBorder(searchBox, nil, nil, nil, resultsContainer)
}
