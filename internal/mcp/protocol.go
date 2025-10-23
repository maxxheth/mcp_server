package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Protocol implements the Model Context Protocol
type Protocol struct {
	server        *Server
	wolframClient *WolframAlphaClient
}

// NewProtocol creates a new MCP protocol handler
func NewProtocol(server *Server) *Protocol {
	p := &Protocol{
		server: server,
	}
	// Initialize Wolfram client if API key is available
	apiKey := os.Getenv("WOLFRAM_API_KEY")
	if apiKey != "" {
		p.wolframClient = NewWolframAlphaClient(apiKey)
	}
	return p
}

// Request represents an MCP request
type Request struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params,omitempty"`
	ID     interface{}            `json:"id,omitempty"`
}

// Response represents an MCP response
type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  *Error      `json:"error,omitempty"`
	ID     interface{} `json:"id,omitempty"`
}

// Error represents an MCP error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Handle processes an MCP request
func (p *Protocol) Handle(req Request) Response {
	log.Printf("Handling MCP request: %s", req.Method)

	switch req.Method {
	// MCP Standard Lifecycle Methods
	case "initialize":
		return p.handleInitialize(req)
	case "initialized":
		return p.handleInitialized(req)
	case "resources/list":
		return p.handleResourcesList(req)
	case "resources/read":
		return p.handleResourcesRead(req)
	case "tools/list":
		return p.handleToolsList(req)
	case "tools/call":
		return p.handleToolsCall(req)
	case "ping":
		return p.handlePing(req)
	case "server_info":
		return p.handleServerInfo(req)
	default:
		return Response{
			Error: &Error{Code: -32601, Message: "Method not found"},
			ID:    req.ID,
		}
	}
}

// handleInitialize handles MCP initialization
func (p *Protocol) handleInitialize(req Request) Response {
	return Response{
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]bool{"tools": true},
			"serverInfo": map[string]interface{}{
				"name":    "MCP Wolfram Server",
				"version": "1.0.0",
			},
		},
		ID: req.ID,
	}
}

// handleInitialized handles post-initialization
func (p *Protocol) handleInitialized(req Request) Response {
	return Response{Result: nil, ID: req.ID}
}

// handleResourcesList lists available resources
func (p *Protocol) handleResourcesList(req Request) Response {
	return Response{
		Result: map[string]interface{}{
			"resources": []map[string]interface{}{},
		},
		ID: req.ID,
	}
}

// handleResourcesRead reads a resource
func (p *Protocol) handleResourcesRead(req Request) Response {
	return Response{
		Result: map[string]interface{}{
			"contents": []map[string]interface{}{},
		},
		ID: req.ID,
	}
}

// handleToolsList lists available tools
func (p *Protocol) handleToolsList(req Request) Response {
	tools := []map[string]interface{}{
		{
			"name": "query_wolfram",
			"description": "Query Wolfram Alpha for mathematical, scientific, and computational answers. " +
				"Supports math calculations, symbolic operations, scientific data, unit conversions, and natural language questions.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The query to send to Wolfram Alpha",
					},
				},
				"required": []string{"query"},
			},
		},
	}

	return Response{
		Result: map[string]interface{}{
			"tools": tools,
		},
		ID: req.ID,
	}
}

// handleToolsCall handles tool calls
func (p *Protocol) handleToolsCall(req Request) Response {
	toolName, ok := req.Params["name"].(string)
	if !ok {
		return Response{
			Error: &Error{Code: -32602, Message: "Invalid parameters: name is required"},
			ID:    req.ID,
		}
	}

	arguments, ok := req.Params["arguments"].(map[string]interface{})
	if !ok {
		return Response{
			Error: &Error{Code: -32602, Message: "Invalid parameters: arguments is required"},
			ID:    req.ID,
		}
	}

	switch toolName {
	case "query_wolfram":
		query, ok := arguments["query"].(string)
		if !ok {
			return Response{
				Error: &Error{Code: -32602, Message: "Invalid arguments: query is required"},
				ID:    req.ID,
			}
		}

		if p.wolframClient == nil {
			return Response{
				Error: &Error{Code: -32603, Message: "Wolfram Alpha API key not configured"},
				ID:    req.ID,
			}
		}

		result, err := p.wolframClient.Query(query)
		if err != nil {
			return Response{
				Error: &Error{Code: -32603, Message: fmt.Sprintf("Wolfram query failed: %v", err)},
				ID:    req.ID,
			}
		}

		resultText := FormatResultAsText(result)

		return Response{
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": resultText,
					},
				},
			},
			ID: req.ID,
		}
	default:
		return Response{
			Error: &Error{Code: -32601, Message: fmt.Sprintf("Tool not found: %s", toolName)},
			ID:    req.ID,
		}
	}
}

// handlePing handles ping requests
func (p *Protocol) handlePing(req Request) Response {
	return Response{
		Result: map[string]interface{}{
			"pong": true,
		},
		ID: req.ID,
	}
}

// handleServerInfo returns server information
func (p *Protocol) handleServerInfo(req Request) Response {
	return Response{
		Result: map[string]interface{}{
			"name":    "MCP Wolfram Server",
			"version": "1.0.0",
			"type":    "wolfram",
		},
		ID: req.ID,
	}
}

// ServeStdio handles stdio-based MCP communication
func (p *Protocol) ServeStdio(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	encoder := json.NewEncoder(writer)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("Failed to parse request: %v", err)
			continue
		}

		response := p.Handle(req)
		if err := encoder.Encode(response); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	}

	return scanner.Err()
}

// ServeHTTP serves MCP protocol over HTTP
func (p *Protocol) ServeHTTP(port string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Error: &Error{Code: -32700, Message: "Parse error"},
			})
			return
		}

		response := p.Handle(req)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Printf("MCP HTTP server listening on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}
