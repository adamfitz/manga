.PHONY: deb windows all clean

all: deb windows

deb:
	$(MAKE) -C builds/deb

clean:
	$(MAKE) -C builds/deb clean

