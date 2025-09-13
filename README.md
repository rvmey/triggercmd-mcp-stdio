
# triggercmd-mcp-stdio

Local stdio MCP server for TRIGGERcmd

![Claude Desktop using TRIGGERcmd MCP server](Claude_Desktop_using_TRIGGERcmd_MCP_server.png)

## Authentication

If the TRIGGERCMD_TOKEN environment variable is set, it will use that, otherwise it will try to read the token from your ~/.TRIGGERcmdData/token.tkn file in your user's home directory.

## AI prompt examples

```
Use the #triggercmd mcp server to run notepad on laptop with parameter newnote.
```

You can also tell the LLM what the command does, and how the parameters work before you run your command. For example:
```
Prompt 1: The parameter of the "lights" command on the "office" computer in the #triggercmd mcp server is used to specify what color (red, white, green, or blue) to make the lights in my office. It also accepts on or off to turn the lights on or off.

Prompt 2: Turn the lights in my office green.

Prompt 3: Turn the lights in my office off.
```

## Usage

If using Claude Desktop, add an entry like this to your claude_desktop_config.json file: 
```
{
  "mcpServers": {
    "triggercmd": {
      "command": "C:\\Users\\johndoe\\Downloads\\triggercmd-mcp-windows-amd64.exe"
    }
  }
}
```

If using VS Code, add an entry like this to your mcp.json file:
```
{
	"servers": {
		"triggercmd": {
			"type": "stdio",
			"command": "C:\\Users\\johndoe\\Downloads\\triggercmd-mcp-windows-amd64.exe"
		}
	}
}
```
## Downloads

Mac:

[darwin-arm64](https://agents.triggercmd.com/triggercmd-mcp/triggercmd-mcp-darwin-arm64)

[darwin-amd64](https://agents.triggercmd.com/triggercmd-mcp/triggercmd-mcp-darwin-amd64)

Linux:

[linux-386](https://agents.triggercmd.com/triggercmd-mcp/triggercmd-mcp-linux-386)

[linux-amd64](https://agents.triggercmd.com/triggercmd-mcp/triggercmd-mcp-linux-amd64)

[linux-arm](https://agents.triggercmd.com/triggercmd-mcp/triggercmd-mcp-linux-arm)

[linux-arm64](https://agents.triggercmd.com/triggercmd-mcp/triggercmd-mcp-linux-arm64)

Windows:

[windows-386.exe](https://agents.triggercmd.com/triggercmd-mcp/triggercmd-mcp-windows-386.exe)

[windows-amd64.exe](https://agents.triggercmd.com/triggercmd-mcp/triggercmd-mcp-windows-amd64.exe)

[windows-arm64.exe](https://agents.triggercmd.com/triggercmd-mcp/triggercmd-mcp-windows-arm64.exe)

NOTE: On Mac and Linux, you'll need to make the binary executable with a command like this:
```
chmod +x triggercmd-mcp-darwin-arm64
```