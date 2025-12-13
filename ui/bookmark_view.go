package ui

import (
	"fmt"
	"manga/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (a *App) createBookmarkView() *container.Split {
	list := widget.NewList(
		func() int {
			bookmarks, _ := a.bookmarkService.GetAll()
			return len(bookmarks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			bookmarks, _ := a.bookmarkService.GetAll()
			if id < len(bookmarks) && bookmarks[id].Manga != nil {
				obj.(*widget.Label).SetText(bookmarks[id].Manga.Title)
			}
		},
	)

	detailContainer := container.NewVBox(
		widget.NewLabel("Select a bookmark to view details"),
	)

	list.OnSelected = func(id widget.ListItemID) {
		bookmarks, _ := a.bookmarkService.GetAll()
		if id < len(bookmarks) {
			detailContainer.Objects = a.createBookmarkDetail(&bookmarks[id], list)
			detailContainer.Refresh()
		}
	}

	return container.NewHSplit(list, detailContainer)
}

func (a *App) createBookmarkDetail(bookmark *models.Bookmark, list *widget.List) []fyne.CanvasObject {
	if bookmark.Manga == nil {
		return []fyne.CanvasObject{widget.NewLabel("No manga data available")}
	}

	title := widget.NewLabel(fmt.Sprintf("Title: %s", bookmark.Manga.Title))
	title.Wrapping = fyne.TextWrapWord

	author := widget.NewLabel(fmt.Sprintf("Author: %s", bookmark.Manga.Author))

	note := widget.NewLabel(fmt.Sprintf("Note: %s", bookmark.Note))
	note.Wrapping = fyne.TextWrapWord

	deleteBtn := widget.NewButton("Remove Bookmark", func() {
		dialog.ShowConfirm("Remove Bookmark",
			"Are you sure you want to remove this bookmark?",
			func(ok bool) {
				if !ok {
					return
				}
				if err := a.bookmarkService.Delete(bookmark.ID); err != nil {
					dialog.ShowError(err, a.mainWindow)
					return
				}
				dialog.ShowInformation("Success", "Bookmark removed", a.mainWindow)
				list.Refresh()
				list.UnselectAll()
			}, a.mainWindow)
	})

	return []fyne.CanvasObject{
		title,
		author,
		widget.NewSeparator(),
		note,
		widget.NewSeparator(),
		deleteBtn,
	}
}
