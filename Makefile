# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_DIR = bin
BINARY_NAME=test
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
		rm -rf $(BINARY_DIR)
		$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
vendor:
		go mod vendor -v
clean:
		$(GOCLEAN)
		rm -rf $(BINARY_DIR)