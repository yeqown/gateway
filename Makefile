CMD=go
BUILD=$(CMD) build
BIN=gateway
MAIN=.

default: clean build

build:
	@ echo "build"
	$(BUILD) -o $(BIN)

clean:
	- cd $(MAIN) && rm gateway