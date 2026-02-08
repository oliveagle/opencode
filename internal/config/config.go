package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/anomalyco/opencode/internal/installation"
)

// Info represents the configuration
type Info struct {
	// Schema reference
	Schema string `json:"$schema,omitempty"`

	// Model configuration
	Model *ModelConfig `json:"model,omitempty"`

	// Provider configurations
	Providers map[string]ProviderConfig `json:"providers,omitempty"`

	// Agent configurations
	Agents map[string]AgentConfig `json:"agents,omitempty"`

	// MCP server configurations
	MCP map[string]MCPServerConfig `json:"mcp,omitempty"`

	// Permission rules
	Permission map[string]interface{} `json:"permission,omitempty"`

	// Custom instructions
	Instructions []string `json:"instructions,omitempty"`

	// Plugin configurations
	Plugins []string `json:"plugin,omitempty"`

	// Theme settings
	Theme *ThemeConfig `json:"theme,omitempty"`

	// Keybindings
	Keybindings map[string]string `json:"keybinds,omitempty"`
}

// ModelConfig represents model selection
type ModelConfig struct {
	ModelID    string `json:"modelId,omitempty"`
	ProviderID string `json:"providerId,omitempty"`
}

// ProviderConfig represents provider configuration
type ProviderConfig struct {
	APIKey  string            `json:"apiKey,omitempty"`
	BaseURL string            `json:"baseUrl,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Disabled bool             `json:"disabled,omitempty"`
}

// AgentConfig represents agent configuration
type AgentConfig struct {
	Model       *ModelConfig         `json:"model,omitempty"`
	Permission  map[string]interface{} `json:"permission,omitempty"`
	Prompt      string               `json:"prompt,omitempty"`
	Temperature float64              `json:"temperature,omitempty"`
	TopP        float64              `json:"topP,omitempty"`
}

// MCPServerConfig represents MCP server configuration
type MCPServerConfig struct {
	Type    string            `json:"type"` // stdio, sse, ws
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	URL     string            `json:"url,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Disabled bool             `json:"disabled,omitempty"`
}

// ThemeConfig represents theme configuration
type ThemeConfig struct {
	Mode  string `json:"mode,omitempty"` // light, dark, system
	Color string `json:"color,omitempty"`
}

// Manager handles configuration loading and saving
type Manager struct {
	mu     sync.RWMutex
	config *Info
	path   string
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		config: &Info{
			Providers:   make(map[string]ProviderConfig),
			Agents:      make(map[string]AgentConfig),
			MCP:         make(map[string]MCPServerConfig),
			Plugins:     []string{},
			Instructions: []string{},
			Keybindings: make(map[string]string),
			Permission:  make(map[string]interface{}),
		},
	}
}

// Load loads configuration from all sources
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load in order of precedence (low to high)
	// 1. Global config
	if err := m.loadGlobal(); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 2. Project config
	if err := m.loadProject(); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 3. Environment config
	m.loadFromEnv()

	return nil
}

func (m *Manager) loadGlobal() error {
	configDir := installation.ConfigDir()
	configPath := filepath.Join(configDir, "opencode.json")
	
	// Try .jsonc if .json doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(configDir, "opencode.jsonc")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return os.ErrNotExist
		}
	}

	return m.loadFile(configPath)
}

func (m *Manager) loadProject() error {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, "opencode.json")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(cwd, "opencode.jsonc")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return os.ErrNotExist
		}
	}

	return m.loadFile(configPath)
}

func (m *Manager) loadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse JSON (TODO: handle JSONC with comments)
	var cfg Info
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	// Merge configurations
	m.merge(&cfg)
	m.path = path

	return nil
}

func (m *Manager) loadFromEnv() {
	// Load API keys from environment
	envKeys := map[string]string{
		"ANTHROPIC_API_KEY": "anthropic",
		"OPENAI_API_KEY":    "openai",
		"GEMINI_API_KEY":    "google",
		"AZURE_API_KEY":     "azure",
		"GROQ_API_KEY":      "groq",
		"MISTRAL_API_KEY":   "mistral",
		"COHERE_API_KEY":    "cohere",
		"XAI_API_KEY":       "xai",
		"OPENROUTER_API_KEY": "openrouter",
	}

	for envKey, providerID := range envKeys {
		if apiKey := os.Getenv(envKey); apiKey != "" {
			if _, exists := m.config.Providers[providerID]; !exists {
				m.config.Providers[providerID] = ProviderConfig{APIKey: apiKey}
			}
		}
	}
}

func (m *Manager) merge(other *Info) {
	if other.Schema != "" {
		m.config.Schema = other.Schema
	}
	if other.Model != nil {
		m.config.Model = other.Model
	}
	for k, v := range other.Providers {
		m.config.Providers[k] = v
	}
	for k, v := range other.Agents {
		m.config.Agents[k] = v
	}
	for k, v := range other.MCP {
		m.config.MCP[k] = v
	}
	if other.Permission != nil {
		m.config.Permission = other.Permission
	}
	if len(other.Instructions) > 0 {
		m.config.Instructions = append(m.config.Instructions, other.Instructions...)
	}
	if len(other.Plugins) > 0 {
		m.config.Plugins = append(m.config.Plugins, other.Plugins...)
	}
	if other.Theme != nil {
		m.config.Theme = other.Theme
	}
	for k, v := range other.Keybindings {
		m.config.Keybindings[k] = v
	}
}

// Get returns the current configuration
func (m *Manager) Get() *Info {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// Save saves the configuration to file
func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.path == "" {
		// Default to global config
		configDir := installation.ConfigDir()
		os.MkdirAll(configDir, 0755)
		m.path = filepath.Join(configDir, "opencode.json")
	}

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0644)
}

// Set updates a configuration value
func (m *Manager) Set(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO: Implement nested key setting
	return nil
}

// Global returns the global configuration manager
var Global = NewManager()

// Get is a convenience function
func Get() *Info {
	return Global.Get()
}
