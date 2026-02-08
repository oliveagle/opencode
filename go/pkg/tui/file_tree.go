package tui

import (
	"fmt"

	"github.com/anomalyco/opencode/pkg/api"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FileTree struct {
	view        *tview.TreeView
	client      *api.Client
	currentPath string
	onSelected  func(path string, isDir bool)
}

func NewFileTree(client *api.Client) *FileTree {
	tree := tview.NewTreeView()
	ft := &FileTree{
		view:        tree,
		client:      client,
		currentPath: ".",
	}

	// Root node
	root := tview.NewTreeNode(".").
		SetColor(tcell.ColorGreen)
	tree.SetRoot(root)
	tree.SetCurrentNode(root)

	// Load initial directory
	ft.loadDirectory(root, ".")

	// Selection handler
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if ft.onSelected != nil {
			path := node.GetReference().(string)
			// Check if it's a directory
			isDir := len(node.GetChildren()) > 0
			ft.onSelected(path, isDir)
		}
	})

	return ft
}

func (ft *FileTree) loadDirectory(node *tview.TreeNode, path string) {
	files, err := ft.client.ListDir(path)
	if err != nil {
		node.SetText(fmt.Sprintf("Error: %v", err))
		return
	}

	// Clear existing children
	node.ClearChildren()

	// Add files to tree
	for _, f := range files {
		color := tcell.ColorWhite
		if f.IsDir {
			color = tcell.ColorBlue
		}

		child := tview.NewTreeNode(f.Name).
			SetColor(color).
			SetReference(f.Path).
			SetSelectable(true)

		if f.IsDir {
			child.SetChildren([]*tview.TreeNode{})
		}

		node.AddChild(child)
	}
}

func (ft *FileTree) View() *tview.TreeView {
	return ft.view
}

func (ft *FileTree) SetSelectionChangedFunc(fn func(path string, isDir bool)) {
	ft.onSelected = fn
}

func (ft *FileTree) Refresh() {
	root := ft.view.GetRoot()
	if root != nil {
		ft.loadDirectory(root, ".")
	}
}
