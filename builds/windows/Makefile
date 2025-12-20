APP_NAME := mangadb
ICON     := ../../assets/icons/mangadb.ico

all: exe

exe:
	fyne package \
		-os windows \
		-name $(APP_NAME) \
		-icon $(ICON) \
		-release

	mv *.exe ./$(APP_NAME).exe || true

clean:
	rm -f *.exe
	rm -rf fyne-*
