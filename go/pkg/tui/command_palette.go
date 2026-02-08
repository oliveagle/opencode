package tui

import (
	"github.com/rivo/tview"
)

type CommandPalette struct {
	view     *tview.InputField
	commands map[string]func()
}

func NewCommandPalette() *CommandPalette {
	input := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(50).
		SetPlaceholder("Type command...")

	cp := &CommandPalette{
		view:     input,
		commands: make(map[string]func()),
	}

	cp.registerDefaultCommands()

	return cp
}

func (cp *CommandPalette) registerDefaultCommands() {
	cp.commands["new"] = func() {
		// TODO: implement new file
	}
	cp.commands["open"] = func() {
		// TODO: implement open file
	}
	cp.commands["save"] = func() {
		// TODO: implement save
	}
	cp.commands["search"] = func() {
		// TODO: implement search
	}
	cp.commands["format"] = func() {
		// TODO: implement format
	}
	cp.commands["quit"] = func() {
		// TODO: implement quit
	}
}

func (cp *CommandPalette) View() *tview.InputField {
	return cp.view
}

func (cp *CommandPalette) GetInput() string {
	return cp.view.GetText()
}

func (cp *CommandPalette) Clear() {
	cp.view.SetText("")
}

func (cp *CommandPalette) SetInputChangeFunc(fn func(text string)) {
	cp.view.SetChangedFunc(fn)
}

func (cp *CommandPalette) ExecuteCommand(cmd string) error {
	if fn, ok := cp.commands[cmd]; ok {
		fn()
	}
	return nil
}
