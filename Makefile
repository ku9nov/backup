# Specify the name of the output binary
BINARY_NAME := backup
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
# Go source files
SOURCES := $(wildcard *.go)

# Default target: build the binary
.PHONY: all
all: clean $(BINARY_NAME)

# Target to build the binary
$(BINARY_NAME): $(SOURCES)
	go build -o $(BINARY_NAME) $(SOURCES)

# Target to clean up generated files
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
