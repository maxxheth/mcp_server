# MCP Server - Independent Model Context Protocol Server Implementation

This is a standalone Model Context Protocol (MCP) server implementation extracted from the ciwg-zoho-migration project. It provides MCP server implementations for:

- **Database Access**: Connect to source and target PostgreSQL databases
- **Wolfram Alpha**: Query computational intelligence for mathematical and scientific answers
- **SSH Tunneling**: Support for remote database connections via SSH

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
- PostgreSQL (for database MCP servers)
- Wolfram Alpha API key (for Wolfram MCP server)

### Installation

```bash
git clone <repo-url>
cd mcp-server
go mod tidy
go build -o mcp-server ./cmd/cli
```

### Usage

#### Start Database MCP Server (Source)

```bash
export MCP_SOURCE_PORT=8080
export SOURCE_DB_HOST=localhost
export SOURCE_DB_PORT=5432
export SOURCE_DB_USERNAME=user
export SOURCE_DB_PASSWORD=password
export SOURCE_DB_DATABASE=source_db

./mcp-server mcp source
```

#### Start Database MCP Server (Target)

```bash
export MCP_TARGET_PORT=8082
export TARGET_DB_HOST=localhost
export TARGET_DB_PORT=5432
export TARGET_DB_USERNAME=user
export TARGET_DB_PASSWORD=password
export TARGET_DB_DATABASE=target_db

./mcp-server mcp target
```

#### Start Wolfram Alpha MCP Server

```bash
export WOLFRAM_API_KEY=your-api-key-here
export MCP_WOLFRAM_PORT=8083

./mcp-server mcp wolfram
```

## Configuration

### Environment Variables

#### MCP Server Ports
- `MCP_SOURCE_PORT` - Port for source database MCP server (default: stdio)
- `MCP_TARGET_PORT` - Port for target database MCP server (default: stdio)
- `MCP_WOLFRAM_PORT` - Port for Wolfram Alpha MCP server (default: stdio)

#### Source Database
- `SOURCE_DB_HOST` - Database host
- `SOURCE_DB_PORT` - Database port
- `SOURCE_DB_USERNAME` - Database username
- `SOURCE_DB_PASSWORD` - Database password
- `SOURCE_DB_DATABASE` - Database name
- `SOURCE_DB_SCHEMA` - Database schema (default: public)

#### Target Database
- `TARGET_DB_HOST` - Database host
- `TARGET_DB_PORT` - Database port
- `TARGET_DB_USERNAME` - Database username
- `TARGET_DB_PASSWORD` - Database password
- `TARGET_DB_DATABASE` - Database name
- `TARGET_DB_SCHEMA` - Database schema (default: public)

#### SSH Tunnel (Source or Target)
- `SSH_SOURCE_REMOTE_HOST` / `SSH_TARGET_REMOTE_HOST` - SSH server host
- `SSH_SOURCE_REMOTE_PORT` / `SSH_TARGET_REMOTE_PORT` - SSH server port
- `SSH_SOURCE_REMOTE_USER` / `SSH_TARGET_REMOTE_USER` - SSH username
- `SSH_SOURCE_PRIVATE_KEY_PATH` / `SSH_TARGET_PRIVATE_KEY_PATH` - Path to SSH private key
- `SSH_SOURCE_LOCAL_PORT` / `SSH_TARGET_LOCAL_PORT` - Local port for tunnel
- `SSH_SOURCE_HOST` / `SSH_TARGET_HOST` - Target database host (as seen from SSH server)
- `SSH_SOURCE_PORT` / `SSH_TARGET_PORT` - Target database port (as seen from SSH server)

#### Wolfram Alpha
- `WOLFRAM_API_KEY` - API key from Wolfram Alpha (required)

## Available Tools

### Database Servers (Source/Target)

- **query** - Execute SQL queries on the database
- **list_tables** - List all tables in the database
- **table_info** - Get detailed information about a table (columns, types, constraints)
- **row_count** - Get the number of rows in a table

### Wolfram Alpha Server

- **query_wolfram** - Query Wolfram Alpha for computational answers

## Examples

### Query a Database via HTTP

```bash
curl -X POST http://localhost:8080/mcp \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "tools/call",
    "params": {
      "name": "query",
      "arguments": {"sql": "SELECT * FROM users LIMIT 10"}
    }
  }'
```

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
    "database-source": {
      "command": "bash",
      "args": [
        "-c",
        "MCP_SOURCE_PORT=8080 SOURCE_DB_HOST=localhost SOURCE_DB_USERNAME=user SOURCE_DB_PASSWORD=pass SOURCE_DB_DATABASE=mydb /path/to/mcp-server mcp source"
      ]
    },
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
