package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

type SplitDirection int

const (
	SplitHorizontal SplitDirection = iota
	SplitVertical
)

type SplitView struct {
	root       *tview.Flex
	direction  SplitDirection
	primary    tview.Primitive
	secondary  tview.Primitive
	ratio      float64 // 0.0 to 1.0, primary takes this proportion
	resizable  bool
	dragActive bool
	dragStart  int
}

func NewSplitView(direction SplitDirection, primary, secondary tview.Primitive, ratio float64) *SplitView {
	if ratio < 0.1 || ratio > 0.9 {
		ratio = 0.5 // Default to 50/50
	}

	flex := tview.NewFlex()

	sv := &SplitView{
		root:      flex,
		direction: direction,
		primary:   primary,
		secondary: secondary,
		ratio:     ratio,
		resizable: true,
	}

	sv.updateLayout()
	return sv
}

func (sv *SplitView) updateLayout() {
	sv.root.Clear()

	if sv.direction == SplitHorizontal {
		sv.root.SetDirection(tview.FlexRow)
	} else {
		sv.root.SetDirection(tview.FlexColumn)
	}

	// Convert ratio (0.0-1.0) to tview flex units
	primarySize := int(float64(100) * sv.ratio)

	sv.root.AddItem(sv.primary, primarySize, 1, true)
	sv.root.AddItem(sv.secondary, 100-primarySize, 1, false)
}

func (sv *SplitView) SetResizable(resizable bool) {
	sv.resizable = resizable
}

func (sv *SplitView) GetView() *tview.Flex {
	return sv.root
}

func (sv *SplitView) SetRatio(ratio float64) {
	if ratio < 0.1 || ratio > 0.9 {
		return
	}
	sv.ratio = ratio
	sv.updateLayout()
}

func (sv *SplitView) GetRatio() float64 {
	return sv.ratio
}

func (sv *SplitView) SwapPanes() {
	sv.primary, sv.secondary = sv.secondary, sv.primary
	sv.updateLayout()
}

func (sv *SplitView) GetPrimary() tview.Primitive {
	return sv.primary
}

func (sv *SplitView) GetSecondary() tview.Primitive {
	return sv.secondary
}

func (sv *SplitView) SetPrimary(p tview.Primitive) {
	sv.primary = p
	sv.updateLayout()
}

func (sv *SplitView) SetSecondary(s tview.Primitive) {
	sv.secondary = s
	sv.updateLayout()
}

// LayoutConfig represents a saved layout configuration
type LayoutConfig struct {
	Name          string                 `json:"name"`
	Type          string                 `json:"type"` // "single", "hsplit", "vsplit", "grid"
	SplitRatio    float64                `json:"split_ratio"`
	PrimaryTitle  string                 `json:"primary_title"`
	SecondaryTitle string                `json:"secondary_title"`
	SavedAt       int64                  `json:"saved_at"`
}

// SplitViewManager manages multiple split configurations
type SplitViewManager struct {
	activeLayout string
	layouts      map[string]*SplitView
	configs      map[string]*LayoutConfig
}

func NewSplitViewManager() *SplitViewManager {
	return &SplitViewManager{
		layouts: map[string]*SplitView{},
		configs: map[string]*LayoutConfig{},
	}
}

func (m *SplitViewManager) CreateLayout(name string, direction SplitDirection, primary, secondary tview.Primitive) error {
	if _, exists := m.layouts[name]; exists {
		return fmt.Errorf("layout %q already exists", name)
	}

	split := NewSplitView(direction, primary, secondary, 0.5)
	m.layouts[name] = split
	m.activeLayout = name

	return nil
}

func (m *SplitViewManager) GetLayout(name string) *SplitView {
	return m.layouts[name]
}

func (m *SplitViewManager) GetActiveLayout() *SplitView {
	return m.layouts[m.activeLayout]
}

func (m *SplitViewManager) SetActive(name string) error {
	if _, exists := m.layouts[name]; !exists {
		return fmt.Errorf("layout %q not found", name)
	}
	m.activeLayout = name
	return nil
}

func (m *SplitViewManager) DeleteLayout(name string) error {
	if name == m.activeLayout {
		return fmt.Errorf("cannot delete active layout")
	}
	if _, exists := m.layouts[name]; !exists {
		return fmt.Errorf("layout %q not found", name)
	}
	delete(m.layouts, name)
	delete(m.configs, name)
	return nil
}

func (m *SplitViewManager) SaveConfig(name string, config *LayoutConfig) {
	config.Name = name
	m.configs[name] = config
}

func (m *SplitViewManager) GetConfig(name string) *LayoutConfig {
	return m.configs[name]
}

func (m *SplitViewManager) ListLayouts() []string {
	var names []string
	for name := range m.layouts {
		names = append(names, name)
	}
	return names
}

// PaneManager tracks multiple panes for focus management
type PaneManager struct {
	panes      []tview.Primitive
	activePaneIndex int
}

func NewPaneManager() *PaneManager {
	return &PaneManager{
		panes: []tview.Primitive{},
	}
}

func (p *PaneManager) AddPane(pane tview.Primitive) {
	p.panes = append(p.panes, pane)
	if len(p.panes) == 1 {
		p.activePaneIndex = 0
	}
}

func (p *PaneManager) NextPane() tview.Primitive {
	if len(p.panes) == 0 {
		return nil
	}

	p.activePaneIndex = (p.activePaneIndex + 1) % len(p.panes)
	return p.panes[p.activePaneIndex]
}

func (p *PaneManager) PrevPane() tview.Primitive {
	if len(p.panes) == 0 {
		return nil
	}

	p.activePaneIndex--
	if p.activePaneIndex < 0 {
		p.activePaneIndex = len(p.panes) - 1
	}

	return p.panes[p.activePaneIndex]
}

func (p *PaneManager) GetActivePane() tview.Primitive {
	if len(p.panes) > p.activePaneIndex {
		return p.panes[p.activePaneIndex]
	}
	return nil
}

func (p *PaneManager) SetActivePane(index int) error {
	if index < 0 || index >= len(p.panes) {
		return fmt.Errorf("pane index %d out of range", index)
	}
	p.activePaneIndex = index
	return nil
}

func (p *PaneManager) RemovePane(index int) error {
	if index < 0 || index >= len(p.panes) {
		return fmt.Errorf("pane index %d out of range", index)
	}

	p.panes = append(p.panes[:index], p.panes[index+1:]...)

	if p.activePaneIndex >= len(p.panes) {
		p.activePaneIndex = len(p.panes) - 1
	}

	return nil
}

func (p *PaneManager) Count() int {
	return len(p.panes)
}
