package tui

import (
	"log"

	"github.com/rivo/tview"
)

type App struct {
	root *tview.Application
}

func NewApp() *App {
	app := tview.NewApplication()
	
	return &App{
		root: app,
	}
}

func (a *App) Run() error {
	log.Println("Starting opencode TUI")
	
	// Create main layout
	mainView := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(createHeader(), 1, 0, false).
		AddItem(createMainPane(), 0, 1, true).
		AddItem(createFooter(), 1, 0, false)
	
	if err := a.root.SetRoot(mainView, true).Run(); err != nil {
		return err
	}
	return nil
}

func createHeader() *tview.TextView {
	header := tview.NewTextView().
		SetText("OpenCode TUI - v0.1.0").
		SetTextAlign(tview.AlignCenter)
	header.SetBorder(true).SetTitle("Header")
	return header
}

func createMainPane() *tview.Flex {
	// TODO: implement main content pane with file explorer, editor, etc.
	content := tview.NewTextView().
		SetText("Main content area - TUI under development")
	return tview.NewFlex().AddItem(content, 0, 1, true)
}

func createFooter() *tview.TextView {
	footer := tview.NewTextView().
		SetText("Status bar - Ready")
	return footer
}
