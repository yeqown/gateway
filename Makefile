CMD=go
BUILD=$(CMD) build
BIN=gateway
MAIN=.
WEBVIEW=webview
WEBVIEW_TAR=webview.tar
BUILDOUT=./.build

LINUX=linux
OSX="osx"

default: clean build

build: build-osx build-linux build-frontend
	@ echo "build done"

build-linux:
	@ echo "build linux"
	- mkdir -p $(BUILDOUT)/$(LINUX)
	@ echo "build linux"
	GOOS=linux GOARCH=amd64 $(BUILD) -o $(BIN)
	mv $(BIN) $(BUILDOUT)/$(LINUX)

build-osx:
	@ echo "build osx"
	- mkdir -p $(BUILDOUT)/$(OSX) 
	@ echo "build osx"
	GOOS=darwin GOARCH=amd64 $(BUILD) -o $(BIN)
	mv $(BIN) $(BUILDOUT)/$(OSX)

build-frontend:
	@ mkdir -p $(BUILDOUT)/$(WEBVIEW)
	cd $(WEBVIEW) && npm run build
	mv -f $(WEBVIEW)/dist/* $(BUILDOUT)/$(WEBVIEW)
	cd $(BUILDOUT)/$(WEBVIEW) && tar -zcvf $(WEBVIEW_TAR) ./*
	mv $(BUILDOUT)/$(WEBVIEW)/$(WEBVIEW_TAR) $(BUILDOUT)

clean:
	- rm -fr $(BUILDOUT)
