package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/anomalyco/opencode/internal/installation"
)

// Storage provides persistent key-value storage
type Storage struct {
	dir string
	mu  sync.RWMutex
}

// New creates a new storage instance
func New(name string) *Storage {
	dir := filepath.Join(installation.DataDir(), name)
	os.MkdirAll(dir, 0755)
	return &Storage{dir: dir}
}

// Get retrieves a value
func (s *Storage) Get(key string, dest interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.path(key)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Set stores a value
func (s *Storage) Set(key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.path(key)
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Delete removes a value
func (s *Storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.path(key)
	return os.Remove(path)
}

// Exists checks if a key exists
func (s *Storage) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.path(key)
	_, err := os.Stat(path)
	return err == nil
}

// List lists all keys
func (s *Storage) List() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			keys = append(keys, entry.Name()[:len(entry.Name())-5])
		}
	}
	return keys, nil
}

func (s *Storage) path(key string) string {
	return filepath.Join(s.dir, key+".json")
}

// SessionStorage provides session-specific storage
type SessionStorage struct {
	*Storage
}

// NewSessionStorage creates session storage
func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		Storage: New("sessions"),
	}
}

// GlobalStorage provides global storage
type GlobalStorage struct {
	*Storage
}

// NewGlobalStorage creates global storage
func NewGlobalStorage() *GlobalStorage {
	return &GlobalStorage{
		Storage: New("global"),
	}
}

// CacheStorage provides cache storage
type CacheStorage struct {
	*Storage
}

// NewCacheStorage creates cache storage
func NewCacheStorage() *CacheStorage {
	return &CacheStorage{
		Storage: New("cache"),
	}
}
