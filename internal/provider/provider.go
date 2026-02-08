package provider

import (
	"context"
	"time"
)

// Provider represents an AI provider
type Provider interface {
	// Name returns the provider name
	Name() string
	
	// Chat sends a chat request and returns the response
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	
	// Stream sends a chat request and streams the response
	Stream(ctx context.Context, req *ChatRequest, handler StreamHandler) error
	
	// Models returns available models
	Models() []Model
}

// Model represents an AI model
type Model struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Provider    string   `json:"provider"`
	ContextSize int      `json:"contextSize"`
	MaxOutput   int      `json:"maxOutput"`
	InputPrice  float64  `json:"inputPrice"`  // per 1M tokens
	OutputPrice float64  `json:"outputPrice"` // per 1M tokens
	Modalities  []string `json:"modalities"`  // text, image, audio
}

// ChatRequest represents a chat request
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []Message     `json:"messages"`
	Tools       []Tool        `json:"tools,omitempty"`
	MaxTokens   int           `json:"maxTokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        float64       `json:"topP,omitempty"`
	Stop        []string      `json:"stop,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string     `json:"role"` // system, user, assistant
	Content []Content  `json:"content"`
}

// Content represents message content
type Content struct {
	Type string `json:"type"` // text, image, tool_use, tool_result
	Text string `json:"text,omitempty"`
	
	// For images
	ImageURL string `json:"imageUrl,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	
	// For tool use
	ToolID   string      `json:"toolId,omitempty"`
	ToolName string      `json:"toolName,omitempty"`
	ToolArgs interface{} `json:"toolArgs,omitempty"`
	
	// For tool result
	ToolResult interface{} `json:"toolResult,omitempty"`
	IsError    bool        `json:"isError,omitempty"`
}

// Tool represents a tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	ID        string          `json:"id"`
	Model     string          `json:"model"`
	Created   time.Time       `json:"created"`
	Message   Message         `json:"message"`
	Usage     Usage           `json:"usage"`
	ToolCalls []ToolCall      `json:"toolCalls,omitempty"`
	FinishReason string       `json:"finishReason"`
}

// Usage represents token usage
type Usage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
	TotalTokens  int `json:"totalTokens"`
}

// ToolCall represents a tool call
type ToolCall struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Arguments interface{} `json:"arguments"`
}

// StreamHandler handles streaming responses
type StreamHandler func(event StreamEvent) error

// StreamEvent represents a streaming event
type StreamEvent struct {
	Type      string      `json:"type"` // content, tool_use, done, error
	Content   string      `json:"content,omitempty"`
	ToolCall  *ToolCall   `json:"toolCall,omitempty"`
	Delta     *Content    `json:"delta,omitempty"`
	Error     error       `json:"error,omitempty"`
	Usage     *Usage      `json:"usage,omitempty"`
}

// Info represents provider configuration info
type Info struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	EnvVar   string            `json:"envVar"`   // Environment variable for API key
	BaseURL  string            `json:"baseUrl"`
	Headers  map[string]string `json:"headers"`
	Disabled bool              `json:"disabled"`
}

// Registry holds all registered providers
type Registry struct {
	providers map[string]Provider
	models    map[string]Model
}

// NewRegistry creates a new provider registry
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
		models:    make(map[string]Model),
	}
}

// Register registers a provider
func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
	for _, m := range p.Models() {
		r.models[m.ID] = m
	}
}

// Get retrieves a provider by name
func (r *Registry) Get(name string) (Provider, bool) {
	p, ok := r.providers[name]
	return p, ok
}

// GetModel retrieves a model by ID
func (r *Registry) GetModel(id string) (Model, bool) {
	m, ok := r.models[id]
	return m, ok
}

// ListModels returns all available models
func (r *Registry) ListModels() []Model {
	models := make([]Model, 0, len(r.models))
	for _, m := range r.models {
		models = append(models, m)
	}
	return models
}
