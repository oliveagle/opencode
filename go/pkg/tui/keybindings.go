package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type KeyBindingAction string

const (
	// Editor actions
	ActionUndo      KeyBindingAction = "undo"
	ActionRedo      KeyBindingAction = "redo"
	ActionFind      KeyBindingAction = "find"
	ActionReplace   KeyBindingAction = "replace"
	ActionSave      KeyBindingAction = "save"
	ActionQuit      KeyBindingAction = "quit"
	ActionNewFile   KeyBindingAction = "newfile"

	// File tree actions
	ActionOpenFile  KeyBindingAction = "openfile"
	ActionDeleteFile KeyBindingAction = "deletefile"
	ActionCreateFile KeyBindingAction = "createfile"
	ActionRefresh   KeyBindingAction = "refresh"

	// Global actions
	ActionCommandPalette KeyBindingAction = "command_palette"

	// Split actions
	ActionSplitHorizontal KeyBindingAction = "split_h"
	ActionSplitVertical   KeyBindingAction = "split_v"
	ActionCloseSplit      KeyBindingAction = "close_split"
	ActionSwitchPane      KeyBindingAction = "switch_pane"
)

type KeyBinding struct {
	Action      KeyBindingAction
	Keys        []tcell.Key
	Runes       []rune
	Modifiers   tcell.ModMask
	Description string
}

type KeyBindingSet struct {
	name        string
	bindings    map[KeyBindingAction]*KeyBinding
	actionToKey map[KeyBindingAction]string // For display
}

func NewKeyBindingSet(name string) *KeyBindingSet {
	return &KeyBindingSet{
		name:        name,
		bindings:    map[KeyBindingAction]*KeyBinding{},
		actionToKey: map[KeyBindingAction]string{},
	}
}

func (k *KeyBindingSet) Register(action KeyBindingAction, key tcell.Key, mod tcell.ModMask, desc string) {
	k.bindings[action] = &KeyBinding{
		Action:      action,
		Keys:        []tcell.Key{key},
		Modifiers:   mod,
		Description: desc,
	}
	k.actionToKey[action] = k.encodeKey(key, mod)
}

func (k *KeyBindingSet) RegisterRune(action KeyBindingAction, r rune, mod tcell.ModMask, desc string) {
	k.bindings[action] = &KeyBinding{
		Action:      action,
		Runes:       []rune{r},
		Modifiers:   mod,
		Description: desc,
	}
	k.actionToKey[action] = k.encodeRune(r, mod)
}

func (k *KeyBindingSet) Get(action KeyBindingAction) *KeyBinding {
	return k.bindings[action]
}

func (k *KeyBindingSet) GetByKey(key tcell.Key, mod tcell.ModMask) KeyBindingAction {
	for action, binding := range k.bindings {
		if len(binding.Keys) > 0 && binding.Keys[0] == key && binding.Modifiers == mod {
			return action
		}
	}
	return ""
}

func (k *KeyBindingSet) GetByRune(r rune, mod tcell.ModMask) KeyBindingAction {
	for action, binding := range k.bindings {
		if len(binding.Runes) > 0 && binding.Runes[0] == r && binding.Modifiers == mod {
			return action
		}
	}
	return ""
}

func (k *KeyBindingSet) List() []*KeyBinding {
	var bindings []*KeyBinding
	for _, b := range k.bindings {
		bindings = append(bindings, b)
	}

	sort.Slice(bindings, func(i, j int) bool {
		return bindings[i].Description < bindings[j].Description
	})

	return bindings
}

func (k *KeyBindingSet) ValidateNoConflicts() error {
	used := map[string]bool{}

	for action, binding := range k.bindings {
		key := k.encodeBinding(binding)
		if used[key] {
			return fmt.Errorf("key binding conflict: %s for action %s", key, action)
		}
		used[key] = true
	}

	return nil
}

func (k *KeyBindingSet) encodeBinding(b *KeyBinding) string {
	if len(b.Keys) > 0 {
		return k.encodeKey(b.Keys[0], b.Modifiers)
	}
	if len(b.Runes) > 0 {
		return k.encodeRune(b.Runes[0], b.Modifiers)
	}
	return ""
}

func (k *KeyBindingSet) encodeKey(key tcell.Key, mod tcell.ModMask) string {
	parts := []string{}

	if mod&tcell.ModCtrl != 0 {
		parts = append(parts, "Ctrl")
	}
	if mod&tcell.ModAlt != 0 {
		parts = append(parts, "Alt")
	}
	if mod&tcell.ModShift != 0 {
		parts = append(parts, "Shift")
	}

	// Convert key to string representation
	keyStr := ""
	switch key {
	case tcell.KeyUp:
		keyStr = "Up"
	case tcell.KeyDown:
		keyStr = "Down"
	case tcell.KeyLeft:
		keyStr = "Left"
	case tcell.KeyRight:
		keyStr = "Right"
	case tcell.KeyEnter:
		keyStr = "Enter"
	case tcell.KeyTab:
		keyStr = "Tab"
	case tcell.KeyCtrlC:
		keyStr = "CtrlC"
	case tcell.KeyCtrlQ:
		keyStr = "CtrlQ"
	case tcell.KeyEscape:
		keyStr = "Esc"
	case tcell.KeyDelete:
		keyStr = "Delete"
	default:
		keyStr = "Unknown"
	}

	parts = append(parts, keyStr)
	return strings.Join(parts, "+")
}

func (k *KeyBindingSet) encodeRune(r rune, mod tcell.ModMask) string {
	parts := []string{}

	if mod&tcell.ModCtrl != 0 {
		parts = append(parts, "Ctrl")
	}
	if mod&tcell.ModAlt != 0 {
		parts = append(parts, "Alt")
	}
	if mod&tcell.ModShift != 0 {
		parts = append(parts, "Shift")
	}

	parts = append(parts, string(r))
	return strings.Join(parts, "+")
}

// DefaultKeyBindingSet returns the default VIM-like keybindings
func DefaultKeyBindingSet() *KeyBindingSet {
	ks := NewKeyBindingSet("default")

	// Editor keybindings
	ks.RegisterRune(ActionUndo, 'z', tcell.ModCtrl, "Undo (Ctrl+Z)")
	ks.RegisterRune(ActionRedo, 'Z', tcell.ModCtrl, "Redo (Ctrl+Shift+Z)")
	ks.RegisterRune(ActionFind, 'f', tcell.ModCtrl, "Find (Ctrl+F)")
	ks.RegisterRune(ActionReplace, 'h', tcell.ModCtrl, "Replace (Ctrl+H)")
	ks.RegisterRune(ActionSave, 's', tcell.ModCtrl, "Save (Ctrl+S)")
	ks.RegisterRune(ActionQuit, 'q', tcell.ModCtrl, "Quit (Ctrl+Q)")
	ks.RegisterRune(ActionNewFile, 'n', tcell.ModCtrl, "New File (Ctrl+N)")

	// File tree keybindings
	ks.RegisterRune(ActionOpenFile, 'o', tcell.ModCtrl, "Open File (Ctrl+O)")
	ks.RegisterRune(ActionDeleteFile, 'd', tcell.ModCtrl, "Delete File (Ctrl+D)")
	ks.RegisterRune(ActionCreateFile, 'c', tcell.ModCtrl, "Create File (Ctrl+C)")

	// Global keybindings
	ks.RegisterRune(ActionCommandPalette, 'k', tcell.ModCtrl, "Command Palette (Ctrl+K)")

	// Split window keybindings
	ks.RegisterRune(ActionSplitHorizontal, 'j', tcell.ModAlt, "Split Horizontal (Alt+J)")
	ks.RegisterRune(ActionSplitVertical, 'l', tcell.ModAlt, "Split Vertical (Alt+L)")
	ks.RegisterRune(ActionSwitchPane, 'w', tcell.ModCtrl, "Switch Pane (Ctrl+W)")

	return ks
}

// EmacsKeyBindingSet returns Emacs-like keybindings
func EmacsKeyBindingSet() *KeyBindingSet {
	ks := NewKeyBindingSet("emacs")

	// Similar to default but with Emacs conventions
	ks.RegisterRune(ActionUndo, '_', tcell.ModCtrl, "Undo (Ctrl+_)")
	ks.RegisterRune(ActionFind, 's', tcell.ModCtrl, "Find (Ctrl+S)")
	ks.RegisterRune(ActionReplace, 'h', tcell.ModCtrl, "Replace (Ctrl+H)")
	ks.RegisterRune(ActionSave, 's', tcell.ModCtrl|tcell.ModShift, "Save (Ctrl+Shift+S)")
	ks.RegisterRune(ActionQuit, 'c', tcell.ModCtrl, "Quit (Ctrl+C)")
	ks.RegisterRune(ActionNewFile, 'n', tcell.ModCtrl, "New File (Ctrl+N)")
	ks.RegisterRune(ActionCommandPalette, 'p', tcell.ModCtrl, "Command Palette (Ctrl+P)")

	return ks
}

// VimKeyBindingSet returns VIM-like keybindings
func VimKeyBindingSet() *KeyBindingSet {
	ks := NewKeyBindingSet("vim")

	// VIM-style keybindings (using : for command mode simulation)
	ks.RegisterRune(ActionUndo, 'u', tcell.ModCtrl, "Undo (Ctrl+U)")
	ks.RegisterRune(ActionRedo, 'r', tcell.ModCtrl, "Redo (Ctrl+R)")
	ks.RegisterRune(ActionFind, '/', 0, "Find (/)")
	ks.RegisterRune(ActionReplace, '&', 0, "Replace (&)")
	ks.RegisterRune(ActionSave, 'w', 0, "Save (:w)")
	ks.RegisterRune(ActionQuit, 'q', 0, "Quit (:q)")
	ks.RegisterRune(ActionCommandPalette, ':', 0, "Command Palette (:)")

	return ks
}

// KeyBindingManager manages all keybinding sets
type KeyBindingManager struct {
	sets    map[string]*KeyBindingSet
	active  *KeyBindingSet
}

func NewKeyBindingManager() *KeyBindingManager {
	m := &KeyBindingManager{
		sets:  map[string]*KeyBindingSet{},
	}

	// Register default sets
	m.Register(DefaultKeyBindingSet())
	m.Register(EmacsKeyBindingSet())
	m.Register(VimKeyBindingSet())

	// Set default as active
	m.active = m.Get("default")

	return m
}

func (m *KeyBindingManager) Register(set *KeyBindingSet) error {
	if err := set.ValidateNoConflicts(); err != nil {
		return err
	}
	m.sets[set.name] = set
	return nil
}

func (m *KeyBindingManager) Get(name string) *KeyBindingSet {
	return m.sets[name]
}

func (m *KeyBindingManager) Activate(name string) error {
	set, ok := m.sets[name]
	if !ok {
		return fmt.Errorf("keybinding set %q not found", name)
	}
	m.active = set
	return nil
}

func (m *KeyBindingManager) GetActive() *KeyBindingSet {
	return m.active
}

func (m *KeyBindingManager) GetActiveName() string {
	return m.active.name
}

func (m *KeyBindingManager) ListSets() []string {
	var names []string
	for name := range m.sets {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (m *KeyBindingManager) GetActionForEvent(event *tcell.EventKey) KeyBindingAction {
	active := m.active
	if event.Key() != tcell.KeyRune {
		return active.GetByKey(event.Key(), event.Modifiers())
	}
	return active.GetByRune(event.Rune(), event.Modifiers())
}
