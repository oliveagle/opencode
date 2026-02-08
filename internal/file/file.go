package file

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/anomalyco/opencode/internal/util/log"
)

// IgnoreMatcher matches files to ignore
type IgnoreMatcher interface {
	Match(path string, isDir bool) bool
}

// Manager handles file operations
type Manager struct {
	rootDir string
	ignore  IgnoreMatcher
	watcher *Watcher
	mu      sync.RWMutex
}

// NewManager creates a new file manager
func NewManager(rootDir string) *Manager {
	return &Manager{
		rootDir: rootDir,
		ignore:  NewGitIgnore(rootDir),
	}
}

// SetIgnore sets the ignore matcher
func (m *Manager) SetIgnore(ignore IgnoreMatcher) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ignore = ignore
}

// Read reads a file's contents
func (m *Manager) Read(path string) (string, error) {
	fullPath := m.resolvePath(path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadLines reads specific lines from a file
func (m *Manager) ReadLines(path string, offset, limit int) (string, error) {
	fullPath := m.resolvePath(path)
	file, err := os.Open(fullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum < offset {
			continue
		}
		if limit > 0 && len(lines) >= limit {
			break
		}
		lines = append(lines, scanner.Text())
	}

	return strings.Join(lines, "\n"), scanner.Err()
}

// Write writes content to a file
func (m *Manager) Write(path, content string) error {
	fullPath := m.resolvePath(path)
	
	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}

// Edit edits a file by replacing old string with new string
func (m *Manager) Edit(path, oldString, newString string) (bool, error) {
	content, err := m.Read(path)
	if err != nil {
		return false, err
	}

	// Check if old string exists
	if !strings.Contains(content, oldString) {
		return false, nil
	}

	newContent := strings.Replace(content, oldString, newString, 1)
	return true, m.Write(path, newContent)
}

// Delete deletes a file
func (m *Manager) Delete(path string) error {
	fullPath := m.resolvePath(path)
	return os.Remove(fullPath)
}

// Exists checks if a file exists
func (m *Manager) Exists(path string) bool {
	fullPath := m.resolvePath(path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// IsDir checks if a path is a directory
func (m *Manager) IsDir(path string) bool {
	fullPath := m.resolvePath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// List lists directory contents
func (m *Manager) List(path string) ([]fs.DirEntry, error) {
	fullPath := m.resolvePath(path)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	// Filter ignored files
	m.mu.RLock()
	ignore := m.ignore
	m.mu.RUnlock()

	var result []fs.DirEntry
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if ignore != nil && ignore.Match(entryPath, entry.IsDir()) {
			continue
		}
		result = append(result, entry)
	}

	return result, nil
}

// Walk walks the directory tree
func (m *Manager) Walk(root string, fn filepath.WalkFunc) error {
	fullPath := m.resolvePath(root)

	m.mu.RLock()
	ignore := m.ignore
	m.mu.RUnlock()

	return filepath.Walk(fullPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(fullPath, path)
		if relPath == "." {
			return fn(path, info, nil)
		}

		if ignore != nil && ignore.Match(relPath, info.IsDir()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		return fn(path, info, nil)
	})
}

// Glob finds files matching a pattern
func (m *Manager) Glob(pattern string) ([]string, error) {
	fullPattern := m.resolvePath(pattern)
	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		return nil, err
	}

	// Make paths relative
	var result []string
	for _, match := range matches {
		rel, _ := filepath.Rel(m.rootDir, match)
		result = append(result, rel)
	}

	return result, nil
}

// Grep searches for a pattern in files
func (m *Manager) Grep(ctx context.Context, pattern string, paths []string, caseSensitive bool) ([]GrepMatch, error) {
	m.mu.RLock()
	ignore := m.ignore
	m.mu.RUnlock()

	var matches []GrepMatch

	for _, searchPath := range paths {
		fullPath := m.resolvePath(searchPath)

		err := filepath.Walk(fullPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				relPath, _ := filepath.Rel(m.rootDir, path)
				if ignore != nil && ignore.Match(relPath, true) {
					return filepath.SkipDir
				}
				return nil
			}

			// Skip non-text files
			if !isTextFile(path) {
				return nil
			}

			// Search in file
			fileMatches, err := grepFile(path, pattern, caseSensitive)
			if err != nil {
				return nil // Continue on error
			}

			for _, m := range fileMatches {
				relPath, _ := filepath.Rel(m.rootDir, path)
				m.Path = relPath
				matches = append(matches, m)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return matches, nil
}

// Watch starts watching for file changes
func (m *Manager) Watch(ctx context.Context, handler WatchHandler) error {
	m.mu.Lock()
	if m.watcher != nil {
		m.watcher.Stop()
	}
	m.watcher = NewWatcher(m.rootDir, m.ignore)
	m.mu.Unlock()

	return m.watcher.Start(ctx, handler)
}

// StopWatch stops watching for changes
func (m *Manager) StopWatch() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.watcher != nil {
		m.watcher.Stop()
		m.watcher = nil
	}
}

func (m *Manager) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(m.rootDir, path)
}

// GrepMatch represents a grep match
type GrepMatch struct {
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Content string `json:"content"`
}

func grepFile(path, pattern string, caseSensitive bool) ([]GrepMatch, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	searchPattern := pattern
	if !caseSensitive {
		searchPattern = strings.ToLower(pattern)
	}

	var matches []GrepMatch
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		searchLine := line
		if !caseSensitive {
			searchLine = strings.ToLower(line)
		}

		if strings.Contains(searchLine, searchPattern) {
			col := strings.Index(searchLine, searchPattern)
			matches = append(matches, GrepMatch{
				Line:    lineNum,
				Column:  col + 1,
				Content: line,
			})
		}
	}

	return matches, scanner.Err()
}

func isTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	textExts := map[string]bool{
		".txt": true, ".md": true, ".json": true, ".yaml": true, ".yml": true,
		".xml": true, ".html": true, ".css": true, ".js": true, ".ts": true,
		".jsx": true, ".tsx": true, ".go": true, ".rs": true, ".py": true,
		".java": true, ".c": true, ".cpp": true, ".h": true, ".hpp": true,
		".sh": true, ".bash": true, ".zsh": true, ".fish": true,
		".toml": true, ".ini": true, ".cfg": true, ".conf": true,
		".gitignore": true, ".dockerignore": true, ".env": true,
		".sql": true, ".proto": true, ".graphql": true, ".gql": true,
	}
	return textExts[ext]
}

// WatchHandler handles file change events
type WatchHandler func(event WatchEvent)

// WatchEvent represents a file change event
type WatchEvent struct {
	Type    WatchEventType
	Path    string
	OldPath string // For rename events
}

type WatchEventType string

const (
	WatchEventCreate WatchEventType = "create"
	WatchEventModify WatchEventType = "modify"
	WatchEventDelete WatchEventType = "delete"
	WatchEventRename WatchEventType = "rename"
)

// Watcher watches for file changes
type Watcher struct {
	rootDir string
	ignore  IgnoreMatcher
	stop    chan struct{}
}

// NewWatcher creates a new file watcher
func NewWatcher(rootDir string, ignore IgnoreMatcher) *Watcher {
	return &Watcher{
		rootDir: rootDir,
		ignore:  ignore,
		stop:    make(chan struct{}),
	}
}

// Start starts watching for changes
func (w *Watcher) Start(ctx context.Context, handler WatchHandler) error {
	log.Default.Info("Starting file watcher", map[string]interface{}{
		"root": w.rootDir,
	})

	// Simple polling-based watcher
	// TODO: Use fsnotify or similar for better performance
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	knownFiles := make(map[string]os.FileInfo)

	// Initial scan
	filepath.Walk(w.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		relPath, _ := filepath.Rel(w.rootDir, path)
		if w.ignore != nil && w.ignore.Match(relPath, info.IsDir()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		knownFiles[path] = info
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.stop:
			return nil
		case <-ticker.C:
			// Check for changes
			currentFiles := make(map[string]os.FileInfo)

			filepath.Walk(w.rootDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				relPath, _ := filepath.Rel(w.rootDir, path)
				if w.ignore != nil && w.ignore.Match(relPath, info.IsDir()) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				currentFiles[path] = info

				// Check if new or modified
				if oldInfo, exists := knownFiles[path]; exists {
					if info.ModTime().After(oldInfo.ModTime()) {
						handler(WatchEvent{
							Type: WatchEventModify,
							Path: relPath,
						})
					}
				} else {
					handler(WatchEvent{
						Type: WatchEventCreate,
						Path: relPath,
					})
				}

				return nil
			})

			// Check for deleted files
			for path := range knownFiles {
				if _, exists := currentFiles[path]; !exists {
					relPath, _ := filepath.Rel(w.rootDir, path)
					handler(WatchEvent{
						Type: WatchEventDelete,
						Path: relPath,
					})
				}
			}

			knownFiles = currentFiles
		}
	}
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	close(w.stop)
}
