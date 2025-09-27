# Testing Dynamic MCP Tools

## How It Works

When the TRIGGERcmd MCP server starts up, it:

1. **Fetches all commands** from the TRIGGERcmd API using your token
2. **Identifies commands with mcpToolDescription** - any command that has this field populated
3. **Creates individual MCP tools** for each such command with the pattern `run_{computer}_{command}`
4. **Uses the mcpToolDescription** as the tool description for better AI understanding

## Example Scenario

Suppose you have these commands in your TRIGGERcmd account:

**Computer: "office"**
- Command: "lights" 
  - mcpToolDescription: "Control office lights - accepts colors (red, blue, green, white) or on/off commands"
- Command: "music"
  - mcpToolDescription: "Play music in the office - accepts playlist names or 'stop'"
- Command: "backup" 
  - mcpToolDescription: "" (empty - will not create dynamic tool)

**Computer: "laptop"**  
- Command: "notepad"
  - mcpToolDescription: "Open notepad with optional filename parameter"

## Generated Dynamic Tools

The MCP server would automatically create these tools:

1. **run_office_lights**
   - Description: "Control office lights - accepts colors (red, blue, green, white) or on/off commands"
   - Parameters: `parameters` (optional string)

2. **run_office_music**  
   - Description: "Play music in the office - accepts playlist names or 'stop'"
   - Parameters: `parameters` (optional string)

3. **run_laptop_notepad**
   - Description: "Open notepad with optional filename parameter"  
   - Parameters: `parameters` (optional string)

The "backup" command would NOT get a dynamic tool since it has no mcpToolDescription.

## Benefits

- **Better AI Understanding**: Each tool has a specific, meaningful description
- **Easier Discovery**: AI can see exactly what each command does without generic descriptions
- **Cleaner Interface**: Direct tool names like `run_office_lights` vs generic `run_command`
- **Parameter Guidance**: Tool descriptions can explain what parameters are expected

## Current Status

✅ **Server Implementation**: Complete - fetches commands and registers dynamic tools
✅ **Startup Logging**: Shows how many dynamic tools were registered  
✅ **Tool Generation**: Creates sanitized tool names and uses custom descriptions
✅ **Error Handling**: Gracefully handles missing tokens or API errors

📋 **To Test**: Set mcpToolDescription fields on some of your TRIGGERcmd commands and restart the server

## Logs Example

When the server starts with dynamic tools enabled:

```
2025/09/27 19:02:51 TriggerCMD MCP Server starting up...
2025/09/27 19:02:51 Fetching commands from TriggerCMD API...
2025/09/27 19:02:52 Fetched 362 commands from API
2025/09/27 19:02:52 Command 'lights' on 'office' has MCP description: 'Control office lights...'
2025/09/27 19:02:52 Command 'music' on 'office' has MCP description: 'Play music in office...'
2025/09/27 19:02:52 Found 2 commands with MCP tool descriptions
2025/09/27 19:02:52 Registering dynamic tool: run_office_lights -> lights on office
2025/09/27 19:02:52 Registering dynamic tool: run_office_music -> music on office  
2025/09/27 19:02:52 Registered 2 dynamic command tools
```