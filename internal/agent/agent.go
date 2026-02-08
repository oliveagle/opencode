package agent

import (
	"github.com/anomalyco/opencode/internal/permission"
)

// Mode represents the agent mode
type Mode string

const (
	ModeSubagent Mode = "subagent"
	ModePrimary  Mode = "primary"
	ModeAll      Mode = "all"
)

// Info represents agent configuration
type Info struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Mode        Mode              `json:"mode"`
	Native      bool              `json:"native,omitempty"`
	Hidden      bool              `json:"hidden,omitempty"`
	TopP        float64           `json:"topP,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	Color       string            `json:"color,omitempty"`
	Permission  *permission.Ruleset `json:"permission"`
	Model       *ModelRef         `json:"model,omitempty"`
	Variant     string            `json:"variant,omitempty"`
	Prompt      string            `json:"prompt,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Steps       int               `json:"steps,omitempty"`
}

// ModelRef references a specific model
type ModelRef struct {
	ModelID    string `json:"modelId"`
	ProviderID string `json:"providerId"`
}

// Registry holds all registered agents
type Registry struct {
	agents map[string]*Info
}

// NewRegistry creates a new agent registry
func NewRegistry() *Registry {
	r := &Registry{
		agents: make(map[string]*Info),
	}
	
	// Register default agents
	r.registerDefaults()
	
	return r
}

func (r *Registry) registerDefaults() {
	// Build agent - default, full-access
	r.Register(&Info{
		Name:        "build",
		Description: "The default agent. Executes tools based on configured permissions.",
		Mode:        ModePrimary,
		Native:      true,
		Permission:  permission.DefaultRuleset(),
		Options:     make(map[string]interface{}),
	})
	
	// Plan agent - read-only
	r.Register(&Info{
		Name:        "plan",
		Description: "Plan mode. Disallows all edit tools.",
		Mode:        ModePrimary,
		Native:      true,
		Permission:  permission.ReadOnlyRuleset(),
		Options:     make(map[string]interface{}),
	})
	
	// General subagent for complex searches
	r.Register(&Info{
		Name:        "general",
		Description: "General subagent for complex searches and multistep tasks.",
		Mode:        ModeSubagent,
		Native:      true,
		Permission:  permission.ReadOnlyRuleset(),
		Options:     make(map[string]interface{}),
	})
}

// Register registers an agent
func (r *Registry) Register(agent *Info) {
	r.agents[agent.Name] = agent
}

// Get retrieves an agent by name
func (r *Registry) Get(name string) (*Info, bool) {
	agent, ok := r.agents[name]
	return agent, ok
}

// List returns all registered agents
func (r *Registry) List() []*Info {
	agents := make([]*Info, 0, len(r.agents))
	for _, agent := range r.agents {
		if !agent.Hidden {
			agents = append(agents, agent)
		}
	}
	return agents
}

// ListPrimary returns all primary agents
func (r *Registry) ListPrimary() []*Info {
	agents := make([]*Info, 0)
	for _, agent := range r.agents {
		if !agent.Hidden && (agent.Mode == ModePrimary || agent.Mode == ModeAll) {
			agents = append(agents, agent)
		}
	}
	return agents
}

// ListSubagents returns all subagents
func (r *Registry) ListSubagents() []*Info {
	agents := make([]*Info, 0)
	for _, agent := range r.agents {
		if !agent.Hidden && (agent.Mode == ModeSubagent || agent.Mode == ModeAll) {
			agents = append(agents, agent)
		}
	}
	return agents
}

// Default returns the default agent
func (r *Registry) Default() *Info {
	return r.agents["build"]
}
