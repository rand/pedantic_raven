# Pedantic Raven

**Interactive Context Engineering Environment**

A next-generation TUI for creating, editing, and refining context documents and project knowledge, built with Go and Bubble Tea.

## Features

### Current (Phase 1 - Foundation)
- ✅ Bubble Tea framework setup
- 🚧 Rich interactive editor
- 🚧 File tree navigation
- 🚧 Command palette

### Planned
- **Direct LLM Integration**: Multi-provider support (Anthropic, OpenAI, Gemini, local models)
- **Semantic Analysis**: Entity extraction, relationship mapping, hole detection
- **Knowledge Graph Visualization**: Interactive force-directed graphs
- **Mnemosyne Integration**: Optional RPC-based memory and orchestration
- **Multi-Buffer Editing**: Work with multiple files simultaneously
- **Extension System**: Plugin API, LSP integration, custom analyzers

## Architecture

- **Language**: Go
- **TUI Framework**: Bubble Tea (Elm Architecture)
- **Styling**: Lipgloss
- **Components**: Bubbles
- **Integration**: gRPC to mnemosyne-rpc (optional)

## Development

```bash
# Run the app
go run main.go

# Build
go build -o pedantic-raven

# Run tests
go test ./...
```

## Project Structure

```
pedantic_raven/
├── main.go                 # Entry point
├── internal/
│   ├── app/               # Application core
│   ├── editor/            # Editor component
│   ├── tree/              # File tree component
│   ├── llm/               # LLM providers
│   ├── graph/             # Knowledge graph
│   └── mnemosyne/         # gRPC client (optional)
├── proto/                 # Protobuf definitions (copied from mnemosyne)
└── README.md
```

## Integration with Mnemosyne

Pedantic Raven works standalone but enhances when mnemosyne-rpc is available:

- **Level 1**: Memory operations (store, recall, search)
- **Level 2**: Deep semantic analysis (LLM + memory context)
- **Level 3**: Multi-agent orchestration (bidirectional event streaming)

See [Mnemosyne RPC Documentation](../mnemosyne-rpc-dev/docs/rpc.md) for setup.

## License

MIT
