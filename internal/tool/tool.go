package tool

import (
	"context"
	"encoding/json"
)

// Info represents a tool definition
type Info struct {
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Parameters  *Schema     `json:"parameters"`
}

// Schema represents a JSON schema for tool parameters
type Schema struct {
	Type       string                 `json:"type"`
	Properties map[string]Property    `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

// Property represents a schema property
type Property struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Items       *Property   `json:"items,omitempty"` // For arrays
}

// Context provides context for tool execution
type Context struct {
	SessionID string
	MessageID string
	Agent     string
	CallID    string
	Extra     map[string]interface{}
}

// Result represents the result of tool execution
type Result struct {
	Title       string      `json:"title"`
	Output      string      `json:"output"`
	Metadata    interface{} `json:"metadata,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents a file attachment
type Attachment struct {
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	Content  []byte `json:"content"`
}

// Executor is the interface for tool execution
type Executor interface {
	// Info returns tool information
	Info(ctx context.Context) (*Info, error)
	
	// Execute runs the tool
	Execute(ctx context.Context, args map[string]interface{}, execCtx *Context) (*Result, error)
	
	// Validate validates the arguments
	Validate(args map[string]interface{}) error
}

// Registry holds all registered tools
type Registry struct {
	tools map[string]Executor
}

// NewRegistry creates a new tool registry
func NewRegistry() *Registry {
	r := &Registry{
		tools: make(map[string]Executor),
	}
	r.registerDefaults()
	return r
}

func (r *Registry) registerDefaults() {
	// Register built-in tools
	r.Register("read", &ReadTool{})
	r.Register("write", &WriteTool{})
	r.Register("edit", &EditTool{})
	r.Register("ls", &LsTool{})
	r.Register("grep", &GrepTool{})
	r.Register("glob", &GlobTool{})
	r.Register("bash", &BashTool{})
	r.Register("webfetch", &WebFetchTool{})
	r.Register("websearch", &WebSearchTool{})
	r.Register("question", &QuestionTool{})
	r.Register("task", &TaskTool{})
	r.Register("todo", &TodoTool{})
	r.Register("lsp", &LspTool{})
	r.Register("apply_patch", &ApplyPatchTool{})
}

// Register registers a tool
func (r *Registry) Register(id string, tool Executor) {
	r.tools[id] = tool
}

// Get retrieves a tool by ID
func (r *Registry) Get(id string) (Executor, bool) {
	tool, ok := r.tools[id]
	return tool, ok
}

// List returns all registered tools
func (r *Registry) List() []string {
	ids := make([]string, 0, len(r.tools))
	for id := range r.tools {
		ids = append(ids, id)
	}
	return ids
}

// ListInfo returns information about all tools
func (r *Registry) ListInfo(ctx context.Context) ([]*Info, error) {
	infos := make([]*Info, 0, len(r.tools))
	for _, tool := range r.tools {
		info, err := tool.Info(ctx)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// ToJSONSchema converts tool info to JSON schema format for AI providers
func (info *Info) ToJSONSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        info.ID,
			"description": info.Description,
			"parameters":  info.Parameters,
		},
	}
}

// ToJSON returns JSON representation
func (info *Info) ToJSON() string {
	b, _ := json.Marshal(info)
	return string(b)
}
