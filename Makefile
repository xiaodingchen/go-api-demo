# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_DIR = bin
BINARY_NAME=test
BINARY_UNIX=$(BINARY_NAME)_unix

all: build
build:
		rm -rf $(BINARY_DIR)
		$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
vendor:
		go mod vendor -v
tidy:
		go mod tidy -v
run_api:
	go run main.go api -c ./config/config.toml
clean:
		$(GOCLEAN)
		rm -rf $(BINARY_DIR)