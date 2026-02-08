.PHONY: build run run-tui run-backend test clean help

BINARY_TUI=opencode-tui
BINARY_BACKEND=opencode-backend
GO=go
GOFLAGS=-v

help:
	@echo "OpenCode Go Project - Available targets:"
	@echo "  make build          - Build both TUI and backend"
	@echo "  make build-tui      - Build only TUI"
	@echo "  make build-backend  - Build only backend"
	@echo "  make run-tui        - Run TUI directly"
	@echo "  make run-backend    - Run backend directly"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean binaries"
	@echo "  make mod-download   - Download dependencies"

build: build-tui build-backend

build-tui:
	$(GO) build $(GOFLAGS) -o bin/$(BINARY_TUI) ./cmd/tui

build-backend:
	$(GO) build $(GOFLAGS) -o bin/$(BINARY_BACKEND) ./cmd/backend

run-tui: build-tui
	./bin/$(BINARY_TUI)

run-backend: build-backend
	./bin/$(BINARY_BACKEND)

test:
	$(GO) test $(GOFLAGS) ./...

clean:
	$(GO) clean
	rm -f bin/$(BINARY_TUI) bin/$(BINARY_BACKEND)

mod-download:
	$(GO) mod download
	$(GO) mod tidy
