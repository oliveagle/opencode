package session

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/anomalyco/opencode/internal/id"
	"github.com/anomalyco/opencode/internal/storage"
)

// Manager manages sessions
type Manager struct {
	sessions map[string]*Info
	storage  *storage.Storage
	mu       sync.RWMutex
}

// NewManager creates a new session manager
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*Info),
		storage:  storage.NewSessionStorage(),
	}
}

// Create creates a new session
func (m *Manager) Create(projectID, directory string) *Info {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := New(projectID, directory)
	m.sessions[session.ID.String()] = session

	// Persist to storage
	m.storage.Set(session.ID.String(), session)

	return session
}

// Get retrieves a session by ID
func (m *Manager) Get(id string) (*Info, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[id]
	if !ok {
		// Try loading from storage
		var s Info
		if err := m.storage.Get(id, &s); err == nil {
			m.sessions[id] = &s
			return &s, true
		}
	}
	return session, ok
}

// List returns all sessions for a project
func (m *Manager) List(projectID string) []*Info {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var sessions []*Info
	for _, s := range m.sessions {
		if projectID == "" || s.ProjectID == projectID {
			sessions = append(sessions, s)
		}
	}
	return sessions
}

// Update updates a session
func (m *Manager) Update(session *Info) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session.Time.Updated = time.Now().UnixMilli()
	m.sessions[session.ID.String()] = session
	m.storage.Set(session.ID.String(), session)
}

// Delete removes a session
func (m *Manager) Delete(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, id)
	m.storage.Delete(id)
}

// Fork creates a fork of a session
func (m *Manager) Fork(id string) (*Info, error) {
	original, ok := m.Get(id)
	if !ok {
		return nil, ErrSessionNotFound
	}

	fork := &Info{
		ID:        id.New(id.PrefixSession),
		Slug:      generateSlug(),
		ProjectID: original.ProjectID,
		Directory: original.Directory,
		ParentID:  &original.ID,
		Title:     GetForkedTitle(original.Title),
		Version:   "1",
		Time: TimeInfo{
			Created: time.Now().UnixMilli(),
			Updated: time.Now().UnixMilli(),
		},
	}

	m.mu.Lock()
	m.sessions[fork.ID.String()] = fork
	m.mu.Unlock()

	m.storage.Set(fork.ID.String(), fork)

	return fork, nil
}

// Archive archives a session
func (m *Manager) Archive(id string) error {
	session, ok := m.Get(id)
	if !ok {
		return ErrSessionNotFound
	}

	session.Time.Archived = time.Now().UnixMilli()
	m.Update(session)
	return nil
}

// ConversationManager manages conversations within sessions
type ConversationManager struct {
	conversations map[string]*Conversation
	mu            sync.RWMutex
}

// NewConversationManager creates a new conversation manager
func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		conversations: make(map[string]*Conversation),
	}
}

// GetOrCreate gets or creates a conversation for a session
func (m *ConversationManager) GetOrCreate(sessionID string) *Conversation {
	m.mu.Lock()
	defer m.mu.Unlock()

	conv, ok := m.conversations[sessionID]
	if !ok {
		conv = &Conversation{
			SessionID: sessionID,
			Messages:  []*Message{},
		}
		m.conversations[sessionID] = conv
	}
	return conv
}

// AddMessage adds a message to a conversation
func (m *ConversationManager) AddMessage(sessionID string, msg *Message) {
	conv := m.GetOrCreate(sessionID)
	m.mu.Lock()
	defer m.mu.Unlock()
	conv.Messages = append(conv.Messages, msg)
}

// GetMessages returns all messages for a session
func (m *ConversationManager) GetMessages(sessionID string) []*Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conv, ok := m.conversations[sessionID]
	if !ok {
		return nil
	}
	return conv.Messages
}

// Clear clears all messages for a session
func (m *ConversationManager) Clear(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	conv, ok := m.conversations[sessionID]
	if ok {
		conv.Messages = []*Message{}
	}
}

// Error definitions
var (
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidSession  = errors.New("invalid session")
)
