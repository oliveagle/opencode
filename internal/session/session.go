package session

import (
	"time"

	"github.com/anomalyco/opencode/internal/id"
	"github.com/anomalyco/opencode/internal/permission"
)

// Info represents session metadata
type Info struct {
	ID        id.ID            `json:"id"`
	Slug      string           `json:"slug"`
	ProjectID string           `json:"projectId"`
	Directory string           `json:"directory"`
	ParentID  *id.ID           `json:"parentId,omitempty"`
	Summary   *Summary         `json:"summary,omitempty"`
	Share     *ShareInfo       `json:"share,omitempty"`
	Title     string           `json:"title"`
	Version   string           `json:"version"`
	Time      TimeInfo         `json:"time"`
	Permission *permission.Ruleset `json:"permission,omitempty"`
	Revert    *RevertInfo      `json:"revert,omitempty"`
}

// Summary contains session summary statistics
type Summary struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	Files     int `json:"files"`
}

// ShareInfo contains sharing information
type ShareInfo struct {
	URL string `json:"url"`
}

// TimeInfo contains timestamps
type TimeInfo struct {
	Created    int64 `json:"created"`
	Updated    int64 `json:"updated"`
	Compacting int64 `json:"compacting,omitempty"`
	Archived   int64 `json:"archived,omitempty"`
}

// RevertInfo contains revert information
type RevertInfo struct {
	MessageID string `json:"messageId"`
	PartID    string `json:"partId,omitempty"`
	Snapshot  string `json:"snapshot,omitempty"`
	Diff      string `json:"diff,omitempty"`
}

// New creates a new session
func New(projectID, directory string) *Info {
	now := time.Now().UnixMilli()
	return &Info{
		ID:        id.New(id.PrefixSession),
		Slug:      generateSlug(),
		ProjectID: projectID,
		Directory: directory,
		Title:     "New session - " + time.Now().UTC().Format(time.RFC3339),
		Version:   "1",
		Time: TimeInfo{
			Created: now,
			Updated: now,
		},
	}
}

// IsDefaultTitle checks if the title is a default generated title
func IsDefaultTitle(title string) bool {
	// Check if title matches default pattern
	return len(title) > 14 && title[:14] == "New session - "
}

// GetForkedTitle generates a forked title
func GetForkedTitle(title string) string {
	// Simple fork numbering
	return title + " (fork #1)"
}

func generateSlug() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
