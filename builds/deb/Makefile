APP_NAME := mangadb
VERSION := $(shell git describe --tags --abbrev=0)
CLEAN_VERSION := $(shell echo $(VERSION) | sed 's/^v//')

GOOS     := linux
GOARCH   := amd64
BINARY   := $(APP_NAME)

ROOTFS       := rootfs
DEBIAN_DIR   := $(ROOTFS)/DEBIAN
BIN_DIR      := $(ROOTFS)/usr/bin
DESKTOP_DIR  := $(ROOTFS)/usr/share/applications
ICON_DIR     := $(ROOTFS)/usr/share/icons/hicolor/256x256/apps

all: deb

deb: clean prepare build install_files package

prepare:
	mkdir -p $(DEBIAN_DIR) $(BIN_DIR) $(DESKTOP_DIR) $(ICON_DIR)
	sed "s/@VERSION@/$(CLEAN_VERSION)/" DEBIAN/control.in > $(DEBIAN_DIR)/control
	cp mangadb.desktop $(DESKTOP_DIR)/$(APP_NAME).desktop
	cp ../../assets/icons/$(APP_NAME).png $(ICON_DIR)/$(APP_NAME).png

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X mangadb/utils.Version=$(VERSION)" -o $(BIN_DIR)/$(BINARY) ../../main.go

install_files:
	chmod 755 $(BIN_DIR)/$(BINARY)
	chmod 644 $(DESKTOP_DIR)/$(APP_NAME).desktop
	chmod 644 $(ICON_DIR)/$(APP_NAME).png

package:
	dpkg-deb --build $(ROOTFS) ../$(APP_NAME)_$(VERSION)_amd64.deb

clean:
	rm -rf $(ROOTFS)
	rm -f ../../builds/*.deb
	rm -rf mangadb_*/
