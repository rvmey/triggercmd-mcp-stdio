#!/bin/bash
set -e

# TRIGGERcmd MCP Server Docker Build Script

echo "Building TRIGGERcmd MCP Server Docker Image..."

# Change to the parent directory (where go.mod and Dockerfile are located)
cd "$(dirname "$0")/.."

# Build the Docker image
docker build -t triggercmd-mcp:latest .

echo "Docker image built successfully!"
echo ""
echo "Usage examples:"
echo ""
echo "1. Run with environment variable token:"
echo "   docker run -e TRIGGERCMD_TOKEN='your-token-here' triggercmd-mcp:latest"
echo ""
echo "2. Run with token file mount:"
echo "   docker run -v ~/.TRIGGERcmdData/token.tkn:/home/triggercmd/.TRIGGERcmdData/token.tkn:ro triggercmd-mcp:latest"
echo ""
echo "3. Use with Docker Compose:"
echo "   cd docker && docker-compose up"
echo ""
echo "4. Test the image:"
echo "   docker run --rm -e TRIGGERCMD_TOKEN='your-token' triggercmd-mcp:latest"