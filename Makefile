.PHONY: build run clean help deps

BINARY_NAME=dataforge
BINARY_UNIX=$(BINARY_NAME)
BINARY_WIN=$(BINARY_NAME).exe

help:
	@echo "Available targets:"
	@echo "  make build     - Build the CLI tool for current OS"
	@echo "  make build-win - Build Windows binary"
	@echo "  make build-linux - Build Linux binary"
	@echo "  make build-mac - Build macOS binary"
	@echo "  make run       - Build and run with sample database"
	@echo "  make clean     - Remove built binaries"
	@echo "  make deps      - Download dependencies"
	@echo "  make test      - Run tests"

deps:
	go mod download
	go mod tidy

build:
	go build -o $(BINARY_NAME) -v

build-win:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_WIN) -v

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_UNIX) -v

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_UNIX) -v

run: build
	./$(BINARY_NAME) start --db "postgres://postgres:postgres@localhost:5432/shop" --interval 5s

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WIN)
	rm -f $(BINARY_UNIX)

test:
	go test -v ./...
