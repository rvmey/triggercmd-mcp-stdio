package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listInput struct{}

type runInput struct {
	Command    string `json:"command" jsonschema:"Name of the command to run"`
	Computer   string `json:"computer" jsonschema:"Name of the computer that owns the command"`
	Parameters string `json:"parameters,omitempty" jsonschema:"Optional parameters to pass to the command"`
}

// getTriggerCmdToken returns the TriggerCMD token, preferring the TRIGGERCMD_TOKEN
// environment variable and falling back to ~/.TRIGGERcmdData/token.tkn.
// Works on Windows, macOS, and Linux via os.UserHomeDir.
func getTriggerCmdToken() (string, error) {
	if tok := strings.TrimSpace(os.Getenv("TRIGGERCMD_TOKEN")); tok != "" {
		return tok, nil
	}
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		if err == nil {
			err = fmt.Errorf("empty home directory")
		}
		return "", fmt.Errorf("unable to resolve home directory: %w", err)
	}
	path := filepath.Join(home, ".TRIGGERcmdData", "token.tkn")
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed reading %s: %w", path, err)
	}
	tok := strings.TrimSpace(string(b))
	if tok == "" {
		return "", fmt.Errorf("token file %s is empty", path)
	}
	return tok, nil
}

func listCommands(ctx context.Context, _ *mcp.CallToolRequest, _ listInput) (*mcp.CallToolResult, any, error) {
	log.Println("TriggerCMD list_commands tool called")
	url := "https://triggercmd.com/api/command/list"
	token, err := getTriggerCmdToken()
	if err != nil {
		// Tool error (user-visible)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Missing TriggerCMD token: " + err.Error()}},
			IsError: true,
		}, nil, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var jsonResp map[string]any
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		log.Println("Failed to parse JSON. Raw response:", string(body))
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(body)}}}, nil, nil
	}

	// Build a simplified list: name, voice, computer.name
	simplified := make([]map[string]any, 0)
	if recs, ok := jsonResp["records"].([]any); ok {
		for _, r := range recs {
			if m, ok := r.(map[string]any); ok {
				var compName any
				if comp, ok := m["computer"].(map[string]any); ok {
					compName = comp["name"]
				}
				simplified = append(simplified, map[string]any{
					"name":     m["name"],
					"voice":    m["voice"],
					"computer": compName,
				})
			}
		}
	} else {
		log.Println("Unexpected API response format:", jsonResp)
	}
	b, _ := json.MarshalIndent(simplified, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, nil, nil
}

func runCommand(ctx context.Context, _ *mcp.CallToolRequest, in runInput) (*mcp.CallToolResult, any, error) {
	log.Println("TriggerCMD run_command tool called.")
	url := "https://triggercmd.com/api/run/trigger"
	token, err := getTriggerCmdToken()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Missing TriggerCMD token: " + err.Error()}},
			IsError: true,
		}, nil, nil
	}

	payload := map[string]any{
		"command":  in.Command,
		"computer": in.Computer,
		"params":   in.Parameters,
		"sender":   "MCP",
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var result any
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println("Non-JSON response from run endpoint:", string(respBody))
		result = map[string]any{"status": resp.StatusCode, "ok": resp.StatusCode >= 200 && resp.StatusCode < 300, "body": string(respBody)}
	}
	statusText, _ := json.MarshalIndent(result, "", "  ")
	msg := fmt.Sprintf("Triggered '%s' on '%s'.\nResponse: %s", in.Command, in.Computer, string(statusText))
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: msg}}}, nil, nil
}

func main() {
	// Ensure logs go to stderr so stdout remains clean for MCP stdio protocol
	log.SetOutput(os.Stderr)
	log.Println("TriggerCMD MCP Server starting up...")
	server := mcp.NewServer(&mcp.Implementation{Name: "triggercmd", Version: "1.0.0"}, &mcp.ServerOptions{HasTools: true})

	// Register tools with typed handlers (schemas inferred from input structs)
	mcp.AddTool(server, &mcp.Tool{Name: "list_commands", Description: "List available TriggerCMD commands"}, listCommands)
	mcp.AddTool(server, &mcp.Tool{Name: "run_command", Description: "Run a TriggerCMD command by computer and command name with optional parameters"}, runCommand)

	// Run server over stdio transport
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
