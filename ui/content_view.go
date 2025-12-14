package ui

import (
	"fmt"
	"manga/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) createContentView(contentType models.ContentType) *fyne.Container {
	// Search section
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Substring search in name or alt_name...")

	var resultsList *widget.List
	var searchResults []models.Content

	searchButton := widget.NewButton("Search", func() {
		if searchEntry.Text == "" {
			dialog.ShowInformation("Info", "Please enter a search term", a.mainWindow)
			return
		}

		results, err := a.contentService.Search(contentType, searchEntry.Text)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}

		searchResults = results
		resultsList.Refresh()
		dialog.ShowInformation("Results", fmt.Sprintf("Found %d entries", len(results)), a.mainWindow)
	})

	// Lookup section
	lookupEntry := widget.NewEntry()
	lookupEntry.SetPlaceHolder("Exact match: ID, name, or alt_name...")

	lookupButton := widget.NewButton("Lookup", func() {
		if lookupEntry.Text == "" {
			dialog.ShowInformation("Info", "Please enter a lookup value", a.mainWindow)
			return
		}

		result, err := a.contentService.Lookup(contentType, lookupEntry.Text)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}

		searchResults = []models.Content{*result}
		resultsList.Refresh()
	})

	// Load All button
	loadAllButton := widget.NewButton("Load All", func() {
		results, err := a.contentService.GetAll(contentType)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}

		searchResults = results
		resultsList.Refresh()
		dialog.ShowInformation("Success", fmt.Sprintf("Loaded %d entries", len(results)), a.mainWindow)
	})

	// Results list
	resultsList = widget.NewList(
		func() int {
			return len(searchResults)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Template Name"),
				widget.NewLabel("Template Alt Name"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(searchResults) {
				box := obj.(*fyne.Container)
				box.Objects[0].(*widget.Label).SetText(searchResults[id].Name)
				altName := searchResults[id].AltName
				if altName == "" {
					altName = "(no alt name)"
				}
				box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("Alt: %s", altName))
			}
		},
	)

	detailContainer := container.NewVBox(
		widget.NewLabel("Select an entry to view details"),
	)

	resultsList.OnSelected = func(id widget.ListItemID) {
		if id < len(searchResults) {
			detailContainer.Objects = a.createContentDetail(&searchResults[id], contentType, resultsList, &searchResults)
			detailContainer.Refresh()
		}
	}

	// Layout
	searchBox := container.NewBorder(nil, nil, nil, searchButton, searchEntry)
	lookupBox := container.NewBorder(nil, nil, nil, lookupButton, lookupEntry)
	controlsBox := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s Database", contentType)),
		widget.NewSeparator(),
		searchBox,
		lookupBox,
		loadAllButton,
		widget.NewSeparator(),
	)

	resultsContainer := container.NewHSplit(resultsList, detailContainer)
	resultsContainer.SetOffset(0.3)

	return container.NewBorder(controlsBox, nil, nil, nil, resultsContainer)
}

func (a *App) createContentDetail(content *models.Content, contentType models.ContentType, list *widget.List, searchResults *[]models.Content) []fyne.CanvasObject {
	objects := []fyne.CanvasObject{
		widget.NewLabel(fmt.Sprintf("ID: %d", content.ID)),
		widget.NewLabel(fmt.Sprintf("Name: %s", content.Name)),
		widget.NewLabel(fmt.Sprintf("Alt Name: %s", func() string {
			if content.AltName == "" {
				return "(none)"
			}
			return content.AltName
		}())),
		widget.NewLabel(fmt.Sprintf("URL: %s", content.URL)),
	}

	// Add MangaDex ID if applicable
	if contentType.HasMangadexID() && content.MangadexID != "" {
		objects = append(objects, widget.NewLabel(fmt.Sprintf("MangaDex ID: %s", content.MangadexID)))
	}

	// Add manga-specific fields
	if contentType == models.TypeManga {
		if content.Author != "" {
			objects = append(objects, widget.NewLabel(fmt.Sprintf("Author: %s", content.Author)))
		}
		if content.Description != "" {
			desc := widget.NewLabel(fmt.Sprintf("Description: %s", content.Description))
			desc.Wrapping = fyne.TextWrapWord
			objects = append(objects, desc)
		}
		if content.CoverURL != "" {
			objects = append(objects, widget.NewLabel(fmt.Sprintf("Cover URL: %s", content.CoverURL)))
		}
	}

	objects = append(objects, widget.NewLabel(fmt.Sprintf("Status: %s", content.Status)))
	objects = append(objects, widget.NewSeparator())

	// Delete button
	deleteBtn := widget.NewButton("Delete Entry", func() {
		dialog.ShowConfirm("Delete Entry",
			fmt.Sprintf("Are you sure you want to delete '%s'?", content.Name),
			func(ok bool) {
				if !ok {
					return
				}

				if err := a.contentService.Delete(contentType, content.ID); err != nil {
					dialog.ShowError(err, a.mainWindow)
					return
				}

				// Remove from searchResults
				for i, c := range *searchResults {
					if c.ID == content.ID {
						*searchResults = append((*searchResults)[:i], (*searchResults)[i+1:]...)
						break
					}
				}

				dialog.ShowInformation("Success", "Entry deleted", a.mainWindow)
				list.Refresh()
				list.UnselectAll()
			}, a.mainWindow)
	})

	objects = append(objects, deleteBtn)

	return objects
}
