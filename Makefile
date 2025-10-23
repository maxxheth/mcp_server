.PHONY: build dev start stop clean help

# Build output
BINARY_NAME=mcp-server

help:
	@echo "MCP Server - Available Commands"
	@echo ""
	@echo "make build       - Build the binary"
	@echo "make dev        - Run in development mode with hot reload"
	@echo "make start      - Start the Wolfram Alpha MCP server"
	@echo "make stop       - Stop the running MCP server"
	@echo "make clean      - Clean build artifacts"
	@echo "make help       - Show this help message"

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd/cli
	@echo "✓ Build complete: ./$(BINARY_NAME)"

dev:
	@echo "Starting development server with hot reload..."
	@which air > /dev/null || (echo "Installing air for hot reload..." && go install github.com/cosmtrek/air@latest)
	air

start:
	@echo "Starting Wolfram Alpha MCP server..."
	@if [ -z "$$WOLFRAM_API_KEY" ]; then \
		echo "Error: WOLFRAM_API_KEY environment variable not set"; \
		exit 1; \
	fi
	./$(BINARY_NAME) mcp wolfram

stop:
	@echo "Stopping MCP server..."
	pkill -f "$(BINARY_NAME) mcp wolfram" || echo "No process found"
	@echo "✓ Server stopped"

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	go clean
	@echo "✓ Clean complete"
