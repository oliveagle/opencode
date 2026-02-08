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

// GoogleProvider implements the Provider interface for Google AI
type GoogleProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
	models  []Model
}

// NewGoogleProvider creates a new Google provider
func NewGoogleProvider(apiKey string) *GoogleProvider {
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			apiKey = os.Getenv("GOOGLE_API_KEY")
		}
	}
	
	return &GoogleProvider{
		apiKey:  apiKey,
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
		client: &http.Client{Timeout: 5 * time.Minute},
		models: []Model{
			{ID: "gemini-2.5-pro-latest", Name: "Gemini 2.5 Pro", Provider: "google", ContextSize: 1048576, MaxOutput: 65536, InputPrice: 1.25, OutputPrice: 10.0, Modalities: []string{"text", "image", "audio", "video"}},
			{ID: "gemini-2.0-flash", Name: "Gemini 2.0 Flash", Provider: "google", ContextSize: 1048576, MaxOutput: 8192, InputPrice: 0.1, OutputPrice: 0.4, Modalities: []string{"text", "image", "audio"}},
			{ID: "gemini-2.0-flash-lite", Name: "Gemini 2.0 Flash Lite", Provider: "google", ContextSize: 1048576, MaxOutput: 8192, InputPrice: 0.075, OutputPrice: 0.3, Modalities: []string{"text", "image"}},
		},
	}
}

// Name returns the provider name
func (p *GoogleProvider) Name() string {
	return "google"
}

// Models returns available models
func (p *GoogleProvider) Models() []Model {
	return p.models
}

// Chat sends a chat request
func (p *GoogleProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// Convert to Google format
	googleReq := map[string]interface{}{
		"contents": convertMessagesForGoogle(req.Messages),
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": req.MaxTokens,
		},
	}

	if len(req.Tools) > 0 {
		googleReq["tools"] = convertToolsForGoogle(req.Tools)
	}

	body, err := json.Marshal(googleReq)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.baseURL, req.Model, p.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google API error: %s - %s", resp.Status, string(respBody))
	}

	var googleResp googleResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, err
	}

	return convertGoogleResponse(&googleResp, req.Model), nil
}

// Stream sends a chat request with streaming
func (p *GoogleProvider) Stream(ctx context.Context, req *ChatRequest, handler StreamHandler) error {
	googleReq := map[string]interface{}{
		"contents": convertMessagesForGoogle(req.Messages),
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": req.MaxTokens,
		},
	}

	if len(req.Tools) > 0 {
		googleReq["tools"] = convertToolsForGoogle(req.Tools)
	}

	body, err := json.Marshal(googleReq)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/models/%s:streamGenerateContent?key=%s&alt=sse", p.baseURL, req.Model, p.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse SSE stream
	decoder := json.NewDecoder(resp.Body)
	for {
		var event googleResponse
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		if len(event.Candidates) > 0 {
			candidate := event.Candidates[0]
			if len(candidate.Content.Parts) > 0 {
				for _, part := range candidate.Content.Parts {
					if part.Text != "" {
						handler(StreamEvent{
							Type:    "content",
							Content: part.Text,
						})
					}
					if part.FunctionCall != nil {
						handler(StreamEvent{
							Type: "tool_use",
							ToolCall: &ToolCall{
								ID:        part.FunctionCall.Name,
								Name:      part.FunctionCall.Name,
								Arguments: part.FunctionCall.Args,
							},
						})
					}
				}
			}
			if candidate.FinishReason != "" {
				handler(StreamEvent{Type: "done"})
			}
		}
	}

	return nil
}

type googleResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text,omitempty"`
				FunctionCall *struct {
					Name string                 `json:"name"`
					Args map[string]interface{} `json:"args"`
				} `json:"functionCall,omitempty"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata *struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

func convertMessagesForGoogle(messages []Message) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(messages))
	for _, msg := range messages {
		parts := make([]map[string]interface{}, 0)
		for _, part := range msg.Content {
			switch part.Type {
			case "text":
				parts = append(parts, map[string]interface{}{
					"text": part.Text,
				})
			case "image":
				parts = append(parts, map[string]interface{}{
					"inlineData": map[string]interface{}{
						"mimeType": part.MimeType,
						"data":     part.ImageURL,
					},
				})
			case "tool_result":
				parts = append(parts, map[string]interface{}{
					"functionResponse": map[string]interface{}{
						"name":    part.ToolID,
						"content": part.ToolResult,
					},
				})
			}
		}

		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}

		result = append(result, map[string]interface{}{
			"role":  role,
			"parts": parts,
		})
	}
	return result
}

func convertToolsForGoogle(tools []Tool) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(tools))
	for _, tool := range tools {
		result = append(result, map[string]interface{}{
			"functionDeclarations": []map[string]interface{}{
				{
					"name":        tool.Name,
					"description": tool.Description,
					"parameters":  tool.Parameters,
				},
			},
		})
	}
	return result
}

func convertGoogleResponse(resp *googleResponse, model string) *ChatResponse {
	if len(resp.Candidates) == 0 {
		return &ChatResponse{Model: model}
	}

	candidate := resp.Candidates[0]
	msg := Message{
		Role: "assistant",
	}

	var toolCalls []ToolCall
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			msg.Content = append(msg.Content, Content{
				Type: "text",
				Text: part.Text,
			})
		}
		if part.FunctionCall != nil {
			toolCalls = append(toolCalls, ToolCall{
				ID:        part.FunctionCall.Name,
				Name:      part.FunctionCall.Name,
				Arguments: part.FunctionCall.Args,
			})
		}
	}

	usage := Usage{}
	if resp.UsageMetadata != nil {
		usage.InputTokens = resp.UsageMetadata.PromptTokenCount
		usage.OutputTokens = resp.UsageMetadata.CandidatesTokenCount
	}

	return &ChatResponse{
		Model:        model,
		Created:      time.Now(),
		Message:      msg,
		Usage:        usage,
		ToolCalls:    toolCalls,
		FinishReason: candidate.FinishReason,
	}
}
