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
			// Add just a small spacer after the label
			spacer := widget.NewLabel(" ")
			return container.NewVBox(nameLabel, spacer)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(searchResults) {
				vbox := obj.(*fyne.Container)
				nameLabel := vbox.Objects[0].(*widget.Label)

				// Use name, fallback to alt_name if name is empty
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
			// Reset scroll position to top when selecting new entry
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
	// Pane headers - BOLD and styled
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
	// Wrap panes with headers in bordered containers
	// ----------------------
	resultsScroll := container.NewScroll(resultsList)

	leftPane := container.NewBorder(
		container.NewVBox(leftHeader, widget.NewSeparator()),
		nil,
		nil,
		nil,
		resultsScroll,
	)

	rightPane := container.NewBorder(
		container.NewVBox(rightHeader, widget.NewSeparator()),
		nil,
		nil,
		nil,
		detailScroll,
	)

	// ----------------------
	// Left / Right split - Use GridWithColumns for truly fixed 50/50
	// ----------------------
	resultsContainer := container.NewGridWithColumns(2, leftPane, rightPane)

	// ----------------------
	// Buttons and Actions
	// ----------------------
	performSearch := func() {
		if searchEntry.Text == "" {
			dialog.ShowInformation("Info", "Please enter a search term", a.mainWindow)
			return
		}
		// Save current results before searching
		if len(searchResults) > 0 && len(previousResults) == 0 {
			previousResults = searchResults
		}
		results, err := a.contentService.Search(contentType, searchEntry.Text)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}
		// Sort alphabetically by name (or alt_name if name is empty)
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
		// Save current results before lookup
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

	// Add Enter key binding to search entry
	searchEntry.OnSubmitted = func(s string) {
		performSearch()
	}

	// Add Enter key binding to lookup entry
	lookupEntry.OnSubmitted = func(s string) {
		performLookup()
	}

	searchButton := widget.NewButton("Search", performSearch)
	searchButton.Alignment = widget.ButtonAlignCenter

	lookupButton := widget.NewButton("Lookup", performLookup)
	lookupButton.Alignment = widget.ButtonAlignCenter

	loadAllButton := widget.NewButton("Load All", func() {
		results, err := a.contentService.GetAll(contentType)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
			return
		}
		// Sort alphabetically by name (or alt_name if name is empty)
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
		previousResults = nil // Clear previous since we loaded all
		resultsList.Refresh()
		dialog.ShowInformation(
			"Success",
			fmt.Sprintf("Loaded %d entries", len(results)),
			a.mainWindow,
		)
	})
	loadAllButton.Alignment = widget.ButtonAlignCenter

	clearButton := widget.NewButton("Clear", func() {
		// Clear the search/lookup fields
		searchEntry.SetText("")
		lookupEntry.SetText("")

		// Restore previous results or clear completely
		if len(previousResults) > 0 {
			searchResults = previousResults
			previousResults = nil
		} else {
			searchResults = []models.Content{}
		}

		// Clear selection and details in one go
		clearSelection()
		resultsList.Refresh()
	})
	clearButton.Alignment = widget.ButtonAlignCenter

	// ----------------------
	// Controls (top) - Fixed alignment with proper grid layout
	// ----------------------
	// Put buttons in containers with fixed minimum width
	searchButtonContainer := container.NewStack(searchButton)
	lookupButtonContainer := container.NewStack(lookupButton)

	// Use HBox to create aligned rows
	searchBox := container.NewBorder(nil, nil, nil, searchButtonContainer, searchEntry)
	lookupBox := container.NewBorder(nil, nil, nil, lookupButtonContainer, lookupEntry)

	// Create button row with Load All and Clear side by side
	buttonRow := container.NewGridWithColumns(2, loadAllButton, clearButton)

	controlsBox := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("%s Database", contentType)),
		widget.NewSeparator(),
		searchBox,
		lookupBox,
		buttonRow,
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

	// Helper to create selectable, copyable text that looks like a label
	newSelectableText := func(text string) *widget.Entry {
		entry := widget.NewEntry()
		entry.SetText(text)
		entry.Wrapping = fyne.TextWrapWord
		// Make it look like a label but still selectable
		entry.OnChanged = func(s string) {
			// Prevent editing by resetting to original text
			if s != text {
				entry.SetText(text)
			}
		}
		return entry
	}

	// Helper to create clickable hyperlink with wrapping
	newHyperlink := func(label, urlStr string) fyne.CanvasObject {
		if urlStr == "" {
			return newSelectableText(label + " ")
		}
		// Parse the URL using net/url
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			// If URL is invalid, just show as selectable text
			return newSelectableText(label + " " + urlStr)
		}

		// Create hyperlink without the label prefix
		link := widget.NewHyperlink(urlStr, parsedURL)
		link.Wrapping = fyne.TextWrapWord

		// Show label and link separately
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

	deleteBtn := widget.NewButton("Delete Entry", func() {
		// unchanged
	})
	objects = append(objects, deleteBtn)

	return objects
}
