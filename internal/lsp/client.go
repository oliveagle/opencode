package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Client represents an LSP client
type Client struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.Reader
	stderr io.Reader

	mu       sync.Mutex
	nextID   int
	pending  map[int]chan *Response
	handlers map[string]NotificationHandler

	initialized bool
	capabilities *ServerCapabilities
}

// Response represents an LSP response
type Response struct {
	ID     int             `json:"id,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *ResponseError  `json:"error,omitempty"`
}

// ResponseError represents an LSP error
type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NotificationHandler handles LSP notifications
type NotificationHandler func(params json.RawMessage)

// ServerCapabilities represents server capabilities
type ServerCapabilities struct {
	TextDocumentSync           interface{} `json:"textDocumentSync,omitempty"`
	HoverProvider              bool        `json:"hoverProvider,omitempty"`
	CompletionProvider         interface{} `json:"completionProvider,omitempty"`
	SignatureHelpProvider      interface{} `json:"signatureHelpProvider,omitempty"`
	DefinitionProvider         bool        `json:"definitionProvider,omitempty"`
	TypeDefinitionProvider     bool        `json:"typeDefinitionProvider,omitempty"`
	ImplementationProvider     bool        `json:"implementationProvider,omitempty"`
	ReferencesProvider         bool        `json:"referencesProvider,omitempty"`
	DocumentHighlightProvider  bool        `json:"documentHighlightProvider,omitempty"`
	DocumentSymbolProvider     bool        `json:"documentSymbolProvider,omitempty"`
	WorkspaceSymbolProvider    bool        `json:"workspaceSymbolProvider,omitempty"`
	CodeActionProvider         interface{} `json:"codeActionProvider,omitempty"`
	CodeLensProvider           interface{} `json:"codeLensProvider,omitempty"`
	DocumentFormattingProvider bool        `json:"documentFormattingProvider,omitempty"`
	DocumentRangeFormattingProvider bool `json:"documentRangeFormattingProvider,omitempty"`
	RenameProvider             interface{} `json:"renameProvider,omitempty"`
	ExecuteCommandProvider     interface{} `json:"executeCommandProvider,omitempty"`
}

// NewClient creates a new LSP client
func NewClient(command string, args ...string) *Client {
	return &Client{
		cmd:      exec.Command(command, args...),
		pending:  make(map[int]chan *Response),
		handlers: make(map[string]NotificationHandler),
	}
}

// Start starts the LSP server
func (c *Client) Start(ctx context.Context) error {
	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return err
	}
	c.stdin = stdin

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	c.stdout = stdout

	stderr, err := c.cmd.StderrPipe()
	if err != nil {
		return err
	}
	c.stderr = stderr

	if err := c.cmd.Start(); err != nil {
		return err
	}

	// Start reading responses
	go c.readLoop()

	// Log stderr
	go func() {
		scanner := bufio.NewScanner(c.stderr)
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "[LSP stderr] %s\n", scanner.Text())
		}
	}()

	return nil
}

// Stop stops the LSP server
func (c *Client) Stop() error {
	if c.cmd.Process != nil {
		// Send shutdown request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c.Shutdown(ctx)
		c.Exit(ctx)
		return c.cmd.Process.Kill()
	}
	return nil
}

// Request sends an LSP request
func (c *Client) Request(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	id := c.nextID
	c.nextID++
	ch := make(chan *Response, 1)
	c.pending[id] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
	}()

	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
		"params":  params,
	}

	if err := c.writeMessage(req); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-ch:
		if resp.Error != nil {
			return nil, fmt.Errorf("LSP error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	}
}

// Notify sends an LSP notification
func (c *Client) Notify(method string, params interface{}) error {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}
	return c.writeMessage(req)
}

// OnNotification registers a notification handler
func (c *Client) OnNotification(method string, handler NotificationHandler) {
	c.handlers[method] = handler
}

func (c *Client) writeMessage(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return err
	}
	_, err = c.stdin.Write(data)
	return err
}

func (c *Client) readLoop() {
	reader := bufio.NewReader(c.stdout)
	for {
		// Read Content-Length header
		var length int
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if line == "\r\n" {
				break
			}
			if len(line) > 16 && line[:16] == "Content-Length: " {
				fmt.Sscanf(line, "Content-Length: %d\r\n", &length)
			}
		}

		// Read body
		body := make([]byte, length)
		if _, err := io.ReadFull(reader, body); err != nil {
			return
		}

		// Parse message
		var msg struct {
			Jsonrpc string          `json:"jsonrpc"`
			ID      *int            `json:"id,omitempty"`
			Method  string          `json:"method,omitempty"`
			Params  json.RawMessage `json:"params,omitempty"`
			Result  json.RawMessage `json:"result,omitempty"`
			Error   *ResponseError  `json:"error,omitempty"`
		}
		if err := json.Unmarshal(body, &msg); err != nil {
			continue
		}

		if msg.ID != nil {
			// Response
			c.mu.Lock()
			if ch, ok := c.pending[*msg.ID]; ok {
				ch <- &Response{
					ID:     *msg.ID,
					Result: msg.Result,
					Error:  msg.Error,
				}
			}
			c.mu.Unlock()
		} else if msg.Method != "" {
			// Notification
			if handler, ok := c.handlers[msg.Method]; ok {
				handler(msg.Params)
			}
		}
	}
}

// Initialize sends the initialize request
func (c *Client) Initialize(ctx context.Context, rootURI string) error {
	params := map[string]interface{}{
		"processId":    os.Getpid(),
		"rootUri":      rootURI,
		"capabilities": map[string]interface{}{},
	}

	result, err := c.Request(ctx, "initialize", params)
	if err != nil {
		return err
	}

	var initResult struct {
		Capabilities ServerCapabilities `json:"capabilities"`
	}
	if err := json.Unmarshal(result, &initResult); err != nil {
		return err
	}

	c.capabilities = &initResult.Capabilities
	c.initialized = true

	// Send initialized notification
	return c.Notify("initialized", map[string]interface{}{})
}

// Shutdown sends the shutdown request
func (c *Client) Shutdown(ctx context.Context) error {
	if !c.initialized {
		return nil
	}
	_, err := c.Request(ctx, "shutdown", nil)
	return err
}

// Exit sends the exit notification
func (c *Client) Exit(ctx context.Context) error {
	return c.Notify("exit", nil)
}

// DidOpen notifies the server that a document was opened
func (c *Client) DidOpen(uri, languageID, content string) error {
	return c.Notify("textDocument/didOpen", map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri":        uri,
			"languageId": languageID,
			"version":    1,
			"text":       content,
		},
	})
}

// DidChange notifies the server that a document was changed
func (c *Client) DidChange(uri string, version int, content string) error {
	return c.Notify("textDocument/didChange", map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri":     uri,
			"version": version,
		},
		"contentChanges": []map[string]interface{}{
			{"text": content},
		},
	})
}

// DidClose notifies the server that a document was closed
func (c *Client) DidClose(uri string) error {
	return c.Notify("textDocument/didClose", map[string]interface{}{
		"textDocument": map[string]interface{}{
			"uri": uri,
		},
	})
}

// Definition requests the definition of a symbol
func (c *Client) Definition(ctx context.Context, uri string, line, character int) (json.RawMessage, error) {
	return c.Request(ctx, "textDocument/definition", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"position":     map[string]interface{}{"line": line, "character": character},
	})
}

// Hover requests hover information
func (c *Client) Hover(ctx context.Context, uri string, line, character int) (json.RawMessage, error) {
	return c.Request(ctx, "textDocument/hover", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"position":     map[string]interface{}{"line": line, "character": character},
	})
}

// Completion requests completion items
func (c *Client) Completion(ctx context.Context, uri string, line, character int) (json.RawMessage, error) {
	return c.Request(ctx, "textDocument/completion", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"position":     map[string]interface{}{"line": line, "character": character},
	})
}

// References requests all references to a symbol
func (c *Client) References(ctx context.Context, uri string, line, character int, includeDeclaration bool) (json.RawMessage, error) {
	return c.Request(ctx, "textDocument/references", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"position":     map[string]interface{}{"line": line, "character": character},
		"context": map[string]interface{}{
			"includeDeclaration": includeDeclaration,
		},
	})
}

// Formatting requests document formatting
func (c *Client) Formatting(ctx context.Context, uri string, options map[string]interface{}) (json.RawMessage, error) {
	return c.Request(ctx, "textDocument/formatting", map[string]interface{}{
		"textDocument": map[string]interface{}{"uri": uri},
		"options":      options,
	})
}
