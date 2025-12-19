package ui

import (
	"fmt"
	"manga/models"
	"net/url"
	"sort"
	"strings"

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
	var previousResults []models.Content // Track what was loaded before search

	// ----------------------
	// Results list - ONLY SHOW NAME/TITLE
	// ----------------------
	resultsList := widget.NewList(
		func() int {
			return len(searchResults)
		},
		func() fyne.CanvasObject {
			nameLabel := widget.NewLabel("Template Name")
			nameLabel.Wrapping = fyne.TextWrapWord
			spacer := widget.NewLabel(" ")
			return container.NewVBox(nameLabel, spacer)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(searchResults) {
				vbox := obj.(*fyne.Container)
				nameLabel := vbox.Objects[0].(*widget.Label)

				displayName := searchResults[id].Name
				if displayName == "" {
					displayName = searchResults[id].AltName
				}
				nameLabel.SetText(displayName)
			}
		},
	)

	// ----------------------
	// Detail Container
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
			detailScroll.ScrollToTop()
		}
	}

	// Function to clear selection
	clearSelection := func() {
		resultsList.UnselectAll()
		detailContainer.Objects = []fyne.CanvasObject{
			widget.NewLabel("Select an entry to view details"),
		}
		detailContainer.Refresh()
	}

	// ----------------------
	// Pane headers
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
	// Wrap panes with headers
	// ----------------------
	resultsScroll := container.NewScroll(resultsList)

	leftPane := container.NewBorder(
		container.NewVBox(leftHeader, widget.NewSeparator()),
		nil, nil, nil,
		resultsScroll,
	)

	rightPane := container.NewBorder(
		container.NewVBox(rightHeader, widget.NewSeparator()),
		nil, nil, nil,
		detailScroll,
	)

	// ----------------------
	// Left / Right split
	// ----------------------
	resultsContainer := container.NewGridWithColumns(2, leftPane, rightPane)

	// ----------------------
	// Helper function to refresh results
	// ----------------------
	refreshResults := func() {
		results, err := a.contentService.GetAll(contentType)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}
		sort.Slice(results, func(i, j int) bool {
			nameI := results[i].Name
			if nameI == "" {
				nameI = results[i].AltName
			}
			nameJ := results[j].Name
			if nameJ == "" {
				nameJ = results[j].AltName
			}
			return strings.ToLower(nameI) < strings.ToLower(nameJ)
		})
		searchResults = results
		previousResults = nil
		clearSelection()
		resultsList.Refresh()
	}

	// ----------------------
	// Buttons and Actions
	// ----------------------
	performSearch := func() {
		if searchEntry.Text == "" {
			dialog.ShowInformation("Info", "Please enter a search term", a.mainWindow)
			return
		}
		if len(searchResults) > 0 && len(previousResults) == 0 {
			previousResults = searchResults
		}
		results, err := a.contentService.Search(contentType, searchEntry.Text)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}
		sort.Slice(results, func(i, j int) bool {
			nameI := results[i].Name
			if nameI == "" {
				nameI = results[i].AltName
			}
			nameJ := results[j].Name
			if nameJ == "" {
				nameJ = results[j].AltName
			}
			return strings.ToLower(nameI) < strings.ToLower(nameJ)
		})
		searchResults = results
		resultsList.Refresh()
		dialog.ShowInformation(
			"Results",
			fmt.Sprintf("Found %d entries", len(results)),
			a.mainWindow,
		)
	}

	performLookup := func() {
		if lookupEntry.Text == "" {
			dialog.ShowInformation("Info", "Please enter a lookup value", a.mainWindow)
			return
		}
		if len(searchResults) > 0 && len(previousResults) == 0 {
			previousResults = searchResults
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
	}

	searchEntry.OnSubmitted = func(s string) {
		performSearch()
	}

	lookupEntry.OnSubmitted = func(s string) {
		performLookup()
	}

	searchButton := widget.NewButton("Search", performSearch)
	searchButton.Alignment = widget.ButtonAlignCenter

	lookupButton := widget.NewButton("Lookup", performLookup)
	lookupButton.Alignment = widget.ButtonAlignCenter

	loadAllButton := widget.NewButton("Load All", func() {
		refreshResults()
		dialog.ShowInformation(
			"Success",
			fmt.Sprintf("Loaded %d entries", len(searchResults)),
			a.mainWindow,
		)
	})
	loadAllButton.Alignment = widget.ButtonAlignCenter

	clearButton := widget.NewButton("Clear", func() {
		searchEntry.SetText("")
		lookupEntry.SetText("")

		if len(previousResults) > 0 {
			searchResults = previousResults
			previousResults = nil
		} else {
			searchResults = []models.Content{}
		}

		clearSelection()
		resultsList.Refresh()
	})
	clearButton.Alignment = widget.ButtonAlignCenter

	// NEW: Add Entry Button
	addButton := widget.NewButton("Add New", func() {
		a.showAddContentDialog(contentType, refreshResults)
	})
	addButton.Alignment = widget.ButtonAlignCenter

	// NEW: Find Duplicates Button
	duplicatesButton := widget.NewButton("Find Duplicates", func() {
		a.findDuplicates(contentType, &searchResults, resultsList, clearSelection)
	})
	duplicatesButton.Alignment = widget.ButtonAlignCenter

	// ----------------------
	// Controls (top)
	// ----------------------
	searchButtonContainer := container.NewStack(searchButton)
	lookupButtonContainer := container.NewStack(lookupButton)

	searchBox := container.NewBorder(nil, nil, nil, searchButtonContainer, searchEntry)
	lookupBox := container.NewBorder(nil, nil, nil, lookupButtonContainer, lookupEntry)

	// Updated button row with 4 buttons
	buttonRow1 := container.NewGridWithColumns(2, loadAllButton, clearButton)
	buttonRow2 := container.NewGridWithColumns(2, addButton, duplicatesButton)

	controlsBox := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s Database", contentType)),
		widget.NewSeparator(),
		searchBox,
		lookupBox,
		buttonRow1,
		buttonRow2,
		widget.NewSeparator(),
	)

	// ----------------------
	// Final layout
	// ----------------------
	content := container.NewBorder(
		controlsBox,
		nil, nil, nil,
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

	// Helper to create selectable text
	newSelectableText := func(text string) *widget.Entry {
		entry := widget.NewEntry()
		entry.SetText(text)
		entry.Wrapping = fyne.TextWrapWord
		entry.OnChanged = func(s string) {
			if s != text {
				entry.SetText(text)
			}
		}
		return entry
	}

	// Helper to create hyperlink
	newHyperlink := func(label, urlStr string) fyne.CanvasObject {
		if urlStr == "" {
			return newSelectableText(label + " ")
		}
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return newSelectableText(label + " " + urlStr)
		}

		link := widget.NewHyperlink(urlStr, parsedURL)
		link.Wrapping = fyne.TextWrapWord

		labelWidget := widget.NewLabel(label)
		return container.NewVBox(labelWidget, link)
	}

	objects := []fyne.CanvasObject{
		newSelectableText(fmt.Sprintf("ID: %d", content.ID)),
		newSelectableText(fmt.Sprintf("Name: %s", content.Name)),
		newSelectableText(fmt.Sprintf("Alt Name: %s", content.AltName)),
		newHyperlink("URL:", content.URL),
	}

	if contentType == models.TypeManga {
		objects = append(objects,
			newSelectableText(fmt.Sprintf("Author: %s", content.Author)),
			newSelectableText(fmt.Sprintf("Description: %s", content.Description)),
			newHyperlink("Cover URL:", content.CoverURL),
		)
	}

	if contentType.HasMangadexID() && content.MangadexID != "" {
		objects = append(objects,
			newSelectableText(fmt.Sprintf("MangaDex ID: %s", content.MangadexID)),
		)
	}

	objects = append(objects,
		newSelectableText(fmt.Sprintf("Status: %s", content.Status)),
		widget.NewSeparator(),
	)

	// Action buttons
	editBtn := widget.NewButton("Edit Entry", func() {
		a.showEditContentDialog(content, contentType, func() {
			// Refresh the list
			results, err := a.contentService.GetAll(contentType)
			if err != nil {
				dialog.ShowError(err, a.mainWindow)
				return
			}
			*searchResults = results
			list.Refresh()
			// Refresh the detail view
			for i, c := range *searchResults {
				if c.ID == content.ID {
					list.Select(i)
					break
				}
			}
		})
	})

	deleteBtn := widget.NewButton("Delete Entry", func() {
		dialog.ShowConfirm(
			"Confirm Delete",
			fmt.Sprintf("Are you sure you want to delete '%s'?", content.Name),
			func(confirmed bool) {
				if !confirmed {
					return
				}

				err := a.contentService.Delete(contentType, content.ID)
				if err != nil {
					dialog.ShowError(err, a.mainWindow)
					return
				}

				dialog.ShowInformation("Success", "Entry deleted successfully", a.mainWindow)

				// Remove from searchResults and refresh
				newResults := make([]models.Content, 0)
				for _, c := range *searchResults {
					if c.ID != content.ID {
						newResults = append(newResults, c)
					}
				}
				*searchResults = newResults
				list.Refresh()
				list.UnselectAll()
			},
			a.mainWindow,
		)
	})

	buttonRow := container.NewGridWithColumns(2, editBtn, deleteBtn)
	objects = append(objects, buttonRow)

	return objects
}

// findDuplicates searches for duplicate entries
func (a *App) findDuplicates(
	contentType models.ContentType,
	searchResults *[]models.Content,
	list *widget.List,
	clearSelection func(),
) {
	all, err := a.contentService.GetAll(contentType)
	if err != nil {
		dialog.ShowError(err, a.mainWindow)
		return
	}

	// Map to track duplicates
	seen := make(map[string][]models.Content)

	for _, c := range all {
		// Create a key based on name, alt_name, and URL
		key := fmt.Sprintf("%s|%s|%s",
			strings.ToLower(strings.TrimSpace(c.Name)),
			strings.ToLower(strings.TrimSpace(c.AltName)),
			strings.ToLower(strings.TrimSpace(c.URL)),
		)

		seen[key] = append(seen[key], c)
	}

	// Collect duplicates
	var duplicates []models.Content
	for _, group := range seen {
		if len(group) > 1 {
			duplicates = append(duplicates, group...)
		}
	}

	if len(duplicates) == 0 {
		dialog.ShowInformation("No Duplicates", "No duplicate entries found!", a.mainWindow)
		return
	}

	// Sort duplicates
	sort.Slice(duplicates, func(i, j int) bool {
		nameI := duplicates[i].Name
		if nameI == "" {
			nameI = duplicates[i].AltName
		}
		nameJ := duplicates[j].Name
		if nameJ == "" {
			nameJ = duplicates[j].AltName
		}
		return strings.ToLower(nameI) < strings.ToLower(nameJ)
	})

	*searchResults = duplicates
	clearSelection()
	list.Refresh()

	dialog.ShowInformation(
		"Duplicates Found",
		fmt.Sprintf("Found %d duplicate entries", len(duplicates)),
		a.mainWindow,
	)
}
