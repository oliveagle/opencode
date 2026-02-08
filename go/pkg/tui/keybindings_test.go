package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestKeyBindingSet_Register(t *testing.T) {
	ks := NewKeyBindingSet("test")

	ks.RegisterRune(ActionUndo, 'z', tcell.ModCtrl, "Undo")
	binding := ks.Get(ActionUndo)

	if binding == nil {
		t.Error("KeyBinding should be registered")
	}

	if binding.Description != "Undo" {
		t.Errorf("Expected description 'Undo', got %q", binding.Description)
	}
}

func TestKeyBindingSet_GetByRune(t *testing.T) {
	ks := NewKeyBindingSet("test")

	ks.RegisterRune(ActionUndo, 'z', tcell.ModCtrl, "Undo")
	action := ks.GetByRune('z', tcell.ModCtrl)

	if action != ActionUndo {
		t.Errorf("Expected ActionUndo, got %q", action)
	}
}

func TestKeyBindingSet_ValidateNoConflicts(t *testing.T) {
	ks := NewKeyBindingSet("test")

	ks.RegisterRune(ActionUndo, 'z', tcell.ModCtrl, "Undo")
	ks.RegisterRune(ActionRedo, 'z', tcell.ModCtrl, "Redo") // Same key binding

	err := ks.ValidateNoConflicts()
	if err == nil {
		t.Error("ValidateNoConflicts should detect conflict")
	}
}

func TestDefaultKeyBindingSet(t *testing.T) {
	ks := DefaultKeyBindingSet()

	if ks == nil {
		t.Error("DefaultKeyBindingSet returned nil")
	}

	if ks.Get(ActionUndo) == nil {
		t.Error("Default set should have ActionUndo")
	}

	if len(ks.List()) == 0 {
		t.Error("Default set should have bindings")
	}
}

func TestEmacsKeyBindingSet(t *testing.T) {
	ks := EmacsKeyBindingSet()

	if ks == nil {
		t.Error("EmacsKeyBindingSet returned nil")
	}

	if ks.Get(ActionUndo) == nil {
		t.Error("Emacs set should have ActionUndo")
	}
}

func TestVimKeyBindingSet(t *testing.T) {
	ks := VimKeyBindingSet()

	if ks == nil {
		t.Error("VimKeyBindingSet returned nil")
	}

	if ks.Get(ActionUndo) == nil {
		t.Error("Vim set should have ActionUndo")
	}
}

func TestKeyBindingManager_Register(t *testing.T) {
	m := NewKeyBindingManager()

	custom := NewKeyBindingSet("custom")
	custom.RegisterRune(ActionUndo, 'x', tcell.ModCtrl, "Undo")

	err := m.Register(custom)
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}

	if retrieved := m.sets["custom"]; retrieved == nil {
		t.Error("Custom keybinding set not registered")
	}
}

func TestKeyBindingManager_Activate(t *testing.T) {
	m := NewKeyBindingManager()

	err := m.Activate("emacs")
	if err != nil {
		t.Errorf("Activate failed: %v", err)
	}

	if m.GetActiveName() != "emacs" {
		t.Errorf("Expected 'emacs', got %q", m.GetActiveName())
	}
}

func TestKeyBindingManager_ActivateInvalid(t *testing.T) {
	m := NewKeyBindingManager()

	err := m.Activate("nonexistent")
	if err == nil {
		t.Error("Activate should fail for nonexistent set")
	}
}

func TestKeyBindingManager_ListSets(t *testing.T) {
	m := NewKeyBindingManager()

	sets := m.ListSets()
	if len(sets) < 3 {
		t.Errorf("Expected at least 3 sets, got %d", len(sets))
	}

	// Check if default sets are present
	hasDefault := false
	hasEmacs := false
	hasVim := false

	for _, name := range sets {
		if name == "default" {
			hasDefault = true
		}
		if name == "emacs" {
			hasEmacs = true
		}
		if name == "vim" {
			hasVim = true
		}
	}

	if !hasDefault || !hasEmacs || !hasVim {
		t.Error("Not all default keybinding sets found")
	}
}
