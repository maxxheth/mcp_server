package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"mcp-server/internal/mcp"
)

// MCPCmd is the main command for MCP servers
var MCPCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Model Context Protocol (MCP) server",
	Long: `Start MCP servers for Wolfram Alpha computational queries.

The server communicates via JSON-RPC over HTTP or stdin/stdout.`,
}

// mcpWolframCmd starts the MCP server for Wolfram Alpha queries
var mcpWolframCmd = &cobra.Command{
	Use:   "wolfram",
	Short: "Start MCP server for Wolfram Alpha queries",
	Long: `Starts an MCP server that provides access to Wolfram Alpha for mathematical, 
scientific, and computational queries.

The server communicates via JSON-RPC over HTTP (if MCP_WOLFRAM_PORT is set)
or stdin/stdout (default). It provides the query_wolfram tool for:
- Mathematical calculations and symbolic operations
- Scientific data and constants
- Unit conversions
- Weather, geography, and other real-world data
- Computational answers to natural language questions

Environment variables:
  WOLFRAM_API_KEY         - Your Wolfram Alpha API key (required for queries)
  MCP_WOLFRAM_PORT        - Port for HTTP server (optional, default: stdio mode)`,
	RunE: runMCPWolfram,
}

func init() {
	MCPCmd.AddCommand(mcpWolframCmd)
}

// runMCPWolfram starts the MCP server for Wolfram Alpha
func runMCPWolfram(cmd *cobra.Command, args []string) error {
	log.Println("Starting MCP server for Wolfram Alpha...")

	// Check API key
	apiKey := os.Getenv("WOLFRAM_API_KEY")
	if apiKey == "" {
		log.Println("Warning: WOLFRAM_API_KEY environment variable not set. Queries will fail.")
	}

	config := mcp.ServerConfig{
		Name: "mcp-wolfram-server",
		Type: "wolfram",
	}

	server, err := mcp.NewWolframServer(config)
	if err != nil {
		return fmt.Errorf("failed to create MCP Wolfram server: %w", err)
	}
	defer server.Close()

	protocol := mcp.NewProtocol(server)

	// Use HTTP server if port is configured, otherwise use stdio
	if server.Config.MCPPort != "" {
		log.Printf("MCP Wolfram server ready. Listening on HTTP port %s...", server.Config.MCPPort)
		if err := protocol.ServeHTTP(server.Config.MCPPort); err != nil {
			return fmt.Errorf("HTTP server error: %w", err)
		}
	} else {
		log.Println("MCP Wolfram server ready. Listening on stdin/stdout...")
		log.Println("Available tools: query_wolfram")
		if err := protocol.ServeStdio(os.Stdin, os.Stdout); err != nil {
			return fmt.Errorf("stdio server error: %w", err)
		}
	}

	return nil
}
