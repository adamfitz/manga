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
	// ----------------------
	// Search Section
	// ----------------------
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Substring search in name or alt_name...")

	lookupEntry := widget.NewEntry()
	lookupEntry.SetPlaceHolder("Exact match: ID, name, or alt_name...")

	var searchResults []models.Content

	// ----------------------
	// Results list (UNCHANGED)
	// ----------------------
	resultsList := widget.NewList(
		func() int {
			return len(searchResults)
		},
		func() fyne.CanvasObject {
			nameLabel := widget.NewLabel("Template Name")
			nameLabel.Wrapping = fyne.TextWrapWord

			altLabel := widget.NewLabel("Template Alt Name")
			altLabel.Wrapping = fyne.TextWrapWord

			return container.NewVBox(nameLabel, altLabel)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(searchResults) {
				box := obj.(*fyne.Container)
				nameLabel := box.Objects[0].(*widget.Label)
				altLabel := box.Objects[1].(*widget.Label)

				nameLabel.SetText(searchResults[id].Name)

				altName := searchResults[id].AltName
				if altName == "" {
					altName = "(no alt name)"
				}
				altLabel.SetText(fmt.Sprintf("Alt: %s", altName))
			}
		},
	)

	// ----------------------
	// Detail Container (UNCHANGED)
	// ----------------------
	detailContainer := container.NewVBox(
		widget.NewLabel("Select an entry to view details"),
	)
	detailScroll := container.NewScroll(detailContainer)

	resultsList.OnSelected = func(id widget.ListItemID) {
		if id < len(searchResults) {
			detailContainer.Objects = a.createContentDetail(
				&searchResults[id],
				contentType,
				resultsList,
				&searchResults,
			)
			detailContainer.Refresh()
		}
	}

	// ----------------------
	// Pane headers (NEW)
	// ----------------------
	leftHeader := widget.NewLabelWithStyle(
		"Title",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	rightHeader := widget.NewLabelWithStyle(
		"Details",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	// ----------------------
	// Wrap panes with headers (NEW)
	// ----------------------
	resultsScroll := container.NewScroll(resultsList)

	leftPane := container.NewBorder(
		leftHeader,
		nil,
		nil,
		nil,
		resultsScroll,
	)

	rightPane := container.NewBorder(
		rightHeader,
		nil,
		nil,
		nil,
		detailScroll,
	)

	// ----------------------
	// Left / Right split (UNCHANGED STRUCTURE)
	// ----------------------
	resultsContainer := container.NewHSplit(leftPane, rightPane)
	resultsContainer.SetOffset(0.5)

	// ----------------------
	// Buttons and Actions
	// ----------------------
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
		dialog.ShowInformation(
			"Results",
			fmt.Sprintf("Found %d entries", len(results)),
			a.mainWindow,
		)
	})

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
		if result != nil {
			searchResults = []models.Content{*result}
		} else {
			searchResults = []models.Content{}
		}
		resultsList.Refresh()
	})

	loadAllButton := widget.NewButton("Load All", func() {
		results, err := a.contentService.GetAll(contentType)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}
		searchResults = results
		resultsList.Refresh()
		dialog.ShowInformation(
			"Success",
			fmt.Sprintf("Loaded %d entries", len(results)),
			a.mainWindow,
		)
	})

	// ----------------------
	// Controls (top)
	// ----------------------
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

	// ----------------------
	// Final layout
	// ----------------------
	content := container.NewBorder(
		controlsBox,
		nil,
		nil,
		nil,
		resultsContainer,
	)

	return content
}

func (a *App) createContentDetail(
	content *models.Content,
	contentType models.ContentType,
	list *widget.List,
	searchResults *[]models.Content,
) []fyne.CanvasObject {

	newWrappedLabel := func(text string) *widget.Label {
		l := widget.NewLabel(text)
		l.Wrapping = fyne.TextWrapWord
		return l
	}

	objects := []fyne.CanvasObject{
		newWrappedLabel(fmt.Sprintf("ID: %d", content.ID)),
		newWrappedLabel(fmt.Sprintf("Name: %s", content.Name)),
		newWrappedLabel(fmt.Sprintf("Alt Name: %s", content.AltName)),
		newWrappedLabel(fmt.Sprintf("URL: %s", content.URL)),
	}

	if contentType == models.TypeManga {
		objects = append(objects,
			newWrappedLabel(fmt.Sprintf("Author: %s", content.Author)),
			newWrappedLabel(fmt.Sprintf("Description: %s", content.Description)),
			newWrappedLabel(fmt.Sprintf("Cover URL: %s", content.CoverURL)),
		)
	}

	if contentType.HasMangadexID() && content.MangadexID != "" {
		objects = append(objects,
			newWrappedLabel(fmt.Sprintf("MangaDex ID: %s", content.MangadexID)),
		)
	}

	objects = append(objects,
		newWrappedLabel(fmt.Sprintf("Status: %s", content.Status)),
		widget.NewSeparator(),
	)

	deleteBtn := widget.NewButton("Delete Entry", func() {
		// unchanged
	})
	objects = append(objects, deleteBtn)

	return objects
}
