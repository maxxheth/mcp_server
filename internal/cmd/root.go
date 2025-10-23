package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for the MCP server CLI.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "Model Context Protocol (MCP) Server - Start MCP servers for databases and APIs",
		Long: `mcp-server provides Model Context Protocol (MCP) server implementations for:
- Database access (source and target databases)
- Wolfram Alpha computational queries
- SSH tunnel support for remote database access

Use this tool to start MCP servers that can be integrated with AI assistants
and other tools that support the Model Context Protocol.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Default action when no subcommand is provided
			fmt.Println("Welcome to the MCP Server!")
			fmt.Println("Use 'help' to see available commands.")
			fmt.Println("")
			fmt.Println("Available servers:")
			fmt.Println("  - mcp-server mcp source    Start source database MCP server")
			fmt.Println("  - mcp-server mcp target    Start target database MCP server")
			fmt.Println("  - mcp-server mcp wolfram   Start Wolfram Alpha MCP server")
		},
	}

	// Add MCP subcommand only (focused on MCP functionality)
	rootCmd.AddCommand(MCPCmd)

	return rootCmd
}

// Execute runs the root command.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
