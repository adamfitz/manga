APP_NAME := mangadb
ICON     := ../../assets/icons/mangadb.ico
VERSION := $(shell git describe --tags --abbrev=0)

all: exe

exe:
	fyne package \
		-os windows \
		-name $(APP_NAME) \
		-icon $(ICON) \
		-release \
		-ldflags "-X mangadb/utils.Version=$(VERSION)"

	mv *.exe ./$(APP_NAME).exe || true

clean:
	rm -f *.exe
	rm -rf fyne-*
