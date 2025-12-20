.PHONY: deb windows all clean

all: deb windows

deb:
	$(MAKE) -C builds/deb

windows:
	$(MAKE) -C builds/windows

clean:
	$(MAKE) -C builds/deb clean
	$(MAKE) -C builds/windows clean

