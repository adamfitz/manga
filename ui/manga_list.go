package ui

import (
	"fmt"
	"manga/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) createMangaListView() *container.Split {
	list := widget.NewList(
		func() int {
			if a.mangaService == nil {
				return 0
			}
			mangas, _ := a.mangaService.GetAll()
			return len(mangas)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if a.mangaService == nil {
				return
			}
			mangas, _ := a.mangaService.GetAll()
			if id < len(mangas) {
				obj.(*widget.Label).SetText(mangas[id].Title)
			}
		},
	)

	detailContainer := container.NewVBox(
		widget.NewLabel("Select a manga to view details"),
	)

	list.OnSelected = func(id widget.ListItemID) {
		if a.mangaService == nil {
			return
		}
		mangas, _ := a.mangaService.GetAll()
		if id < len(mangas) {
			detailContainer.Objects = a.createMangaDetail(&mangas[id], list)
			detailContainer.Refresh()
		}
	}

	addButton := widget.NewButton("Add Manga", func() {
		a.showAddMangaDialog(list)
	})

	leftSide := container.NewBorder(nil, addButton, nil, nil, list)
	return container.NewHSplit(leftSide, detailContainer)
}

func (a *App) createMangaDetail(manga *models.Manga, list *widget.List) []fyne.CanvasObject {
	title := widget.NewLabel(fmt.Sprintf("Title: %s", manga.Title))
	title.Wrapping = fyne.TextWrapWord

	author := widget.NewLabel(fmt.Sprintf("Author: %s", manga.Author))
	status := widget.NewLabel(fmt.Sprintf("Status: %s", manga.Status))

	description := widget.NewLabel(fmt.Sprintf("Description:\n%s", manga.Description))
	description.Wrapping = fyne.TextWrapWord

	// Check if bookmarked
	bookmark, _ := a.bookmarkService.GetByMangaID(manga.ID)
	bookmarkBtn := widget.NewButton("Add Bookmark", func() {
		a.showAddBookmarkDialog(manga.ID, list)
	})

	if bookmark != nil {
		bookmarkBtn.SetText("Remove Bookmark")
		bookmarkBtn.OnTapped = func() {
			if err := a.bookmarkService.Delete(bookmark.ID); err != nil {
				dialog.ShowError(err, a.mainWindow)
				return
			}
			dialog.ShowInformation("Success", "Bookmark removed", a.mainWindow)
			list.Refresh()
		}
	}

	deleteBtn := widget.NewButton("Delete Manga", func() {
		dialog.ShowConfirm("Delete Manga",
			fmt.Sprintf("Are you sure you want to delete '%s'?", manga.Title),
			func(ok bool) {
				if !ok {
					return
				}
				if err := a.mangaService.Delete(manga.ID); err != nil {
					dialog.ShowError(err, a.mainWindow)
					return
				}
				dialog.ShowInformation("Success", "Manga deleted", a.mainWindow)
				list.Refresh()
				list.UnselectAll()
			}, a.mainWindow)
	})

	buttons := container.NewHBox(bookmarkBtn, deleteBtn)

	return []fyne.CanvasObject{
		title,
		author,
		status,
		widget.NewSeparator(),
		description,
		widget.NewSeparator(),
		buttons,
	}
}

func (a *App) showAddMangaDialog(list *widget.List) {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Manga Title")

	authorEntry := widget.NewEntry()
	authorEntry.SetPlaceHolder("Author")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Description")

	statusSelect := widget.NewSelect([]string{"ongoing", "completed", "hiatus", "cancelled"}, nil)
	statusSelect.SetSelected("ongoing")

	form := dialog.NewForm("Add New Manga", "Add", "Cancel",
		[]*widget.FormItem{
			{Text: "Title", Widget: titleEntry},
			{Text: "Author", Widget: authorEntry},
			{Text: "Status", Widget: statusSelect},
			{Text: "Description", Widget: descEntry},
		},
		func(ok bool) {
			if !ok {
				return
			}

			manga := &models.Manga{
				Title:       titleEntry.Text,
				Author:      authorEntry.Text,
				Description: descEntry.Text,
				Status:      statusSelect.Selected,
			}

			if err := a.mangaService.Create(manga); err != nil {
				dialog.ShowError(err, a.mainWindow)
				return
			}

			dialog.ShowInformation("Success", "Manga added", a.mainWindow)
			list.Refresh()
		}, a.mainWindow)

	form.Resize(fyne.NewSize(500, 400))
	form.Show()
}

func (a *App) showAddBookmarkDialog(mangaID int, list *widget.List) {
	noteEntry := widget.NewMultiLineEntry()
	noteEntry.SetPlaceHolder("Add a note (optional)")

	dialog.ShowForm("Add Bookmark", "Add", "Cancel",
		[]*widget.FormItem{
			{Text: "Note", Widget: noteEntry},
		},
		func(ok bool) {
			if !ok {
				return
			}

			bookmark := &models.Bookmark{
				MangaID: mangaID,
				Note:    noteEntry.Text,
			}

			if err := a.bookmarkService.Create(bookmark); err != nil {
				dialog.ShowError(err, a.mainWindow)
				return
			}

			dialog.ShowInformation("Success", "Bookmark added", a.mainWindow)
			list.Refresh()
		}, a.mainWindow)
}
