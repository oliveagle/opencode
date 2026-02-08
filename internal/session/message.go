package session

import (
	"time"

	"github.com/anomalyco/opencode/internal/id"
)

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"` // user, assistant, system
	Parts     []Part    `json:"parts"`
	CreatedAt int64     `json:"createdAt"`
}

// Part represents a part of a message
type Part struct {
	Type string `json:"type"` // text, image, tool_use, tool_result
	
	// For text parts
	Text string `json:"text,omitempty"`
	
	// For image parts
	ImageURL string `json:"imageUrl,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	
	// For tool use
	ToolID   string      `json:"toolId,omitempty"`
	ToolName string      `json:"toolName,omitempty"`
	ToolArgs interface{} `json:"toolArgs,omitempty"`
	
	// For tool result
	ToolResult interface{} `json:"toolResult,omitempty"`
	IsError    bool        `json:"isError,omitempty"`
	
	// For thinking
	Thinking string `json:"thinking,omitempty"`
}

// NewMessage creates a new message
func NewMessage(role string, parts []Part) *Message {
	return &Message{
		ID:        id.New(id.PrefixMessage).String(),
		Role:      role,
		Parts:     parts,
		CreatedAt: time.Now().UnixMilli(),
	}
}

// TextPart creates a text part
func TextPart(text string) Part {
	return Part{Type: "text", Text: text}
}

// ImagePart creates an image part
func ImagePart(url, mimeType string) Part {
	return Part{Type: "image", ImageURL: url, MimeType: mimeType}
}

// ToolUsePart creates a tool use part
func ToolUsePart(toolID, toolName string, args interface{}) Part {
	return Part{
		Type:     "tool_use",
		ToolID:   toolID,
		ToolName: toolName,
		ToolArgs: args,
	}
}

// ToolResultPart creates a tool result part
func ToolResultPart(toolID string, result interface{}, isError bool) Part {
	return Part{
		Type:       "tool_result",
		ToolID:     toolID,
		ToolResult: result,
		IsError:    isError,
	}
}

// ThinkingPart creates a thinking part
func ThinkingPart(thinking string) Part {
	return Part{Type: "thinking", Thinking: thinking}
}

// Conversation represents a conversation with messages
type Conversation struct {
	SessionID string     `json:"sessionId"`
	Messages  []*Message `json:"messages"`
}

// AddMessage adds a message to the conversation
func (c *Conversation) AddMessage(msg *Message) {
	c.Messages = append(c.Messages, msg)
}

// LastMessage returns the last message
func (c *Conversation) LastMessage() *Message {
	if len(c.Messages) == 0 {
		return nil
	}
	return c.Messages[len(c.Messages)-1]
}

// GetMessagesForAPI returns messages formatted for API calls
func (c *Conversation) GetMessagesForAPI() []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(c.Messages))
	
	for _, msg := range c.Messages {
		apiMsg := map[string]interface{}{
			"role": msg.Role,
		}
		
		// Convert parts to content
		var content []map[string]interface{}
		for _, part := range msg.Parts {
			switch part.Type {
			case "text":
				content = append(content, map[string]interface{}{
					"type": "text",
					"text": part.Text,
				})
			case "image":
				content = append(content, map[string]interface{}{
					"type": "image",
					"image": map[string]interface{}{
						"url": part.ImageURL,
					},
				})
			case "tool_use":
				content = append(content, map[string]interface{}{
					"type":      "tool_use",
					"id":        part.ToolID,
					"name":      part.ToolName,
					"arguments": part.ToolArgs,
				})
			case "tool_result":
				content = append(content, map[string]interface{}{
					"type":       "tool_result",
					"tool_use_id": part.ToolID,
					"content":    part.ToolResult,
					"is_error":   part.IsError,
				})
			}
		}
		
		if len(content) == 1 && content[0]["type"] == "text" {
			apiMsg["content"] = content[0]["text"]
		} else {
			apiMsg["content"] = content
		}
		
		result = append(result, apiMsg)
	}
	
	return result
}
