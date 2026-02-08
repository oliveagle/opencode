package provider

import (
	"fmt"
	"os"
	"sync"
)

// GlobalRegistry is the global provider registry
var GlobalRegistry = NewRegistry()

func init() {
	// Register providers with API keys from environment
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		GlobalRegistry.Register(NewAnthropicProvider(apiKey))
	}
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		GlobalRegistry.Register(NewOpenAIProvider(apiKey))
	}
	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		GlobalRegistry.Register(NewGoogleProvider(apiKey))
	} else if apiKey := os.Getenv("GOOGLE_API_KEY"); apiKey != "" {
		GlobalRegistry.Register(NewGoogleProvider(apiKey))
	}
	if apiKey := os.Getenv("OPENROUTER_API_KEY"); apiKey != "" {
		GlobalRegistry.Register(NewOpenRouterProvider(apiKey))
	}
}

// SetupProviders initializes providers from configuration
func SetupProviders(configs map[string]ProviderConfig) error {
	for providerID, cfg := range configs {
		if cfg.Disabled {
			continue
		}

		var p Provider
		switch providerID {
		case "anthropic":
			p = NewAnthropicProvider(cfg.APIKey)
		case "openai":
			p = NewOpenAIProvider(cfg.APIKey)
		case "google", "gemini":
			p = NewGoogleProvider(cfg.APIKey)
		case "openrouter":
			p = NewOpenRouterProvider(cfg.APIKey)
		default:
			continue
		}

		GlobalRegistry.Register(p)
	}
	return nil
}

// ProviderConfig represents provider configuration
type ProviderConfig struct {
	APIKey   string            `json:"apiKey"`
	BaseURL  string            `json:"baseUrl,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Disabled bool              `json:"disabled,omitempty"`
}

// GetProvider returns a provider by name
func GetProvider(name string) (Provider, error) {
	p, ok := GlobalRegistry.Get(name)
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return p, nil
}

// GetModel returns a model by ID
func GetModel(id string) (Model, error) {
	m, ok := GlobalRegistry.GetModel(id)
	if !ok {
		return Model{}, fmt.Errorf("model %s not found", id)
	}
	return m, nil
}

// ListAllModels returns all available models
func ListAllModels() []Model {
	return GlobalRegistry.ListModels()
}

// Chat sends a chat request using the specified provider and model
func Chat(providerID, modelID string, req *ChatRequest) (*ChatResponse, error) {
	p, err := GetProvider(providerID)
	if err != nil {
		return nil, err
	}
	req.Model = modelID
	return p.Chat(nil, req) // TODO: context propagation
}

// Stream sends a streaming chat request
func Stream(providerID, modelID string, req *ChatRequest, handler StreamHandler) error {
	p, err := GetProvider(providerID)
	if err != nil {
		return err
	}
	req.Model = modelID
	return p.Stream(nil, req, handler) // TODO: context propagation
}

// ProviderInfo returns information about a provider
type ProviderInfo struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Models []string `json:"models"`
}

// ListProviders returns information about all providers
func ListProviders() []ProviderInfo {
	infos := make([]ProviderInfo, 0)
	for name, p := range GlobalRegistry.providers {
		models := p.Models()
		modelIDs := make([]string, len(models))
		for i, m := range models {
			modelIDs[i] = m.ID
		}
		infos = append(infos, ProviderInfo{
			ID:     name,
			Name:   p.Name(),
			Models: modelIDs,
		})
	}
	return infos
}

// RegistryWithLock provides thread-safe access to the registry
type RegistryWithLock struct {
	registry *Registry
	mu       sync.RWMutex
}

// NewRegistryWithLock creates a thread-safe registry wrapper
func NewRegistryWithLock() *RegistryWithLock {
	return &RegistryWithLock{
		registry: NewRegistry(),
	}
}

// Register registers a provider
func (r *RegistryWithLock) Register(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registry.Register(p)
}

// Get retrieves a provider by name
func (r *RegistryWithLock) Get(name string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.registry.Get(name)
}

// GetModel retrieves a model by ID
func (r *RegistryWithLock) GetModel(id string) (Model, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.registry.GetModel(id)
}

// ListModels returns all available models
func (r *RegistryWithLock) ListModels() []Model {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.registry.ListModels()
}
