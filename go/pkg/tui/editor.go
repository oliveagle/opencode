package tui

import (
	"strings"

	"github.com/rivo/tview"
)

type Editor struct {
	view    *tview.TextArea
	path    string
	content string
	dirty   bool
}

func NewEditor() *Editor {
	ta := tview.NewTextArea()
	ta.SetBorder(true).SetTitle("Editor")
	ta.SetPlaceholder("No file open. Ctrl+O: Open file")

	return &Editor{
		view: ta,
	}
}

func (e *Editor) SetContent(content, path string) {
	e.content = content
	e.path = path
	e.dirty = false

	e.view.SetText(content, true)
}

func (e *Editor) GetContent() string {
	return e.view.GetText()
}

func (e *Editor) GetPath() string {
	return e.path
}

func (e *Editor) View() *tview.TextArea {
	return e.view
}

func (e *Editor) SetDirty(dirty bool) {
	e.dirty = dirty
	title := "Editor"
	if dirty {
		title = "Editor *"
	}
	e.view.SetTitle(title)
}

func (e *Editor) IsDirty() bool {
	return e.dirty
}

func (e *Editor) Clear() {
	e.content = ""
	e.path = ""
	e.dirty = false
	e.view.SetText("", true)
	e.view.SetTitle("Editor")
}

// SetHighlighting enables syntax highlighting for common languages
func (e *Editor) SetHighlighting(language string) {
	// tview's TextArea doesn't have built-in syntax highlighting
	// This is a placeholder for future enhancement
}

// GetLineCount returns number of lines in editor
func (e *Editor) GetLineCount() int {
	text := e.view.GetText()
	return len(strings.Split(text, "\n"))
}
