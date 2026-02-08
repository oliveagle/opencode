package core

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/anomalyco/opencode/pkg/types"
)

type ProjectService struct {
	configFile string
	projects   map[string]*types.ProjectConfig
}

func NewProjectService(configDir string) (*ProjectService, error) {
	configFile := filepath.Join(configDir, "projects.json")
	ps := &ProjectService{
		configFile: configFile,
		projects:   make(map[string]*types.ProjectConfig),
	}

	if err := ps.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return ps, nil
}

func (p *ProjectService) CreateProject(cfg *types.ProjectConfig) error {
	p.projects[cfg.Name] = cfg
	return p.save()
}

func (p *ProjectService) GetProject(name string) *types.ProjectConfig {
	return p.projects[name]
}

func (p *ProjectService) ListProjects() []*types.ProjectConfig {
	var result []*types.ProjectConfig
	for _, proj := range p.projects {
		result = append(result, proj)
	}
	return result
}

func (p *ProjectService) DeleteProject(name string) error {
	delete(p.projects, name)
	return p.save()
}

func (p *ProjectService) load() error {
	data, err := os.ReadFile(p.configFile)
	if err != nil {
		return err
	}

	var configs []*types.ProjectConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return err
	}

	for _, cfg := range configs {
		p.projects[cfg.Name] = cfg
	}
	return nil
}

func (p *ProjectService) save() error {
	var configs []*types.ProjectConfig
	for _, cfg := range p.projects {
		configs = append(configs, cfg)
	}

	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(p.configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(p.configFile, data, 0644)
}
