GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
BINARY_NAME=webapp
GOLIB=\

BINARY_LINUX=$(BINARY_NAME)_linux_amd64

.PHONY: all build clean run build-linux format

all: build build-linux
build : $(BINARY_NAME)

build-linux : $(BINARY_LINUX)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)

run: $(BINARY_NAME)
	./$(BINARY_NAME)

format: $(BINARY_NAME).go $(GOLIB)
	$(GOFMT) $(BINARY_NAME).go
	$(GOFMT) $(GOLIB)

$(BINARY_NAME): $(BINARY_NAME).go $(GOLIB)
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) $(BINARY_NAME).go

$(BINARY_LINUX): $(BINARY_NAME).go $(GOLIB)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_LINUX) $(BINARY_NAME).go
