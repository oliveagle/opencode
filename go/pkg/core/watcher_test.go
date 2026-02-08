package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileWatcher(t *testing.T) {
	tmpDir := t.TempDir()
	watcher, err := NewFileWatcher()
	if err != nil {
		t.Fatalf("NewFileWatcher failed: %v", err)
	}
	defer watcher.Stop()

	if err := watcher.Watch(tmpDir); err != nil {
		t.Fatalf("Watch failed: %v", err)
	}

	watcher.Start()

	// Wait for watcher to start
	time.Sleep(100 * time.Millisecond)

	// Create a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
		t.Fatalf("Write file failed: %v", err)
	}

	// Wait for event
	timeout := time.NewTimer(2 * time.Second)
	defer timeout.Stop()

	select {
	case event := <-watcher.Events():
		if event.Type != "create" && event.Type != "modify" {
			t.Errorf("Wrong event type: %s, want create or modify", event.Type)
		}
	case <-timeout.C:
		t.Error("Did not receive create event")
	}

	// Modify the file
	if err := os.WriteFile(testFile, []byte("hello world"), 0644); err != nil {
		t.Fatalf("Modify file failed: %v", err)
	}

	timeout.Reset(2 * time.Second)
	select {
	case event := <-watcher.Events():
		if event.Type != "modify" {
			t.Errorf("Wrong event type: want modify, got %s", event.Type)
		}
	case <-timeout.C:
		t.Error("Did not receive modify event")
	}

	// Delete the file
	if err := os.Remove(testFile); err != nil {
		t.Fatalf("Delete file failed: %v", err)
	}

	timeout.Reset(2 * time.Second)
	eventCount := 0
	foundDelete := false

	// fsnotify may send multiple events on delete
	for {
		select {
		case event := <-watcher.Events():
			eventCount++
			if event.Type == "delete" {
				foundDelete = true
			}
			if eventCount >= 1 {
				break
			}
		case <-timeout.C:
			break
		}
		if eventCount > 0 {
			break
		}
	}

	if !foundDelete && eventCount == 0 {
		t.Error("Did not receive delete or modify event")
	}
}

func TestFileWatcherIgnoresTempFiles(t *testing.T) {
	tmpDir := t.TempDir()
	watcher, err := NewFileWatcher()
	if err != nil {
		t.Fatalf("NewFileWatcher failed: %v", err)
	}
	defer watcher.Stop()

	if err := watcher.Watch(tmpDir); err != nil {
		t.Fatalf("Watch failed: %v", err)
	}

	watcher.Start()
	time.Sleep(100 * time.Millisecond)

	// Create a temp file (should be ignored)
	tempFile := filepath.Join(tmpDir, "test.swp")
	if err := os.WriteFile(tempFile, []byte("temp"), 0644); err != nil {
		t.Fatalf("Write temp file failed: %v", err)
	}

	timeout := time.NewTimer(500 * time.Millisecond)
	defer timeout.Stop()

	select {
	case <-watcher.Events():
		t.Error("Should have ignored temp file")
	case <-timeout.C:
		// Expected - no event
	}
}
