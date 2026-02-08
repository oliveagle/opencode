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

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
	models  []Model
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1",
		client: &http.Client{Timeout: 5 * time.Minute},
		models: []Model{
			{ID: "gpt-4o", Name: "GPT-4o", Provider: "openai", ContextSize: 128000, MaxOutput: 16384, InputPrice: 2.5, OutputPrice: 10.0, Modalities: []string{"text", "image", "audio"}},
			{ID: "gpt-4o-mini", Name: "GPT-4o Mini", Provider: "openai", ContextSize: 128000, MaxOutput: 16384, InputPrice: 0.15, OutputPrice: 0.6, Modalities: []string{"text", "image", "audio"}},
			{ID: "gpt-4-turbo", Name: "GPT-4 Turbo", Provider: "openai", ContextSize: 128000, MaxOutput: 4096, InputPrice: 10.0, OutputPrice: 30.0, Modalities: []string{"text", "image"}},
			{ID: "o1", Name: "o1", Provider: "openai", ContextSize: 200000, MaxOutput: 100000, InputPrice: 15.0, OutputPrice: 60.0, Modalities: []string{"text"}},
			{ID: "o1-mini", Name: "o1 Mini", Provider: "openai", ContextSize: 128000, MaxOutput: 65536, InputPrice: 1.5, OutputPrice: 6.0, Modalities: []string{"text"}},
			{ID: "o3-mini", Name: "o3 Mini", Provider: "openai", ContextSize: 200000, MaxOutput: 100000, InputPrice: 1.1, OutputPrice: 4.4, Modalities: []string{"text"}},
		},
	}
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Models returns available models
func (p *OpenAIProvider) Models() []Model {
	return p.models
}

// Chat sends a chat request
func (p *OpenAIProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	openAIReq := map[string]interface{}{
		"model":    req.Model,
		"messages": convertMessagesForOpenAI(req.Messages),
	}

	if req.MaxTokens > 0 {
		openAIReq["max_tokens"] = req.MaxTokens
	}
	if req.Temperature > 0 {
		openAIReq["temperature"] = req.Temperature
	}
	if req.TopP > 0 {
		openAIReq["top_p"] = req.TopP
	}
	if len(req.Tools) > 0 {
		openAIReq["tools"] = convertToolsForOpenAI(req.Tools)
	}

	body, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai API error: %s - %s", resp.Status, string(respBody))
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, err
	}

	return convertOpenAIResponse(&openAIResp), nil
}

// Stream sends a chat request with streaming
func (p *OpenAIProvider) Stream(ctx context.Context, req *ChatRequest, handler StreamHandler) error {
	openAIReq := map[string]interface{}{
		"model":    req.Model,
		"messages": convertMessagesForOpenAI(req.Messages),
		"stream":   true,
	}

	if req.MaxTokens > 0 {
		openAIReq["max_tokens"] = req.MaxTokens
	}
	if len(req.Tools) > 0 {
		openAIReq["tools"] = convertToolsForOpenAI(req.Tools)
	}

	body, err := json.Marshal(openAIReq)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse SSE stream
	reader := json.NewDecoder(resp.Body)
	for {
		var event struct {
			ID      string `json:"id"`
			Object  string `json:"object"`
			Choices []struct {
				Index int `json:"index"`
				Delta struct {
					Role      string `json:"role,omitempty"`
					Content   string `json:"content,omitempty"`
					ToolCalls []struct {
						ID       string          `json:"id,omitempty"`
						Type     string          `json:"type,omitempty"`
						Function struct {
							Name      string          `json:"name,omitempty"`
							Arguments json.RawMessage `json:"arguments,omitempty"`
						} `json:"function"`
					} `json:"tool_calls,omitempty"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
			Usage *struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			} `json:"usage,omitempty"`
		}

		// Read SSE data: prefix
		var line string
		if err := reader.Decode(&line); err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		if line == "" || line == "[DONE]" {
			continue
		}

		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		if len(event.Choices) > 0 {
			choice := event.Choices[0]
			if choice.Delta.Content != "" {
				handler(StreamEvent{
					Type:    "content",
					Content: choice.Delta.Content,
				})
			}
			if len(choice.Delta.ToolCalls) > 0 {
				for _, tc := range choice.Delta.ToolCalls {
					var args interface{}
					if tc.Function.Arguments != nil {
						json.Unmarshal(tc.Function.Arguments, &args)
					}
					handler(StreamEvent{
						Type: "tool_use",
						ToolCall: &ToolCall{
							ID:        tc.ID,
							Name:      tc.Function.Name,
							Arguments: args,
						},
					})
				}
			}
			if choice.FinishReason != "" {
				handler(StreamEvent{Type: "done"})
			}
		}

		if event.Usage != nil {
			handler(StreamEvent{
				Type: "usage",
				Usage: &Usage{
					InputTokens:  event.Usage.PromptTokens,
					OutputTokens: event.Usage.CompletionTokens,
				},
			})
		}
	}

	return nil
}

type openAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Message struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string      `json:"name"`
					Arguments interface{} `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func convertMessagesForOpenAI(messages []Message) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(messages))
	for _, msg := range messages {
		content := make([]interface{}, 0)
		for _, part := range msg.Content {
			switch part.Type {
			case "text":
				content = append(content, map[string]interface{}{
					"type": "text",
					"text": part.Text,
				})
			case "image":
				content = append(content, map[string]interface{}{
					"type": "image_url",
					"image_url": map[string]interface{}{
						"url": part.ImageURL,
					},
				})
			}
		}

		msgMap := map[string]interface{}{
			"role": msg.Role,
		}

		if len(content) == 1 {
			if textContent, ok := content[0].(map[string]interface{}); ok && textContent["type"] == "text" {
				msgMap["content"] = textContent["text"]
			} else {
				msgMap["content"] = content
			}
		} else {
			msgMap["content"] = content
		}

		result = append(result, msgMap)
	}
	return result
}

func convertToolsForOpenAI(tools []Tool) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(tools))
	for _, tool := range tools {
		result = append(result, map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        tool.Name,
				"description": tool.Description,
				"parameters":  tool.Parameters,
			},
		})
	}
	return result
}

func convertOpenAIResponse(resp *openAIResponse) *ChatResponse {
	if len(resp.Choices) == 0 {
		return &ChatResponse{}
	}

	choice := resp.Choices[0]
	msg := Message{
		Role: choice.Message.Role,
	}

	var toolCalls []ToolCall
	for _, tc := range choice.Message.ToolCalls {
		toolCalls = append(toolCalls, ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: tc.Function.Arguments,
		})
	}

	msg.Content = append(msg.Content, Content{
		Type: "text",
		Text: choice.Message.Content,
	})

	return &ChatResponse{
		ID:           resp.ID,
		Model:        resp.Model,
		Created:      time.Unix(resp.Created, 0),
		Message:      msg,
		Usage:        Usage{InputTokens: resp.Usage.PromptTokens, OutputTokens: resp.Usage.CompletionTokens},
		ToolCalls:    toolCalls,
		FinishReason: choice.FinishReason,
	}
}
