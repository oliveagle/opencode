package tui

import (
	"testing"

	"github.com/rivo/tview"
)

func TestSplitView_CreateHorizontal(t *testing.T) {
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	sv := NewSplitView(SplitHorizontal, primary, secondary, 0.5)

	if sv == nil {
		t.Error("SplitView should be created")
	}

	if sv.GetRatio() != 0.5 {
		t.Errorf("Expected ratio 0.5, got %v", sv.GetRatio())
	}
}

func TestSplitView_CreateVertical(t *testing.T) {
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	sv := NewSplitView(SplitVertical, primary, secondary, 0.6)

	if sv.GetRatio() != 0.6 {
		t.Errorf("Expected ratio 0.6, got %v", sv.GetRatio())
	}
}

func TestSplitView_InvalidRatio(t *testing.T) {
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	// Should default to 0.5 for invalid ratios
	sv := NewSplitView(SplitHorizontal, primary, secondary, 1.5)

	if sv.GetRatio() != 0.5 {
		t.Errorf("Expected ratio 0.5 for invalid input, got %v", sv.GetRatio())
	}
}

func TestSplitView_SetRatio(t *testing.T) {
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	sv := NewSplitView(SplitHorizontal, primary, secondary, 0.5)
	sv.SetRatio(0.7)

	if sv.GetRatio() != 0.7 {
		t.Errorf("Expected ratio 0.7, got %v", sv.GetRatio())
	}
}

func TestSplitView_SwapPanes(t *testing.T) {
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	sv := NewSplitView(SplitHorizontal, primary, secondary, 0.5)

	if sv.GetPrimary() != primary {
		t.Error("Primary pane should be primary")
	}

	sv.SwapPanes()

	if sv.GetPrimary() != secondary {
		t.Error("Primary pane should be swapped with secondary")
	}
}

func TestSplitViewManager_CreateLayout(t *testing.T) {
	m := NewSplitViewManager()
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	err := m.CreateLayout("test", SplitHorizontal, primary, secondary)

	if err != nil {
		t.Errorf("CreateLayout failed: %v", err)
	}

	if layout := m.GetLayout("test"); layout == nil {
		t.Error("Layout not found after creation")
	}
}

func TestSplitViewManager_CreateDuplicateLayout(t *testing.T) {
	m := NewSplitViewManager()
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	m.CreateLayout("test", SplitHorizontal, primary, secondary)
	err := m.CreateLayout("test", SplitHorizontal, primary, secondary)

	if err == nil {
		t.Error("CreateLayout should fail for duplicate names")
	}
}

func TestSplitViewManager_SetActive(t *testing.T) {
	m := NewSplitViewManager()
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	m.CreateLayout("layout1", SplitHorizontal, primary, secondary)
	m.CreateLayout("layout2", SplitVertical, primary, secondary)

	err := m.SetActive("layout2")

	if err != nil {
		t.Errorf("SetActive failed: %v", err)
	}

	if active := m.GetActiveLayout(); active == nil {
		t.Error("Active layout should be set")
	}
}

func TestSplitViewManager_DeleteLayout(t *testing.T) {
	m := NewSplitViewManager()
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	m.CreateLayout("test1", SplitHorizontal, primary, secondary)
	m.CreateLayout("test2", SplitHorizontal, primary, secondary)

	// Set test2 as active so we can delete test1
	m.SetActive("test2")

	err := m.DeleteLayout("test1")

	if err != nil {
		t.Errorf("DeleteLayout failed: %v", err)
	}

	if layout := m.GetLayout("test1"); layout != nil {
		t.Error("Layout should be deleted")
	}
}

func TestSplitViewManager_DeleteActiveLayout(t *testing.T) {
	m := NewSplitViewManager()
	primary := tview.NewTextView()
	secondary := tview.NewTextView()

	m.CreateLayout("test", SplitHorizontal, primary, secondary)
	err := m.DeleteLayout("test")

	if err == nil {
		t.Error("DeleteLayout should fail for active layout")
	}
}

func TestPaneManager_AddPane(t *testing.T) {
	pm := NewPaneManager()
	pane1 := tview.NewTextView()
	pane2 := tview.NewTextView()

	pm.AddPane(pane1)
	pm.AddPane(pane2)

	if pm.Count() != 2 {
		t.Errorf("Expected 2 panes, got %d", pm.Count())
	}
}

func TestPaneManager_NextPane(t *testing.T) {
	pm := NewPaneManager()
	pane1 := tview.NewTextView()
	pane2 := tview.NewTextView()

	pm.AddPane(pane1)
	pm.AddPane(pane2)

	next := pm.NextPane()
	if next != pane2 {
		t.Error("NextPane should return second pane")
	}

	// Should wrap around
	next = pm.NextPane()
	if next != pane1 {
		t.Error("NextPane should wrap around to first pane")
	}
}

func TestPaneManager_PrevPane(t *testing.T) {
	pm := NewPaneManager()
	pane1 := tview.NewTextView()
	pane2 := tview.NewTextView()

	pm.AddPane(pane1)
	pm.AddPane(pane2)

	// Move to second pane
	pm.NextPane()

	prev := pm.PrevPane()
	if prev != pane1 {
		t.Error("PrevPane should return first pane")
	}
}

func TestPaneManager_RemovePane(t *testing.T) {
	pm := NewPaneManager()
	pane1 := tview.NewTextView()
	pane2 := tview.NewTextView()

	pm.AddPane(pane1)
	pm.AddPane(pane2)

	err := pm.RemovePane(0)

	if err != nil {
		t.Errorf("RemovePane failed: %v", err)
	}

	if pm.Count() != 1 {
		t.Errorf("Expected 1 pane after removal, got %d", pm.Count())
	}
}

func TestPaneManager_InvalidIndex(t *testing.T) {
	pm := NewPaneManager()
	pane := tview.NewTextView()

	pm.AddPane(pane)

	err := pm.SetActivePane(10)

	if err == nil {
		t.Error("SetActivePane should fail for invalid index")
	}
}
