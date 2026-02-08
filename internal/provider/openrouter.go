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

// OpenRouterProvider implements the Provider interface for OpenRouter
type OpenRouterProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
	models  []Model
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(apiKey string) *OpenRouterProvider {
	if apiKey == "" {
		apiKey = os.Getenv("OPENROUTER_API_KEY")
	}
	
	return &OpenRouterProvider{
		apiKey:  apiKey,
		baseURL: "https://openrouter.ai/api/v1",
		client: &http.Client{Timeout: 5 * time.Minute},
		models: []Model{
			// OpenRouter provides access to many models through a unified API
			// These are some popular ones
			{ID: "anthropic/claude-sonnet-4", Name: "Claude Sonnet 4 (OpenRouter)", Provider: "openrouter", ContextSize: 200000, MaxOutput: 16000, InputPrice: 3.0, OutputPrice: 15.0, Modalities: []string{"text", "image"}},
			{ID: "anthropic/claude-opus-4", Name: "Claude Opus 4 (OpenRouter)", Provider: "openrouter", ContextSize: 200000, MaxOutput: 16000, InputPrice: 15.0, OutputPrice: 75.0, Modalities: []string{"text", "image"}},
			{ID: "openai/gpt-4o", Name: "GPT-4o (OpenRouter)", Provider: "openrouter", ContextSize: 128000, MaxOutput: 16384, InputPrice: 2.5, OutputPrice: 10.0, Modalities: []string{"text", "image"}},
			{ID: "google/gemini-2.5-pro", Name: "Gemini 2.5 Pro (OpenRouter)", Provider: "openrouter", ContextSize: 1048576, MaxOutput: 65536, InputPrice: 1.25, OutputPrice: 10.0, Modalities: []string{"text", "image"}},
			{ID: "meta-llama/llama-3.3-70b-instruct", Name: "Llama 3.3 70B (OpenRouter)", Provider: "openrouter", ContextSize: 131072, MaxOutput: 8192, InputPrice: 0.35, OutputPrice: 0.4, Modalities: []string{"text"}},
			{ID: "deepseek/deepseek-r1", Name: "DeepSeek R1 (OpenRouter)", Provider: "openrouter", ContextSize: 163840, MaxOutput: 16384, InputPrice: 0.55, OutputPrice: 2.19, Modalities: []string{"text"}},
		},
	}
}

// Name returns the provider name
func (p *OpenRouterProvider) Name() string {
	return "openrouter"
}

// Models returns available models
func (p *OpenRouterProvider) Models() []Model {
	return p.models
}

// Chat sends a chat request
func (p *OpenRouterProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	orReq := map[string]interface{}{
		"model":    req.Model,
		"messages": convertMessagesForOpenAI(req.Messages),
	}

	if req.MaxTokens > 0 {
		orReq["max_tokens"] = req.MaxTokens
	}
	if len(req.Tools) > 0 {
		orReq["tools"] = convertToolsForOpenAI(req.Tools)
	}

	body, err := json.Marshal(orReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://opencode.ai")
	httpReq.Header.Set("X-Title", "OpenCode")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openrouter API error: %s - %s", resp.Status, string(respBody))
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, err
	}

	return convertOpenAIResponse(&openAIResp), nil
}

// Stream sends a chat request with streaming
func (p *OpenRouterProvider) Stream(ctx context.Context, req *ChatRequest, handler StreamHandler) error {
	orReq := map[string]interface{}{
		"model":    req.Model,
		"messages": convertMessagesForOpenAI(req.Messages),
		"stream":   true,
	}

	if req.MaxTokens > 0 {
		orReq["max_tokens"] = req.MaxTokens
	}
	if len(req.Tools) > 0 {
		orReq["tools"] = convertToolsForOpenAI(req.Tools)
	}

	body, err := json.Marshal(orReq)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://opencode.ai")
	httpReq.Header.Set("X-Title", "OpenCode")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse SSE stream - same format as OpenAI
	reader := json.NewDecoder(resp.Body)
	for {
		var event struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content,omitempty"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}

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
			if choice.FinishReason != "" {
				handler(StreamEvent{Type: "done"})
			}
		}
	}

	return nil
}

// FetchModels fetches available models from OpenRouter
func (p *OpenRouterProvider) FetchModels(ctx context.Context) ([]Model, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/models", nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var modelsResp struct {
		Data []struct {
			ID          string  `json:"id"`
			Name        string  `json:"name"`
			ContextSize int     `json:"context_length"`
			Pricing     struct {
				Prompt  float64 `json:"prompt,string"`
				Completion float64 `json:"completion,string"`
			} `json:"pricing"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, err
	}

	models := make([]Model, 0, len(modelsResp.Data))
	for _, m := range modelsResp.Data {
		models = append(models, Model{
			ID:          m.ID,
			Name:        m.Name,
			Provider:    "openrouter",
			ContextSize: m.ContextSize,
			InputPrice:  m.Pricing.Prompt * 1000000,  // Convert to per 1M tokens
			OutputPrice: m.Pricing.Completion * 1000000,
		})
	}

	p.models = models
	return models, nil
}
