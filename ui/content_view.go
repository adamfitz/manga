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

	// Add button
	addButton := widget.NewButton("Add New Entry", func() {
		a.showAddContentDialog(contentType, resultsList, &searchResults)
	})

	// Layout
	searchBox := container.NewBorder(nil, nil, nil, searchButton, searchEntry)
	lookupBox := container.NewBorder(nil, nil, nil, lookupButton, lookupEntry)
	controlsBox := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s Database", contentType)),
		widget.NewSeparator(),
		searchBox,
		lookupBox,
		container.NewHBox(loadAllButton, addButton),
		widget.NewSeparator(),
	)

	resultsContainer := container.NewHSplit(resultsList, detailContainer)
	resultsContainer.SetOffset(0.3)

	return container.NewBorder(controlsBox, nil, nil, nil, resultsContainer)
}

func (a *App) createContentDetail(content *models.Content, contentType models.ContentType, list *widget.List, searchResults *[]models.Content) []fyne.CanvasObject {
	idLabel := widget.NewLabel(fmt.Sprintf("ID: %d", content.ID))

	name := widget.NewLabel(fmt.Sprintf("Name: %s", content.Name))
	name.Wrapping = fyne.TextWrapWord

	altName := content.AltName
	if altName == "" {
		altName = "(none)"
	}
	altNameLabel := widget.NewLabel(fmt.Sprintf("Alt Name: %s", altName))
	altNameLabel.Wrapping = fyne.TextWrapWord

	url := widget.NewLabel(fmt.Sprintf("URL: %s", content.URL))
	url.Wrapping = fyne.TextWrapWord

	status := widget.NewLabel(fmt.Sprintf("Status: %s", content.Status))

	objects := []fyne.CanvasObject{
		idLabel,
		name,
		altNameLabel,
		url,
	}

	// Add mangadex_id if applicable
	if contentType.HasMangadexID() && content.MangadexID != "" {
		mangadexID := widget.NewLabel(fmt.Sprintf("MangaDex ID: %s", content.MangadexID))
		objects = append(objects, mangadexID)
	}

	objects = append(objects, status, widget.NewSeparator())

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

func (a *App) showAddContentDialog(contentType models.ContentType, list *widget.List, searchResults *[]models.Content) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name (required)")

	altNameEntry := widget.NewEntry()
	altNameEntry.SetPlaceHolder("Alt Name (optional)")

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("URL (required)")

	var mangadexIDEntry *widget.Entry
	var formItems []*widget.FormItem

	formItems = append(formItems,
		&widget.FormItem{Text: "Name", Widget: nameEntry},
		&widget.FormItem{Text: "Alt Name", Widget: altNameEntry},
		&widget.FormItem{Text: "URL", Widget: urlEntry},
	)

	// Add mangadex_id field if applicable
	if contentType.HasMangadexID() {
		mangadexIDEntry = widget.NewEntry()
		mangadexIDEntry.SetPlaceHolder("MangaDex ID (optional)")
		formItems = append(formItems, &widget.FormItem{Text: "MangaDex ID", Widget: mangadexIDEntry})
	}

	statusSelect := widget.NewSelect([]string{"ongoing", "completed", "hiatus", "cancelled"}, nil)
	statusSelect.SetSelected("ongoing")
	formItems = append(formItems, &widget.FormItem{Text: "Status", Widget: statusSelect})

	form := dialog.NewForm(fmt.Sprintf("Add New %s", contentType), "Add", "Cancel",
		formItems,
		func(ok bool) {
			if !ok {
				return
			}

			if nameEntry.Text == "" || urlEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("name and URL are required"), a.mainWindow)
				return
			}

			content := &models.Content{
				Name:    nameEntry.Text,
				AltName: altNameEntry.Text,
				URL:     urlEntry.Text,
				Status:  statusSelect.Selected,
			}

			if contentType.HasMangadexID() && mangadexIDEntry != nil {
				content.MangadexID = mangadexIDEntry.Text
			}

			if err := a.contentService.Create(contentType, content); err != nil {
				dialog.ShowError(err, a.mainWindow)
				return
			}

			// Add to searchResults if we're viewing all
			*searchResults = append(*searchResults, *content)

			dialog.ShowInformation("Success", fmt.Sprintf("%s added with ID: %d", contentType, content.ID), a.mainWindow)
			list.Refresh()
		}, a.mainWindow)

	form.Resize(fyne.NewSize(500, 450))
	form.Show()
}
