PROJECT=genfig

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGEN=$(GOCMD) generate
GOTEST=$(GOCMD) test

CMD_DIR=$(CURDIR)
BIN_DIR=$(CURDIR)/bin

BINARY=$(BIN_DIR)/$(PROJECT)

STATUS="-alpha"

all: test build

.PHONY: build
build: 
	$(GOGEN) ./...
	$(GOBUILD) -o $(BINARY)  $(CMD_DIR)

.PHONY: test
test: 
	$(GOTEST) -cover ./...

.PHONY: clean
clean: 
	$(GOCLEAN)
	rm -f $(BINARY)

.PHONY: run
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

.PHONY: version
version:
	git tag `cat VERSION`
	git push origin `cat VERSION`