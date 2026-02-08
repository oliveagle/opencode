package core

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileEvent represents a file system change event
type FileEvent struct {
	Type  string // "create", "modify", "delete", "rename"
	Path  string
	IsDir bool
	Time  time.Time
}

// FileWatcher watches for file system changes
type FileWatcher struct {
	watcher *fsnotify.Watcher
	events  chan FileEvent
	stop    chan struct{}
	mu      sync.RWMutex
	paths   map[string]bool
	running bool
	batch   []FileEvent
	batchMu sync.Mutex
}

// NewFileWatcher creates a new FileWatcher
func NewFileWatcher() (*FileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		watcher: w,
		events:  make(chan FileEvent, 100),
		stop:    make(chan struct{}),
		paths:   make(map[string]bool),
		batch:   make([]FileEvent, 0, 10),
	}

	return fw, nil
}

// Watch adds a path to watch for changes
func (fw *FileWatcher) Watch(path string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.paths[path] {
		return nil // already watching
	}

	if err := fw.watcher.Add(path); err != nil {
		return err
	}

	fw.paths[path] = true
	log.Printf("Watching: %s", path)
	return nil
}

// Start begins watching for changes
func (fw *FileWatcher) Start() {
	fw.mu.Lock()
	if fw.running {
		fw.mu.Unlock()
		return
	}
	fw.running = true
	fw.mu.Unlock()

	go fw.run()
}

// Stop stops watching for changes
func (fw *FileWatcher) Stop() {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if !fw.running {
		return
	}

	fw.running = false
	close(fw.stop)
	fw.watcher.Close()
}

// Events returns the channel with file events
func (fw *FileWatcher) Events() <-chan FileEvent {
	return fw.events
}

func (fw *FileWatcher) run() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// Filter out temporary and hidden files
			if fw.shouldIgnore(event.Name) {
				continue
			}

			fe := fw.convertEvent(event)
			fw.addToBatch(fe)

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)

		case <-ticker.C:
			fw.flushBatch()

		case <-fw.stop:
			fw.flushBatch()
			return
		}
	}
}

func (fw *FileWatcher) convertEvent(event fsnotify.Event) FileEvent {
	fe := FileEvent{
		Path: event.Name,
		Time: time.Now(),
	}

	// Determine event type
	if event.Has(fsnotify.Create) {
		fe.Type = "create"
	} else if event.Has(fsnotify.Write) {
		fe.Type = "modify"
	} else if event.Has(fsnotify.Remove) {
		fe.Type = "delete"
	} else if event.Has(fsnotify.Rename) {
		fe.Type = "rename"
	} else if event.Has(fsnotify.Chmod) {
		fe.Type = "chmod"
	}

	return fe
}

func (fw *FileWatcher) addToBatch(event FileEvent) {
	fw.batchMu.Lock()
	defer fw.batchMu.Unlock()

	// Avoid duplicates in batch window
	for _, e := range fw.batch {
		if e.Path == event.Path && e.Type == event.Type {
			return
		}
	}

	fw.batch = append(fw.batch, event)
}

func (fw *FileWatcher) flushBatch() {
	fw.batchMu.Lock()
	defer fw.batchMu.Unlock()

	if len(fw.batch) == 0 {
		return
	}

	for _, event := range fw.batch {
		select {
		case fw.events <- event:
		case <-fw.stop:
			return
		}
	}

	fw.batch = fw.batch[:0]
}

func (fw *FileWatcher) shouldIgnore(path string) bool {
	// Ignore temporary files
	name := filepath.Base(path)

	ignoredPatterns := []string{
		".git",
		".swp", ".swo", "~",
		".tmp",
		"node_modules",
	}

	for _, pattern := range ignoredPatterns {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
		if matched, _ := filepath.Match("*"+pattern, name); matched {
			return true
		}
	}

	return false
}
