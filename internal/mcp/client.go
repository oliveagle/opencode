package mcp

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

// Client represents an MCP client
type Client struct {
	name    string
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader

	mu          sync.Mutex
	nextID      int
	pending     map[int]chan *Response
	tools       []Tool
	resources   []Resource
	prompts     []Prompt
	initialized bool
}

// Response represents an MCP response
type Response struct {
	ID     int             `json:"id,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *Error          `json:"error,omitempty"`
}

// Error represents an MCP error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// Resource represents an MCP resource
type Resource struct {
	URI         string      `json:"uri"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	MimeType    string      `json:"mimeType,omitempty"`
}

// Prompt represents an MCP prompt
type Prompt struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Arguments   []Argument  `json:"arguments,omitempty"`
}

// Argument represents a prompt argument
type Argument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// ToolResult represents the result of a tool call
type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content represents content in a tool result
type Content struct {
	Type string `json:"type"` // text, image, resource
	Text string `json:"text,omitempty"`
	Data string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// NewClient creates a new MCP client
func NewClient(name string, command string, args []string, env map[string]string) *Client {
	cmd := exec.Command(command, args...)
	
	// Set environment
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	return &Client{
		name:    name,
		cmd:     cmd,
		pending: make(map[int]chan *Response),
	}
}

// Start starts the MCP server
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
			fmt.Fprintf(os.Stderr, "[MCP %s stderr] %s\n", c.name, scanner.Text())
		}
	}()

	return nil
}

// Stop stops the MCP server
func (c *Client) Stop() error {
	if c.cmd.Process != nil {
		return c.cmd.Process.Kill()
	}
	return nil
}

// Request sends an MCP request
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
			return nil, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	}
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
			Error   *Error          `json:"error,omitempty"`
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
		}
		// MCP notifications can be handled here if needed
	}
}

// Initialize initializes the MCP connection
func (c *Client) Initialize(ctx context.Context) error {
	params := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"clientInfo": map[string]interface{}{
			"name":    "opencode",
			"version": "1.0.0",
		},
	}

	result, err := c.Request(ctx, "initialize", params)
	if err != nil {
		return err
	}

	var initResult struct {
		Capabilities struct {
			Tools     interface{} `json:"tools,omitempty"`
			Resources interface{} `json:"resources,omitempty"`
			Prompts   interface{} `json:"prompts,omitempty"`
		} `json:"capabilities"`
	}
	if err := json.Unmarshal(result, &initResult); err != nil {
		return err
	}

	c.initialized = true

	// Send initialized notification
	initialized := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	}
	if err := c.writeMessage(initialized); err != nil {
		return err
	}

	// List tools
	if initResult.Capabilities.Tools != nil {
		if err := c.listTools(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) listTools(ctx context.Context) error {
	result, err := c.Request(ctx, "tools/list", map[string]interface{}{})
	if err != nil {
		return err
	}

	var listResult struct {
		Tools []Tool `json:"tools"`
	}
	if err := json.Unmarshal(result, &listResult); err != nil {
		return err
	}

	c.tools = listResult.Tools
	return nil
}

// CallTool calls an MCP tool
func (c *Client) CallTool(ctx context.Context, name string, arguments map[string]interface{}) (*ToolResult, error) {
	params := map[string]interface{}{
		"name":      name,
		"arguments": arguments,
	}

	result, err := c.Request(ctx, "tools/call", params)
	if err != nil {
		return nil, err
	}

	var toolResult ToolResult
	if err := json.Unmarshal(result, &toolResult); err != nil {
		return nil, err
	}

	return &toolResult, nil
}

// Tools returns available tools
func (c *Client) Tools() []Tool {
	return c.tools
}

// Resources returns available resources
func (c *Client) Resources() []Resource {
	return c.resources
}

// Prompts returns available prompts
func (c *Client) Prompts() []Prompt {
	return c.prompts
}

// Name returns the server name
func (c *Client) Name() string {
	return c.name
}

// Registry manages multiple MCP servers
type Registry struct {
	servers map[string]*Client
}

// NewRegistry creates a new MCP registry
func NewRegistry() *Registry {
	return &Registry{
		servers: make(map[string]*Client),
	}
}

// Register registers an MCP server
func (r *Registry) Register(name string, client *Client) {
	r.servers[name] = client
}

// Get retrieves an MCP server by name
func (r *Registry) Get(name string) (*Client, bool) {
	client, ok := r.servers[name]
	return client, ok
}

// List returns all registered servers
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.servers))
	for name := range r.servers {
		names = append(names, name)
	}
	return names
}

// AllTools returns tools from all servers
func (r *Registry) AllTools() map[string][]Tool {
	tools := make(map[string][]Tool)
	for name, client := range r.servers {
		tools[name] = client.Tools()
	}
	return tools
}

// StartAll starts all MCP servers
func (r *Registry) StartAll(ctx context.Context) error {
	for name, client := range r.servers {
		if err := client.Start(ctx); err != nil {
			return fmt.Errorf("failed to start MCP server %s: %w", name, err)
		}
		if err := client.Initialize(ctx); err != nil {
			return fmt.Errorf("failed to initialize MCP server %s: %w", name, err)
		}
	}
	return nil
}

// StopAll stops all MCP servers
func (r *Registry) StopAll() {
	for _, client := range r.servers {
		client.Stop()
	}
}

// CallTool calls a tool on a specific server
func (r *Registry) CallTool(ctx context.Context, serverName, toolName string, args map[string]interface{}) (*ToolResult, error) {
	client, ok := r.servers[serverName]
	if !ok {
		return nil, fmt.Errorf("MCP server %s not found", serverName)
	}
	return client.CallTool(ctx, toolName, args)
}

// FindTool finds which server has a tool
func (r *Registry) FindTool(toolName string) (string, *Tool, bool) {
	for serverName, client := range r.servers {
		for _, tool := range client.Tools() {
			if tool.Name == toolName {
				return serverName, &tool, true
			}
		}
	}
	return "", nil, false
}
