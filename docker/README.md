# TRIGGERcmd MCP Server - Docker Setup

This directory contains Docker configuration files to build and run the TRIGGERcmd MCP Server as a container.

## Files

- `../Dockerfile` - Multi-stage build configuration for the MCP server (in repository root)
- `../dockerignore` - Files to exclude from the Docker build context (in repository root)
- `docker-compose.yml` - Docker Compose configuration with volume mounts and security settings
- `build.sh` - Build script to create the Docker image
- `README.md` - This documentation file

## Building the Image

### Option 1: Using the build script
```bash
cd docker
./build.sh
```

### Option 2: Manual build
```bash
# From the project root directory
docker build -t triggercmd-mcp:latest .
```

## Running the Container

### Option 1: With environment variable token
```bash
docker run -e TRIGGERCMD_TOKEN='your-token-here' triggercmd-mcp:latest
```

### Option 2: With token file mount (recommended)
```bash
docker run -v ~/.TRIGGERcmdData/token.tkn:/home/triggercmd/.TRIGGERcmdData/token.tkn:ro triggercmd-mcp:latest
```

## Configuration

### Authentication
The container supports two authentication methods:

1. **Environment Variable** (quick testing):
   ```bash
   docker run -e TRIGGERCMD_TOKEN='your-actual-token' triggercmd-mcp:latest
   ```

2. **Token File Mount** (recommended for production):
   ```bash
   docker run -v ~/.TRIGGERcmdData/token.tkn:/home/triggercmd/.TRIGGERcmdData/token.tkn:ro triggercmd-mcp:latest
   ```

### Docker Compose Configuration
Edit `docker-compose.yml` to:
- Uncomment the volume mount line
- Set the correct path to your token file
- Adjust resource limits if needed

## Security Features

The Docker setup includes several security best practices:
- **Non-root user**: Runs as user `triggercmd` (UID 1000)
- **Read-only filesystem**: Container filesystem is read-only
- **No new privileges**: Prevents privilege escalation
- **Minimal base image**: Uses Alpine Linux for smaller attack surface
- **No network access**: Uses `network_mode: none` since MCP uses stdio only

## Image Details

- **Base Image**: Alpine Linux 3.18 (minimal, secure)
- **Binary**: Compiled for Linux AMD64 
- **Size**: Approximately 15-20MB (small footprint)
- **User**: Non-root user `triggercmd`
- **Working Directory**: `/home/triggercmd`

## Usage with MCP Clients

When using with Claude Desktop or other MCP clients, configure the command to run the Docker container:

### Claude Desktop Example
```json
{
  "mcpServers": {
    "triggercmd": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-v", "/home/user/.TRIGGERcmdData/token.tkn:/home/triggercmd/.TRIGGERcmdData/token.tkn:ro",
        "triggercmd-mcp:latest"
      ]
    }
  }
}
```

### VS Code MCP Example
```json
{
  "servers": {
    "triggercmd": {
      "type": "stdio",
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-v", "/home/user/.TRIGGERcmdData/token.tkn:/home/triggercmd/.TRIGGERcmdData/token.tkn:ro",
        "triggercmd-mcp:latest"
      ]
    }
  }
}
```

## Troubleshooting

### Container won't start
- Check that the token file path is correct in the volume mount
- Verify the token file is readable by UID 1000
- Check Docker logs: `docker logs <container-name>`

### MCP client can't connect
- Ensure you're using `-i` (interactive) flag for stdin
- Don't use `-t` (tty) flag with MCP clients - it can interfere with the protocol
- Make sure the container has access to the token file

### Permission issues
- The container runs as UID 1000, ensure your token file is readable by this user
- Use `chmod 644 ~/.TRIGGERcmdData/token.tkn` if needed

## Development

To rebuild after code changes:
```bash
cd docker
./build.sh
```

The multi-stage build ensures a clean, minimal production image while keeping the build process efficient.