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
	"regexp"
	"strings"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listInput struct{}

type runInput struct {
	Command    string `json:"command" jsonschema:"Name of the command to run"`
	Computer   string `json:"computer" jsonschema:"Name of the computer that owns the command"`
	Parameters string `json:"parameters,omitempty" jsonschema:"Optional parameters to pass to the command"`
}

type dynamicCommandInput struct {
	Parameters string `json:"parameters,omitempty" jsonschema:"Optional parameters to pass to the command"`
}

// Command represents a TriggerCMD command from the API
type Command struct {
	Name               string   `json:"name"`
	Voice              string   `json:"voice"`
	McpToolDescription string   `json:"mcpToolDescription"`
	Computer           Computer `json:"computer"`
}

type Computer struct {
	Name string `json:"name"`
}

// CommandResponse represents the API response structure
type CommandResponse struct {
	Records []Command `json:"records"`
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

// fetchCommands retrieves all commands from the TriggerCMD API
func fetchCommands() ([]Command, error) {
	log.Println("Fetching commands from TriggerCMD API...")
	url := "https://triggercmd.com/api/command/list"
	token, err := getTriggerCmdToken()
	if err != nil {
		return nil, fmt.Errorf("missing TriggerCMD token: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "triggercmd-mcp-stdio/1.0.1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cmdResp CommandResponse
	if err := json.Unmarshal(body, &cmdResp); err != nil {
		log.Println("Failed to parse JSON. Raw response:", string(body))
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	log.Printf("Fetched %d commands from API", len(cmdResp.Records))

	// Filter commands that have mcpToolDescription
	var commandsWithMcp []Command
	for _, cmd := range cmdResp.Records {
		if strings.TrimSpace(cmd.McpToolDescription) != "" {
			commandsWithMcp = append(commandsWithMcp, cmd)
			log.Printf("Command '%s' on '%s' has MCP description: '%s'", cmd.Name, cmd.Computer.Name, cmd.McpToolDescription)
		}
	}

	log.Printf("Found %d commands with MCP tool descriptions", len(commandsWithMcp))
	return cmdResp.Records, nil
}

// generateToolName creates a valid tool name from computer and command names
func generateToolName(computerName, commandName string) string {
	// Convert to lowercase and replace non-alphanumeric chars with underscores
	toolName := fmt.Sprintf("run_%s_%s", computerName, commandName)
	toolName = strings.ToLower(toolName)

	// Replace non-alphanumeric characters with underscores
	reg := regexp.MustCompile(`[^a-z0-9_]`)
	toolName = reg.ReplaceAllString(toolName, "_")

	// Replace multiple underscores with single underscore
	reg = regexp.MustCompile(`_+`)
	toolName = reg.ReplaceAllString(toolName, "_")

	// Remove leading/trailing underscores
	toolName = strings.Trim(toolName, "_")

	return toolName
}

// createDynamicCommandHandler creates a handler for a specific command
func createDynamicCommandHandler(command Command) func(context.Context, *mcp.CallToolRequest, dynamicCommandInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input dynamicCommandInput) (*mcp.CallToolResult, any, error) {
		log.Printf("Dynamic command tool called: %s on %s", command.Name, command.Computer.Name)

		// Use the existing runCommand logic but with predefined computer and command
		runInput := runInput{
			Command:    command.Name,
			Computer:   command.Computer.Name,
			Parameters: input.Parameters,
		}

		return runCommand(ctx, nil, runInput)
	}
}

func listCommands(ctx context.Context, _ *mcp.CallToolRequest, _ listInput) (*mcp.CallToolResult, any, error) {
	log.Println("TriggerCMD list_commands tool called")
	commands, err := fetchCommands()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Error fetching commands: " + err.Error()}},
			IsError: true,
		}, nil, nil
	}

	// Build a simplified list: name, voice, computer.name, mcpToolDescription
	simplified := make([]map[string]any, 0)
	for _, cmd := range commands {
		simplified = append(simplified, map[string]any{
			"name":               cmd.Name,
			"voice":              cmd.Voice,
			"computer":           cmd.Computer.Name,
			"mcpToolDescription": cmd.McpToolDescription,
		})
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
	req.Header.Set("User-Agent", "triggercmd-mcp-stdio/1.0.1")
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

	server := mcp.NewServer(&mcp.Implementation{Name: "triggercmd", Version: "1.0.1"}, &mcp.ServerOptions{HasTools: true})

	// Register base tools
	mcp.AddTool(server, &mcp.Tool{Name: "list_commands", Description: "List available TriggerCMD commands"}, listCommands)
	mcp.AddTool(server, &mcp.Tool{Name: "run_command", Description: "Run a TriggerCMD command by computer and command name with optional parameters"}, runCommand)

	// Fetch commands and register dynamic tools for those with mcpToolDescription
	commands, err := fetchCommands()
	if err != nil {
		log.Printf("Warning: Failed to fetch commands for dynamic tools: %v", err)
	} else {
		dynamicToolsRegistered := 0
		for _, cmd := range commands {
			if strings.TrimSpace(cmd.McpToolDescription) != "" {
				toolName := generateToolName(cmd.Computer.Name, cmd.Name)
				description := strings.TrimSpace(cmd.McpToolDescription)

				log.Printf("Registering dynamic tool: %s -> %s on %s", toolName, cmd.Name, cmd.Computer.Name)

				// Create a tool with the custom description
				tool := &mcp.Tool{
					Name:        toolName,
					Description: description,
				}

				// Register the dynamic command handler
				handler := createDynamicCommandHandler(cmd)
				mcp.AddTool(server, tool, handler)
				dynamicToolsRegistered++
			}
		}
		log.Printf("Registered %d dynamic command tools", dynamicToolsRegistered)
	}

	// Run server over stdio transport
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
