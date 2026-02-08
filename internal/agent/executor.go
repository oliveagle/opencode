package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/anomalyco/opencode/internal/provider"
	"github.com/anomalyco/opencode/internal/session"
	"github.com/anomalyco/opencode/internal/tool"
	"github.com/anomalyco/opencode/internal/util/log"
)

// Executor handles agent execution
type Executor struct {
	agent       *Info
	provider    provider.Provider
	model       provider.Model
	toolRegistry *tool.Registry
	session     *session.Info
	conversation *session.Conversation
	mu          sync.Mutex
}

// ExecutorConfig configures the executor
type ExecutorConfig struct {
	Agent       *Info
	Provider    provider.Provider
	Model       provider.Model
	ToolRegistry *tool.Registry
	Session     *session.Info
}

// NewExecutor creates a new agent executor
func NewExecutor(cfg ExecutorConfig) *Executor {
	return &Executor{
		agent:       cfg.Agent,
		provider:    cfg.Provider,
		model:       cfg.Model,
		toolRegistry: cfg.ToolRegistry,
		session:     cfg.Session,
		conversation: &session.Conversation{
			SessionID: cfg.Session.ID.String(),
			Messages:  []*session.Message{},
		},
	}
}

// ExecuteResult represents the result of agent execution
type ExecuteResult struct {
	Message     *session.Message `json:"message"`
	ToolCalls   []provider.ToolCall `json:"toolCalls,omitempty"`
	Finished    bool             `json:"finished"`
	FinishReason string           `json:"finishReason,omitempty"`
}

// Execute executes the agent with a user message
func (e *Executor) Execute(ctx context.Context, userMessage string) (*ExecuteResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Add user message to conversation
	msg := session.NewMessage("user", []session.Part{
		session.TextPart(userMessage),
	})
	e.conversation.AddMessage(msg)

	// Build request
	req := &provider.ChatRequest{
		Model:     e.model.ID,
		Messages:  e.conversation.GetMessagesForAPI(),
		MaxTokens: e.model.MaxOutput,
		Tools:     e.getTools(),
	}

	// Call provider
	resp, err := e.provider.Chat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("provider error: %w", err)
	}

	// Add assistant message to conversation
	assistantMsg := session.NewMessage("assistant", e.convertContentToParts(resp.Message.Content))
	assistantMsg.ID = resp.ID
	e.conversation.AddMessage(assistantMsg)

	return &ExecuteResult{
		Message:     assistantMsg,
		ToolCalls:   resp.ToolCalls,
		Finished:    resp.FinishReason == "stop" || resp.FinishReason == "end_turn",
		FinishReason: resp.FinishReason,
	}, nil
}

// ExecuteWithTools executes and handles tool calls
func (e *Executor) ExecuteWithTools(ctx context.Context, userMessage string, maxSteps int) (*ExecuteResult, error) {
	result, err := e.Execute(ctx, userMessage)
	if err != nil {
		return nil, err
	}

	steps := 0
	for len(result.ToolCalls) > 0 && steps < maxSteps {
		steps++

		// Execute tool calls
		toolResults := make([]session.Part, 0)
		for _, tc := range result.ToolCalls {
			toolResult, err := e.executeTool(ctx, tc)
			if err != nil {
				toolResults = append(toolResults, session.ToolResultPart(tc.ID, err.Error(), true))
			} else {
				toolResults = append(toolResults, toolResult)
			}
		}

		// Add tool results to conversation
		toolMsg := session.NewMessage("user", toolResults)
		e.conversation.AddMessage(toolMsg)

		// Continue execution
		req := &provider.ChatRequest{
			Model:     e.model.ID,
			Messages:  e.conversation.GetMessagesForAPI(),
			MaxTokens: e.model.MaxOutput,
			Tools:     e.getTools(),
		}

		resp, err := e.provider.Chat(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("provider error in step %d: %w", steps, err)
		}

		assistantMsg := session.NewMessage("assistant", e.convertContentToParts(resp.Message.Content))
		assistantMsg.ID = resp.ID
		e.conversation.AddMessage(assistantMsg)

		result = &ExecuteResult{
			Message:     assistantMsg,
			ToolCalls:   resp.ToolCalls,
			Finished:    resp.FinishReason == "stop" || resp.FinishReason == "end_turn",
			FinishReason: resp.FinishReason,
		}
	}

	return result, nil
}

// Stream executes with streaming response
func (e *Executor) Stream(ctx context.Context, userMessage string, handler StreamHandler) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Add user message
	msg := session.NewMessage("user", []session.Part{
		session.TextPart(userMessage),
	})
	e.conversation.AddMessage(msg)

	// Build request
	req := &provider.ChatRequest{
		Model:     e.model.ID,
		Messages:  e.conversation.GetMessagesForAPI(),
		MaxTokens: e.model.MaxOutput,
		Tools:     e.getTools(),
	}

	// Stream from provider
	var fullContent string
	var toolCalls []provider.ToolCall

	err := e.provider.Stream(ctx, req, func(event provider.StreamEvent) error {
		switch event.Type {
		case "content":
			fullContent += event.Content
			handler(StreamEvent{Type: "content", Content: event.Content})
		case "tool_use":
			if event.ToolCall != nil {
				toolCalls = append(toolCalls, *event.ToolCall)
				handler(StreamEvent{Type: "tool_use", ToolCall: event.ToolCall})
			}
		case "done":
			handler(StreamEvent{Type: "done"})
		case "usage":
			handler(StreamEvent{Type: "usage", Usage: event.Usage})
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Add assistant message
	parts := []session.Part{session.TextPart(fullContent)}
	for _, tc := range toolCalls {
		parts = append(parts, session.ToolUsePart(tc.ID, tc.Name, tc.Arguments))
	}
	assistantMsg := session.NewMessage("assistant", parts)
	e.conversation.AddMessage(assistantMsg)

	return nil
}

// StreamHandler handles streaming events
type StreamHandler func(event StreamEvent)

// StreamEvent represents a streaming event
type StreamEvent struct {
	Type     string               `json:"type"`
	Content  string               `json:"content,omitempty"`
	ToolCall *provider.ToolCall   `json:"toolCall,omitempty"`
	Usage    *provider.Usage      `json:"usage,omitempty"`
}

func (e *Executor) getTools() []provider.Tool {
	tools := make([]provider.Tool, 0)
	
	// Get tools from registry
	for _, id := range e.toolRegistry.List() {
		t, ok := e.toolRegistry.Get(id)
		if !ok {
			continue
		}
		
		info, err := t.Info(context.Background())
		if err != nil {
			continue
		}

		// Check permission
		if e.agent.Permission != nil && !e.agent.Permission.IsAllowed(id) {
			continue
		}

		tools = append(tools, provider.Tool{
			Name:        info.ID,
			Description: info.Description,
			Parameters:  info.Parameters,
		})
	}

	return tools
}

func (e *Executor) executeTool(ctx context.Context, tc provider.ToolCall) (session.Part, error) {
	log.Default.Info("Executing tool", map[string]interface{}{
		"tool": tc.Name,
		"args": tc.Arguments,
	})

	t, ok := e.toolRegistry.Get(tc.Name)
	if !ok {
		return session.Part{}, fmt.Errorf("tool %s not found", tc.Name)
	}

	// Convert arguments
	args, ok := tc.Arguments.(map[string]interface{})
	if !ok {
		// Try to convert from JSON
		if jsonBytes, err := json.Marshal(tc.Arguments); err == nil {
			json.Unmarshal(jsonBytes, &args)
		}
	}
	if args == nil {
		args = make(map[string]interface{})
	}

	// Execute tool
	result, err := t.Execute(ctx, args, &tool.Context{
		SessionID: e.session.ID.String(),
		Agent:     e.agent.Name,
	})
	if err != nil {
		return session.Part{}, err
	}

	log.Default.Info("Tool execution complete", map[string]interface{}{
		"tool":   tc.Name,
		"title":  result.Title,
		"outputLen": len(result.Output),
	})

	return session.ToolResultPart(tc.ID, result.Output, false), nil
}

func (e *Executor) convertContentToParts(content []provider.Content) []session.Part {
	parts := make([]session.Part, 0, len(content))
	for _, c := range content {
		switch c.Type {
		case "text":
			parts = append(parts, session.TextPart(c.Text))
		case "tool_use":
			parts = append(parts, session.ToolUsePart(c.ToolID, c.ToolName, c.ToolArgs))
		case "tool_result":
			parts = append(parts, session.ToolResultPart(c.ToolID, c.ToolResult, c.IsError))
		case "image":
			parts = append(parts, session.ImagePart(c.ImageURL, c.MimeType))
		}
	}
	return parts
}

// Conversation returns the current conversation
func (e *Executor) Conversation() *session.Conversation {
	return e.conversation
}

// SetSystemPrompt sets a system prompt
func (e *Executor) SetSystemPrompt(prompt string) {
	// Add system message at the beginning
	systemMsg := session.NewMessage("system", []session.Part{
		session.TextPart(prompt),
	})
	e.conversation.Messages = append([]*session.Message{systemMsg}, e.conversation.Messages...)
}

// MaxSteps returns the maximum number of tool-calling steps
func (e *Executor) MaxSteps() int {
	if e.agent.Steps > 0 {
		return e.agent.Steps
	}
	return 50 // Default max steps
}

// ShouldContinue determines if execution should continue
func (e *Executor) ShouldContinue(result *ExecuteResult, steps int) bool {
	if steps >= e.MaxSteps() {
		return false
	}
	if result.Finished {
		return false
	}
	if len(result.ToolCalls) == 0 {
		return false
	}
	return true
}

// Run runs the complete agent loop
func (e *Executor) Run(ctx context.Context, userMessage string) (*session.Message, error) {
	result, err := e.ExecuteWithTools(ctx, userMessage, e.MaxSteps())
	if err != nil {
		return nil, err
	}
	return result.Message, nil
}

// RunAsync runs the agent asynchronously
func (e *Executor) RunAsync(ctx context.Context, userMessage string, handler func(*ExecuteResult, error)) {
	go func() {
		result, err := e.ExecuteWithTools(ctx, userMessage, e.MaxSteps())
		handler(result, err)
	}()
}

// Abort aborts the current execution
func (e *Executor) Abort() {
	// TODO: Implement cancellation
}

// Status returns the current execution status
func (e *Executor) Status() Status {
	return Status{
		Running:     false, // TODO: Track running state
		Steps:       0,     // TODO: Track steps
		LastUpdate:  time.Now(),
	}
}

// Status represents executor status
type Status struct {
	Running     bool      `json:"running"`
	Steps       int       `json:"steps"`
	LastUpdate  time.Time `json:"lastUpdate"`
}
