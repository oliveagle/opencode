package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// Prefix represents the type prefix for an identifier
type Prefix string

const (
	PrefixSession   Prefix = "session"
	PrefixMessage   Prefix = "msg"
	PrefixProject   Prefix = "project"
	PrefixProvider  Prefix = "provider"
	PrefixTool      Prefix = "tool"
	PrefixAgent     Prefix = "agent"
	PrefixSnapshot  Prefix = "snapshot"
)

// ID represents a typed identifier
type ID struct {
	Prefix Prefix
	Value  string
}

// New generates a new unique ID with the given prefix
func New(prefix Prefix) ID {
	timestamp := time.Now().UnixMicro()
	random := make([]byte, 4)
	rand.Read(random)
	return ID{
		Prefix: prefix,
		Value:  fmt.Sprintf("%d_%s", timestamp, hex.EncodeToString(random)),
	}
}

// Parse parses an ID string
func Parse(s string) (ID, error) {
	for _, prefix := range []Prefix{PrefixSession, PrefixMessage, PrefixProject, PrefixProvider, PrefixTool, PrefixAgent, PrefixSnapshot} {
		prefixStr := string(prefix) + "_"
		if len(s) > len(prefixStr) && s[:len(prefixStr)] == prefixStr {
			return ID{Prefix: prefix, Value: s[len(prefixStr):]}, nil
		}
	}
	return ID{}, fmt.Errorf("invalid ID format: %s", s)
}

// String returns the string representation of the ID
func (id ID) String() string {
	return string(id.Prefix) + "_" + id.Value
}

// Schema returns a schema-compatible ID
func Schema(prefix Prefix) ID {
	return New(prefix)
}
