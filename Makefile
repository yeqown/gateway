CMD=go
BUILD=$(CMD) build
BIN=gateway
MAIN=.
WEBVIEW=webview
BUILDOUT=./.build

LINUX=linux
OSX="osx"

default: clean build

build: build-osx build-linux
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

clean:
	- cd $(MAIN) rm $(BIN)
	- cd $(MAIN) rm -fr $(BUILDOUT)
