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

	// Phase 2 enhancements
	syntax       *SyntaxHighlighter
	undoRedo     *UndoRedoStack
	search       *SearchReplace
	showLineNums bool
	currentLine  int
	language     string
}

func NewEditor() *Editor {
	ta := tview.NewTextArea()
	ta.SetBorder(true).SetTitle("Editor")
	ta.SetPlaceholder("No file open. Ctrl+O: Open file")

	return &Editor{
		view:         ta,
		undoRedo:     NewUndoRedoStack(500),
		search:       NewSearchReplace(),
		showLineNums: true,
		currentLine:  0,
		language:     "text",
	}
}

func (e *Editor) SetContent(content, path string) {
	e.content = content
	e.path = path
	e.dirty = false
	e.currentLine = 0

	e.view.SetText(content, true)

	// Auto-detect language from file extension
	e.detectLanguage(path)
	e.syntax = NewSyntaxHighlighter(e.language)

	// Clear undo/redo for new file
	e.undoRedo.Clear()
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
	e.language = language
	e.syntax = NewSyntaxHighlighter(language)
}

// GetLineCount returns number of lines in editor
func (e *Editor) GetLineCount() int {
	text := e.view.GetText()
	return len(strings.Split(text, "\n"))
}

func (e *Editor) detectLanguage(path string) {
	if path == "" {
		e.language = "text"
		return
	}

	ext := ""
	if idx := strings.LastIndex(path, "."); idx != -1 {
		ext = strings.ToLower(path[idx+1:])
	}

	switch ext {
	case "go":
		e.language = "go"
	case "py":
		e.language = "python"
	case "js":
		e.language = "javascript"
	case "ts":
		e.language = "typescript"
	case "rs":
		e.language = "rust"
	case "yaml", "yml":
		e.language = "yaml"
	case "json":
		e.language = "json"
	default:
		e.language = "text"
	}
}

func (e *Editor) Undo() {
	if e.undoRedo.CanUndo() {
		op := e.undoRedo.Undo()
		if op != nil {
			// Apply undo operation
			text := e.view.GetText()
			if op.Type == "insert" {
				// Remove inserted text
				before := text[:op.Offset]
				after := text[op.Offset+len(op.Text):]
				e.view.SetText(before+after, true)
			} else if op.Type == "delete" {
				// Re-insert deleted text
				before := text[:op.Offset]
				after := text[op.Offset:]
				e.view.SetText(before+op.Text+after, true)
			}
			e.dirty = true
		}
	}
}

func (e *Editor) Redo() {
	if e.undoRedo.CanRedo() {
		op := e.undoRedo.Redo()
		if op != nil {
			// Apply redo operation
			text := e.view.GetText()
			if op.Type == "insert" {
				// Re-insert text
				before := text[:op.Offset]
				after := text[op.Offset:]
				e.view.SetText(before+op.Text+after, true)
			} else if op.Type == "delete" {
				// Remove text again
				before := text[:op.Offset]
				after := text[op.Offset+len(op.Text):]
				e.view.SetText(before+after, true)
			}
			e.dirty = true
		}
	}
}

func (e *Editor) FindAll(pattern string, caseSensitive bool) []SearchMatch {
	e.search.SetCaseSensitive(caseSensitive)
	e.search.SetPattern(pattern, false)
	return e.search.FindAll(e.view.GetText())
}

func (e *Editor) ReplaceAll(pattern, replacement string, regex bool) int {
	e.search.SetPattern(pattern, regex)
	original := e.view.GetText()
	replaced := e.search.ReplaceAll(original, replacement)

	if original != replaced {
		e.view.SetText(replaced, true)
		e.dirty = true
		return strings.Count(original, pattern) // Approximate
	}
	return 0
}

func (e *Editor) GetSyntax() *SyntaxHighlighter {
	return e.syntax
}

func (e *Editor) GetLanguage() string {
	return e.language
}

func (e *Editor) CanUndo() bool {
	return e.undoRedo.CanUndo()
}

func (e *Editor) CanRedo() bool {
	return e.undoRedo.CanRedo()
}
