.PHONY: deb windows all

all: deb windows

deb:
	$(MAKE) -C builds/deb

windows:
	$(MAKE) -C builds/windows
