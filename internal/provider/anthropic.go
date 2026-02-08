package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// AnthropicProvider implements the Provider interface for Anthropic
type AnthropicProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
	models  []Model
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	
	return &AnthropicProvider{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1",
		client: &http.Client{Timeout: 5 * time.Minute},
		models: []Model{
			{ID: "claude-sonnet-4-20250514", Name: "Claude Sonnet 4", Provider: "anthropic", ContextSize: 200000, MaxOutput: 16000, InputPrice: 3.0, OutputPrice: 15.0, Modalities: []string{"text", "image"}},
			{ID: "claude-opus-4-20250514", Name: "Claude Opus 4", Provider: "anthropic", ContextSize: 200000, MaxOutput: 16000, InputPrice: 15.0, OutputPrice: 75.0, Modalities: []string{"text", "image"}},
			{ID: "claude-3-5-sonnet-20241022", Name: "Claude 3.5 Sonnet", Provider: "anthropic", ContextSize: 200000, MaxOutput: 8192, InputPrice: 3.0, OutputPrice: 15.0, Modalities: []string{"text", "image"}},
			{ID: "claude-3-5-haiku-20241022", Name: "Claude 3.5 Haiku", Provider: "anthropic", ContextSize: 200000, MaxOutput: 8192, InputPrice: 0.8, OutputPrice: 4.0, Modalities: []string{"text", "image"}},
		},
	}
}

// Name returns the provider name
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// Models returns available models
func (p *AnthropicProvider) Models() []Model {
	return p.models
}

// Chat sends a chat request
func (p *AnthropicProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// Convert to Anthropic format
	anthropicReq := map[string]interface{}{
		"model":      req.Model,
		"max_tokens": req.MaxTokens,
		"messages":   convertMessagesForAnthropic(req.Messages),
	}

	if req.System != "" {
		anthropicReq["system"] = req.System
	}
	if len(req.Tools) > 0 {
		anthropicReq["tools"] = convertToolsForAnthropic(req.Tools)
	}

	body, err := json.Marshal(anthropicReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("anthropic-beta", "claude-code-20250219,interleaved-thinking-2025-05-14")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anthropic API error: %s - %s", resp.Status, string(respBody))
	}

	var anthropicResp anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return nil, err
	}

	return convertAnthropicResponse(&anthropicResp), nil
}

// Stream sends a chat request with streaming
func (p *AnthropicProvider) Stream(ctx context.Context, req *ChatRequest, handler StreamHandler) error {
	anthropicReq := map[string]interface{}{
		"model":      req.Model,
		"max_tokens": req.MaxTokens,
		"messages":   convertMessagesForAnthropic(req.Messages),
		"stream":     true,
	}

	if len(req.Tools) > 0 {
		anthropicReq["tools"] = convertToolsForAnthropic(req.Tools)
	}

	body, err := json.Marshal(anthropicReq)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse SSE stream
	decoder := json.NewDecoder(resp.Body)
	for {
		var event struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data,omitempty"`
			
			// For content_block_delta
			Delta *struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"delta,omitempty"`
			
			// For message_start
			Message *anthropicResponse `json:"message,omitempty"`
			
			// For content_block_start
			Index         int `json:"index"`
			ContentBlock  *struct {
				Type string `json:"type"`
				Text string `json:"text,omitempty"`
				ID   string `json:"id,omitempty"`
				Name string `json:"name,omitempty"`
				Input json.RawMessage `json:"input,omitempty"`
			} `json:"content_block,omitempty"`
		}

		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch event.Type {
		case "content_block_delta":
			if event.Delta != nil && event.Delta.Type == "text_delta" {
				handler(StreamEvent{
					Type:    "content",
					Content: event.Delta.Text,
				})
			}
		case "content_block_start":
			if event.ContentBlock != nil && event.ContentBlock.Type == "tool_use" {
				var args interface{}
				if event.ContentBlock.Input != nil {
					json.Unmarshal(event.ContentBlock.Input, &args)
				}
				handler(StreamEvent{
					Type: "tool_use",
					ToolCall: &ToolCall{
						ID:        event.ContentBlock.ID,
						Name:      event.ContentBlock.Name,
						Arguments: args,
					},
				})
			}
		case "message_stop":
			handler(StreamEvent{Type: "done"})
		case "message_start":
			if event.Message != nil {
				handler(StreamEvent{
					Type: "usage",
					Usage: &Usage{
						InputTokens:  event.Message.Usage.InputTokens,
						OutputTokens: event.Message.Usage.OutputTokens,
					},
				})
			}
		}
	}

	return nil
}

type anthropicResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
		
		// For tool use
		ID    string          `json:"id,omitempty"`
		Name  string          `json:"name,omitempty"`
		Input json.RawMessage `json:"input,omitempty"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	StopReason string `json:"stop_reason"`
}

func convertMessagesForAnthropic(messages []Message) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(messages))
	for _, msg := range messages {
		content := make([]map[string]interface{}, 0)
		for _, part := range msg.Content {
			switch part.Type {
			case "text":
				content = append(content, map[string]interface{}{
					"type": "text",
					"text": part.Text,
				})
			case "image":
				content = append(content, map[string]interface{}{
					"type": "image",
					"source": map[string]interface{}{
						"type":      "url",
						"url":       part.ImageURL,
						"media_type": part.MimeType,
					},
				})
			case "tool_use":
				content = append(content, map[string]interface{}{
					"type":      "tool_use",
					"id":        part.ToolID,
					"name":      part.ToolName,
					"input":     part.ToolArgs,
				})
			case "tool_result":
				content = append(content, map[string]interface{}{
					"type":        "tool_result",
					"tool_use_id": part.ToolID,
					"content":     part.ToolResult,
					"is_error":    part.IsError,
				})
			}
		}
		result = append(result, map[string]interface{}{
			"role":    msg.Role,
			"content": content,
		})
	}
	return result
}

func convertToolsForAnthropic(tools []Tool) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(tools))
	for _, tool := range tools {
		result = append(result, map[string]interface{}{
			"name":        tool.Name,
			"description": tool.Description,
			"input_schema": tool.Parameters,
		})
	}
	return result
}

func convertAnthropicResponse(resp *anthropicResponse) *ChatResponse {
	msg := Message{
		Role: resp.Role,
	}
	
	var toolCalls []ToolCall
	
	for _, block := range resp.Content {
		if block.Type == "text" {
			msg.Content = append(msg.Content, Content{
				Type: "text",
				Text: block.Text,
			})
		} else if block.Type == "tool_use" {
			var args interface{}
			if block.Input != nil {
				json.Unmarshal(block.Input, &args)
			}
			toolCalls = append(toolCalls, ToolCall{
				ID:        block.ID,
				Name:      block.Name,
				Arguments: args,
			})
		}
	}

	return &ChatResponse{
		ID:           resp.ID,
		Model:        resp.Model,
		Created:      time.Now(),
		Message:      msg,
		Usage:        Usage{InputTokens: resp.Usage.InputTokens, OutputTokens: resp.Usage.OutputTokens},
		ToolCalls:    toolCalls,
		FinishReason: resp.StopReason,
	}
}
