package project

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/anomalyco/opencode/internal/id"
)

// Project represents a project/workspace
type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
	Settings  Settings  `json:"settings,omitempty"`
}

// Settings represents project-specific settings
type Settings struct {
	DefaultModel string            `json:"defaultModel,omitempty"`
	Instructions []string          `json:"instructions,omitempty"`
	Env          map[string]string `json:"env,omitempty"`
}

// Manager manages projects
type Manager struct {
	projects map[string]*Project
	mu       sync.RWMutex
}

// NewManager creates a new project manager
func NewManager() *Manager {
	return &Manager{
		projects: make(map[string]*Project),
	}
}

// Create creates a new project from a directory
func (m *Manager) Create(path string) (*Project, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Check if project already exists
	for _, p := range m.projects {
		if p.Path == abs {
			return p, nil
		}
	}

	// Create new project
	project := &Project{
		ID:        id.New(id.PrefixProject).String(),
		Name:      filepath.Base(abs),
		Path:      abs,
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
		Settings:  Settings{Env: make(map[string]string)},
	}

	m.projects[project.ID] = project
	return project, nil
}

// Get retrieves a project by ID
func (m *Manager) Get(id string) (*Project, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.projects[id]
	return p, ok
}

// GetByPath retrieves a project by path
func (m *Manager) GetByPath(path string) (*Project, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, false
	}

	for _, p := range m.projects {
		if p.Path == abs {
			return p, true
		}
	}

	// Create new project if not exists
	project, err := m.Create(abs)
	if err != nil {
		return nil, false
	}
	return project, true
}

// List returns all projects
func (m *Manager) List() []*Project {
	m.mu.RLock()
	defer m.mu.RUnlock()

	projects := make([]*Project, 0, len(m.projects))
	for _, p := range m.projects {
		projects = append(projects, p)
	}
	return projects
}

// Delete removes a project
func (m *Manager) Delete(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.projects, id)
}

// Update updates a project
func (m *Manager) Update(project *Project) {
	m.mu.Lock()
	defer m.mu.Unlock()
	project.UpdatedAt = time.Now().UnixMilli()
	m.projects[project.ID] = project
}

// DetectProjectType detects the project type based on files present
func DetectProjectType(path string) string {
	// Check for common project files
	indicators := map[string]string{
		"go.mod":          "go",
		"package.json":    "node",
		"Cargo.toml":      "rust",
		"pyproject.toml":  "python",
		"requirements.txt": "python",
		"Gemfile":         "ruby",
		"pom.xml":         "java",
		"build.gradle":    "java",
		"composer.json":   "php",
		"mix.exs":         "elixir",
		"Project.toml":    "julia",
	}

	for file, projectType := range indicators {
		if _, err := os.Stat(filepath.Join(path, file)); err == nil {
			return projectType
		}
	}

	return "unknown"
}

// FindGitRoot finds the root of the git repository
func FindGitRoot(path string) string {
	dir, err := filepath.Abs(path)
	if err != nil {
		return ""
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// Instance represents a project instance with state
type Instance struct {
	Project   *Project
	SessionID string
	State     map[string]interface{}
	mu        sync.RWMutex
}

// NewInstance creates a new project instance
func NewInstance(project *Project) *Instance {
	return &Instance{
		Project: project,
		State:   make(map[string]interface{}),
	}
}

// SetState sets a state value
func (i *Instance) SetState(key string, value interface{}) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.State[key] = value
}

// GetState gets a state value
func (i *Instance) GetState(key string) (interface{}, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	v, ok := i.State[key]
	return v, ok
}

// Context key for project instance
type ctxKey struct{}

// WithInstance adds an instance to context
func WithInstance(ctx context.Context, instance *Instance) context.Context {
	return context.WithValue(ctx, ctxKey{}, instance)
}

// FromContext retrieves instance from context
func FromContext(ctx context.Context) (*Instance, bool) {
	instance, ok := ctx.Value(ctxKey{}).(*Instance)
	return instance, ok
}
