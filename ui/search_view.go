package ui

import (
	"fmt"
	"manga/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) createSearchView() *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search for manga on MangaDex...")

	var resultsList *widget.List
	var searchResults []models.Manga

	searchButton := widget.NewButton("Search", func() {
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

	searchBox := container.NewBorder(nil, nil, nil, searchButton, searchEntry)
	resultsContainer := container.NewHSplit(resultsList, detailContainer)

	return container.NewBorder(searchBox, nil, nil, nil, resultsContainer)
}

func (a *App) createSearchResultDetail(manga *models.Manga) []fyne.CanvasObject {
	title := widget.NewLabel(fmt.Sprintf("Title: %s", manga.Title))
	title.Wrapping = fyne.TextWrapWord

	author := widget.NewLabel(fmt.Sprintf("Author: %s", manga.Author))
	status := widget.NewLabel(fmt.Sprintf("Status: %s", manga.Status))

	description := widget.NewLabel(fmt.Sprintf("Description:\n%s", manga.Description))
	description.Wrapping = fyne.TextWrapWord

	addButton := widget.NewButton("Add to Library", func() {
		if err := a.mangaService.Create(manga); err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}
		dialog.ShowInformation("Success", "Manga added to your library", a.mainWindow)
	})

	return []fyne.CanvasObject{
		title,
		author,
		status,
		widget.NewSeparator(),
		description,
		widget.NewSeparator(),
		addButton,
	}
}
