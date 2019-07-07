PROJECT=genfig

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

CMD_DIR=$(CURDIR)
BIN_DIR=$(CURDIR)/bin

BINARY=$(BIN_DIR)/$(PROJECT)

all: test build

.PHONY: build
build: 
	$(GOBUILD) -o $(BINARY)  $(CMD_DIR)

.PHONY: test
test: 
	$(GOTEST) -cover ./...

clean: 
	$(GOCLEAN)
	rm -f $(BINARY)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)