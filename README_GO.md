# OpenCode Go Implementation

This is the Go implementation of OpenCode, an open-source AI coding agent.

## Project Structure

```
.
├── cmd/
│   └── opencode/          # Main application entry point
│       └── main.go
├── internal/
│   ├── agent/             # Agent definitions and registry
│   ├── cli/               # CLI commands and root
│   │   └── cmd/           # Individual command implementations
│   ├── config/            # Configuration management
│   ├── file/              # File operations and watching
│   ├── id/                # ID generation
│   ├── installation/      # Installation utilities
│   ├── lsp/               # LSP client implementation
│   ├── mcp/               # MCP client implementation
│   ├── permission/        # Permission system
│   ├── provider/          # AI provider implementations
│   │   ├── anthropic.go   # Anthropic/Claude provider
│   │   ├── openai.go      # OpenAI provider
│   │   └── openrouter.go  # OpenRouter provider
│   ├── server/            # HTTP server with SSE
│   ├── session/           # Session and message types
│   ├── storage/           # Persistent storage
│   ├── tool/              # Tool definitions and registry
│   └── util/
│       └── log/           # Logging utilities
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

## Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Or using go directly
go build -o bin/opencode ./cmd/opencode
```

## Running

```bash
# Run the CLI
./bin/opencode

# Run with a prompt
./bin/opencode run "Help me understand this codebase"

# Start the server
./bin/opencode serve --port 3000
```

## Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code (requires golangci-lint)
make lint

# Development with hot reload (requires air)
make dev
```

## Docker

```bash
# Build Docker image
make docker

# Or directly
docker build -t opencode:latest .

# Run container
docker run -it -v $(pwd):/workspace opencode:latest
```

## Configuration

OpenCode can be configured through:

1. **Global config**: `~/.config/opencode/opencode.json`
2. **Project config**: `./opencode.json`
3. **Environment variables**: `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, etc.

Example configuration:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "model": {
    "modelId": "claude-sonnet-4-20250514",
    "providerId": "anthropic"
  },
  "providers": {
    "anthropic": {
      "apiKey": "your-api-key"
    }
  },
  "mcp": {
    "filesystem": {
      "type": "stdio",
      "command": "mcp-filesystem",
      "args": ["--root", "/workspace"]
    }
  }
}
```

## Architecture

### Providers

The provider system supports multiple AI backends:

- **Anthropic**: Claude models (Sonnet, Opus, Haiku)
- **OpenAI**: GPT models (GPT-4o, o1, o3-mini)
- **OpenRouter**: Unified access to many models

### Tools

Built-in tools include:

- `read`: Read file contents
- `write`: Write to files
- `edit`: Edit files with string replacement
- `ls`: List directory contents
- `grep`: Search in files
- `glob`: Find files by pattern
- `bash`: Execute shell commands
- `webfetch`: Fetch web content
- `websearch`: Search the web
- `lsp`: LSP operations

### MCP Integration

Model Context Protocol (MCP) servers can be configured to extend functionality:

```json
{
  "mcp": {
    "github": {
      "type": "stdio",
      "command": "mcp-github",
      "env": {
        "GITHUB_TOKEN": "your-token"
      }
    }
  }
}
```

### LSP Support

The LSP client supports:

- Text document synchronization
- Hover information
- Go to definition
- Find references
- Completion
- Formatting

## License

MIT License - see LICENSE file for details.
