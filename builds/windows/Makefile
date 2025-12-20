APP_NAME := mangadb
ICON     := ../../assets/icons/mangadb.ico
VERSION  := $(shell git describe --tags --abbrev=0)
CLEAN_VERSION := $(shell echo $(VERSION) | sed 's/^v//')

all: exe

exe:
	# Build the Windows binary with version injection
	GOOS=windows GOARCH=amd64 go build -ldflags "-X mangadb/utils.Version=$(VERSION)" -o $(APP_NAME).exe ../../main.go

	# Package using the prebuilt executable (no ldflags here)
	fyne package \
		-os windows \
		-name $(APP_NAME) \
		-icon $(ICON) \
		-release \
		-executable $(APP_NAME).exe \
		-appVersion $(CLEAN_VERSION) \
		-appID com.app.mangadb

	mv *.exe ./$(APP_NAME).exe || true

clean:
	rm -f *.exe
	rm -rf fyne-*
