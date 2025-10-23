package mcp

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// ServerConfig holds the configuration for an MCP server
type ServerConfig struct {
	Name       string
	Type       string // "source", "target", or "wolfram"
	MCPPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSchema   string
	SSHEnabled bool
}

// Server represents an MCP server
type Server struct {
	Config ServerConfig
	db     *sql.DB // Only used for database servers
	ctx    context.Context
}

// NewWolframServer creates a new MCP server for Wolfram Alpha (no database needed)
func NewWolframServer(config ServerConfig) (*Server, error) {
	server := &Server{
		Config: config,
		ctx:    context.Background(),
	}

	// Load MCP port from environment
	if server.Config.MCPPort == "" {
		server.Config.MCPPort = getEnv("MCP_WOLFRAM_PORT", "MCP_PORT")
	}

	log.Printf("MCP Wolfram Server '%s' initialized successfully", server.Config.Name)
	return server, nil
}

// Close closes the server connection
func (s *Server) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// getEnv gets an environment variable with fallback support
func getEnv(keys ...string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return ""
}
