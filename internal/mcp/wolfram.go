package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

// WolframAlphaClient wraps the Wolfram Alpha API
type WolframAlphaClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewWolframAlphaClient creates a new Wolfram Alpha client
func NewWolframAlphaClient(apiKey string) *WolframAlphaClient {
	if apiKey == "" {
		apiKey = os.Getenv("WOLFRAM_API_KEY")
	}
	return &WolframAlphaClient{
		apiKey:     apiKey,
		baseURL:    "https://www.wolframalpha.com/api/v1",
		httpClient: &http.Client{},
	}
}

// QueryResult represents a result from Wolfram Alpha
type QueryResult struct {
	QueryResult struct {
		Success bool  `json:"success"`
		Error   bool  `json:"error"`
		Numpods int   `json:"numpods"`
		Pods    []Pod `json:"pods"`
	} `json:"queryresult"`
}

// Pod represents a "pod" (section) in the Wolfram Alpha response
type Pod struct {
	Title    string   `json:"title"`
	Scanner  string   `json:"scanner"`
	Error    bool     `json:"error"`
	Subpods  []Subpod `json:"subpods"`
	Position int      `json:"position"`
}

// Subpod represents content within a pod
type Subpod struct {
	Title     string `json:"title"`
	Plaintext string `json:"plaintext"`
	Image     Image  `json:"image"`
}

// Image represents an image in a subpod
type Image struct {
	Src    string `json:"src"`
	Alt    string `json:"alt"`
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Query executes a query against Wolfram Alpha
func (w *WolframAlphaClient) Query(input string) (*QueryResult, error) {
	if w.apiKey == "" {
		return nil, fmt.Errorf("WOLFRAM_API_KEY environment variable not set")
	}

	// Build the request URL
	params := url.Values{}
	params.Set("input", input)
	params.Set("appid", w.apiKey)
	params.Set("output", "json")

	fullURL := fmt.Sprintf("%s/query?%s", w.baseURL, params.Encode())

	log.Printf("Querying Wolfram Alpha: %s", input)

	// Make the request
	resp, err := w.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query Wolfram Alpha: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Wolfram Alpha API error: status %d, %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Wolfram Alpha response: %w", err)
	}

	return &result, nil
}

// FormatResultAsText converts a Wolfram Alpha result to readable text
func FormatResultAsText(result *QueryResult) string {
	if !result.QueryResult.Success {
		if result.QueryResult.Error {
			return "Error: Wolfram Alpha could not process this query"
		}
		return "No results found"
	}

	var output string
	output += fmt.Sprintf("Wolfram Alpha Result (%d sections):\n\n", result.QueryResult.Numpods)

	for _, pod := range result.QueryResult.Pods {
		output += fmt.Sprintf("=== %s ===\n", pod.Title)

		for _, subpod := range pod.Subpods {
			if subpod.Plaintext != "" {
				output += subpod.Plaintext + "\n"
			}
			if subpod.Image.Src != "" {
				output += fmt.Sprintf("[Image: %s]\n", subpod.Image.Src)
			}
		}
		output += "\n"
	}

	return output
}
