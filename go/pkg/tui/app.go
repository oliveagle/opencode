package tui

import (
	"context"
	"fmt"
	"log"

	"github.com/anomalyco/opencode/pkg/api"
	"github.com/anomalyco/opencode/pkg/core"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	root       *tview.Application
	client     *api.Client
	fileTree   *FileTree
	editor     *Editor
	cmdPalette *CommandPalette
	statusBar  *StatusBar
	modal      *tview.Modal

	ctx     context.Context
	cancel  context.CancelFunc
	layout  *tview.Flex
	
	// File watching
	watcher  *core.FileWatcher
	eventBus *EventBus
	basePath string
}

func NewApp(serverURL string) (*App, error) {
	root := tview.NewApplication()
	client := api.NewClient(serverURL)

	// Test server connection
	if err := client.Health(); err != nil {
		return nil, fmt.Errorf("cannot connect to server: %w", err)
	}

	// Initialize file watcher
	watcher, err := core.NewFileWatcher()
	if err != nil {
		return nil, fmt.Errorf("cannot create file watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	a := &App{
		root:       root,
		client:     client,
		ctx:        ctx,
		cancel:     cancel,
		fileTree:   NewFileTree(client),
		editor:     NewEditor(),
		cmdPalette: NewCommandPalette(),
		statusBar:  NewStatusBar(),
		watcher:    watcher,
		eventBus:   NewEventBus(),
		basePath:   ".",
	}

	// Build layout
	a.buildLayout()
	a.setupKeyBindings()
	a.setupFileWatching()

	return a, nil
}

func (a *App) buildLayout() {
	// Main content area (split between file tree and editor)
	mainContent := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.fileTree.View(), 30, 0, false).
		AddItem(a.editor.View(), 0, 1, true)

	// Vertical split (content + status)
	a.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(createHeader(), 1, 0, false).
		AddItem(mainContent, 0, 1, true).
		AddItem(a.statusBar.View(), 1, 0, false)

	a.root.SetRoot(a.layout, true)
}

func (a *App) setupKeyBindings() {
	a.root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		modifiers := event.Modifiers()

		if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyCtrlQ {
			a.root.Stop()
			return nil
		}

		if event.Rune() == 'k' && modifiers&tcell.ModCtrl != 0 {
			a.showCommandPalette()
			return nil
		}

		if event.Rune() == 's' && modifiers&tcell.ModCtrl != 0 {
			a.saveCurrentFile()
			return nil
		}

		if event.Rune() == 'n' && modifiers&tcell.ModCtrl != 0 {
			a.showNewFileDialog()
			return nil
		}

		return event
	})

	a.fileTree.SetSelectionChangedFunc(func(path string, isDir bool) {
		if !isDir {
			a.loadFile(path)
		}
	})
}

func (a *App) setupFileWatching() {
	a.watcher.Watch(a.basePath)
	a.watcher.Start()

	// Handle file events in background
	go func() {
		for event := range a.watcher.Events() {
			switch event.Type {
			case "create", "delete", "rename":
				// Refresh file tree
				a.eventBus.Publish(Event{
					Type: EventFileCreated,
					Path: event.Path,
				})
				a.fileTree.Refresh()

			case "modify":
				// Check if it's the current file
				if event.Path == a.editor.GetPath() {
					a.eventBus.Publish(Event{
						Type: EventFileModified,
						Path: event.Path,
					})
					a.statusBar.SetStatus(fmt.Sprintf("File modified externally: %s", event.Path))
				}
			}
		}
	}()
}

func (a *App) loadFile(path string) {
	content, err := a.client.ReadFile(path)
	if err != nil {
		a.statusBar.SetStatus(fmt.Sprintf("Error: %v", err))
		return
	}

	a.editor.SetContent(content, path)
	a.statusBar.SetStatus(fmt.Sprintf("Loaded: %s", path))
}

func (a *App) saveCurrentFile() {
	path := a.editor.GetPath()
	if path == "" {
		a.statusBar.SetStatus("No file open")
		return
	}

	content := a.editor.GetContent()
	if err := a.client.WriteFile(path, content); err != nil {
		a.statusBar.SetStatus(fmt.Sprintf("Save failed: %v", err))
		return
	}

	a.statusBar.SetStatus(fmt.Sprintf("Saved: %s", path))
}

func (a *App) showCommandPalette() {
	a.root.SetFocus(a.cmdPalette.View())
	a.statusBar.SetStatus("Command palette (ESC to close)")
}

func (a *App) showNewFileDialog() {
	modal := tview.NewModal().
		SetText("New file name:").
		AddButtons([]string{"Create", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.root.SetRoot(a.layout, true)
		})

	a.root.SetRoot(modal, false)
}

func (a *App) Run() error {
	log.Println("Starting OpenCode TUI")
	return a.root.Run()
}

func (a *App) Stop() {
	a.cancel()
	a.root.Stop()
}

func createHeader() *tview.TextView {
	header := tview.NewTextView().
		SetText(" OpenCode - TUI v0.1.0 | Ctrl+K: Commands | Ctrl+S: Save | Ctrl+Q: Quit").
		SetTextAlign(tview.AlignLeft)
	return header
}

