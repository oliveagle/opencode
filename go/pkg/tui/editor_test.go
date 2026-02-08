package tui

import (
	"testing"
)

func TestSyntaxHighlighter_Go(t *testing.T) {
	sh := NewSyntaxHighlighter("go")

	code := `package main
func test() {
	x := 42
}`

	colors := sh.ColorizeText(code)
	if len(colors) == 0 {
		t.Error("SyntaxHighlighter returned no colors")
	}
}

func TestSyntaxHighlighter_Python(t *testing.T) {
	sh := NewSyntaxHighlighter("python")

	code := `def hello():
	x = 42
	return x`

	colors := sh.ColorizeText(code)
	if len(colors) == 0 {
		t.Error("SyntaxHighlighter returned no colors")
	}
}

func TestSyntaxHighlighter_JSON(t *testing.T) {
	sh := NewSyntaxHighlighter("json")

	code := `{"name": "test", "count": 42}`

	colors := sh.ColorizeText(code)
	if len(colors) == 0 {
		t.Error("SyntaxHighlighter returned no colors")
	}
}

func TestUndoRedoStack(t *testing.T) {
	stack := NewUndoRedoStack(10)

	if stack.CanUndo() {
		t.Error("New stack should not be able to undo")
	}

	stack.Push(&EditOperation{Type: "insert", Text: "hello", Offset: 0})

	if !stack.CanUndo() {
		t.Error("Stack should be able to undo after push")
	}

	op := stack.Undo()
	if op.Type != "insert" || op.Text != "hello" {
		t.Error("Undo did not return correct operation")
	}

	if !stack.CanRedo() {
		t.Error("Stack should be able to redo")
	}

	op = stack.Redo()
	if op.Type != "insert" || op.Text != "hello" {
		t.Error("Redo did not return correct operation")
	}
}

func TestUndoRedoStack_MaxSize(t *testing.T) {
	stack := NewUndoRedoStack(3)

	stack.Push(&EditOperation{Type: "insert", Text: "1"})
	stack.Push(&EditOperation{Type: "insert", Text: "2"})
	stack.Push(&EditOperation{Type: "insert", Text: "3"})
	stack.Push(&EditOperation{Type: "insert", Text: "4"})

	if stack.UndoCount() > 3 {
		t.Error("Stack exceeded maximum size")
	}
}

func TestSearchReplace_PlainText(t *testing.T) {
	sr := NewSearchReplace()
	sr.SetPattern("hello", false)

	text := "hello world hello"
	matches := sr.FindAll(text)

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}

	if matches[0].Text != "hello" {
		t.Error("Match text incorrect")
	}
}

func TestSearchReplace_CaseInsensitive(t *testing.T) {
	sr := NewSearchReplace()
	sr.SetPattern("Hello", false)
	sr.SetCaseSensitive(false)

	text := "hello world HELLO"
	matches := sr.FindAll(text)

	if len(matches) != 2 {
		t.Errorf("Expected 2 case-insensitive matches, got %d", len(matches))
	}
}

func TestSearchReplace_Regex(t *testing.T) {
	sr := NewSearchReplace()
	sr.SetPattern(`\d+`, true)

	text := "I have 42 apples and 3 oranges"
	matches := sr.FindAll(text)

	if len(matches) != 2 {
		t.Errorf("Expected 2 regex matches, got %d", len(matches))
	}
}

func TestSearchReplace_Replace(t *testing.T) {
	sr := NewSearchReplace()
	sr.SetPattern("world", false)

	text := "hello world"
	result := sr.Replace(text, "OpenCode")

	expected := "hello OpenCode"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSearchReplace_ReplaceAll(t *testing.T) {
	sr := NewSearchReplace()
	sr.SetPattern("o", false)

	text := "foo boo"
	result := sr.ReplaceAll(text, "0")

	expected := "f00 b00"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEditor_LanguageDetection(t *testing.T) {
	e := NewEditor()

	tests := []struct {
		path     string
		expected string
	}{
		{"main.go", "go"},
		{"script.py", "python"},
		{"index.js", "javascript"},
		{"app.ts", "typescript"},
		{"config.yaml", "yaml"},
		{"data.json", "json"},
		{"test.rs", "rust"},
		{"README.md", "text"},
	}

	for _, test := range tests {
		e.detectLanguage(test.path)
		if e.language != test.expected {
			t.Errorf("Path %q: expected %q, got %q", test.path, test.expected, e.language)
		}
	}
}

func TestEditor_UndoRedo(t *testing.T) {
	e := NewEditor()
	e.SetContent("hello", "test.txt")

	if e.CanUndo() {
		t.Error("Editor should not have undo initially")
	}

	// Note: In real use, operations would be pushed via the stack
	// This is a basic test of the methods
}
