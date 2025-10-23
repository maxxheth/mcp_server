# MCP Server - Independent Model Context Protocol Server Implementation

This is a standalone Model Context Protocol (MCP) server implementation. It provides MCP server implementations for:

- **Wolfram Alpha**: Query computational intelligence for mathematical and scientific answers

## Project Structure

```
mcp-server/
├── cmd/
│   └── cli/                  # CLI entry point and command handlers
│       ├── main.go          # Application entry point
│       ├── root.go          # Root command definition
│       └── mcp.go           # MCP server commands (source, target, wolfram)
├── internal/
│   ├── mcp/                 # MCP server implementation
│   │   ├── server.go       # MCP server core
│   │   ├── protocol.go     # MCP protocol handler
│   │   └── wolfram.go      # Wolfram Alpha integration
│   └── migration/           # Shared migration utilities
│       ├── introspection/  # Database introspection
│       └── remote/         # SSH tunnel support
├── go.mod                  # Go module definition
└── README.md              # This file
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Wolfram Alpha API key (for Wolfram MCP server)

### Installation

```bash
git clone <repo-url>
cd mcp-server
go mod tidy
go build -o mcp-server ./cmd/cli
```

### Usage

#### Start Wolfram Alpha MCP Server

```bash
export WOLFRAM_API_KEY=your-api-key-here
export MCP_WOLFRAM_PORT=8083

./mcp-server mcp wolfram
```

## Configuration

### Environment Variables

#### MCP Server Ports
- `MCP_WOLFRAM_PORT` - Port for Wolfram Alpha MCP server (default: stdio)

#### Wolfram Alpha
- `WOLFRAM_API_KEY` - API key from Wolfram Alpha (required)

## Available Tools

### Wolfram Alpha Server

- **query_wolfram** - Query Wolfram Alpha for computational answers

## Examples

### Query Wolfram Alpha

```bash
curl -X POST http://localhost:8083/mcp \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "tools/call",
    "params": {
      "name": "query_wolfram",
      "arguments": {"query": "integrate x^2 from 0 to 1"}
    }
  }'
```

## MCP Protocol Support

- **Initialize**: Setup and handshake
- **Ping**: Health check
- **Tools/List**: Discover available tools
- **Tools/Call**: Execute tools
- **Resources/List**: List available database resources (tables)
- **Resources/Read**: Read table information

## Building

```bash
# Build the binary
go build -o mcp-server ./cmd/cli

# Run tests (if any)
go test ./...

# Create a release build
go build -ldflags="-s -w" -o mcp-server ./cmd/cli
```

## Integration with AI Assistants

### Claude Desktop

Add to `~/.config/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "wolfram-alpha": {
      "command": "bash",
      "args": [
        "-c",
        "WOLFRAM_API_KEY=your-key MCP_WOLFRAM_PORT=8083 /path/to/mcp-server mcp wolfram"
      ]
    }
  }
}
```

## Security Considerations

- Never commit API keys or database credentials to version control
- Use environment variables or `.env` files to manage secrets
- For HTTP mode, ensure the server only listens on localhost or behind authentication
- SSH tunnels provide secure remote database access
- Keep the binary and dependencies updated

## License

[Specify your license here]

## Contributing

[Add contribution guidelines if needed]

## Related Projects

This is extracted from the [ciwg-zoho-migration](https://github.com/maxxheth/ciwg-zoho-migration) project. See that repository for the complete migration system.

## Support

For issues, questions, or suggestions, please refer to the main project repository.
