package permission

// Mode represents the permission mode
type Mode string

const (
	ModeAllow    Mode = "allow"
	ModeDeny     Mode = "deny"
	ModeAsk      Mode = "ask"
	ModeReadOnly Mode = "readonly"
)

// Ruleset represents a set of permission rules
type Ruleset struct {
	Mode       Mode    `json:"mode"`
	AllowList  []Rule  `json:"allowList,omitempty"`
	DenyList   []Rule  `json:"denyList,omitempty"`
	AskList    []Rule  `json:"askList,omitempty"`
}

// Rule represents a permission rule
type Rule struct {
	Type    string `json:"type"`    // tool, file, command
	Pattern string `json:"pattern"` // glob pattern or tool name
	Reason  string `json:"reason,omitempty"`
}

// DefaultRuleset returns the default permission ruleset
func DefaultRuleset() *Ruleset {
	return &Ruleset{
		Mode: ModeAsk,
		AllowList: []Rule{
			{Type: "tool", Pattern: "read"},
			{Type: "tool", Pattern: "ls"},
			{Type: "tool", Pattern: "grep"},
			{Type: "tool", Pattern: "glob"},
		},
		AskList: []Rule{
			{Type: "tool", Pattern: "bash"},
			{Type: "tool", Pattern: "write"},
			{Type: "tool", Pattern: "edit"},
		},
	}
}

// ReadOnlyRuleset returns a read-only permission ruleset
func ReadOnlyRuleset() *Ruleset {
	return &Ruleset{
		Mode: ModeReadOnly,
		AllowList: []Rule{
			{Type: "tool", Pattern: "read"},
			{Type: "tool", Pattern: "ls"},
			{Type: "tool", Pattern: "grep"},
			{Type: "tool", Pattern: "glob"},
		},
		DenyList: []Rule{
			{Type: "tool", Pattern: "bash"},
			{Type: "tool", Pattern: "write"},
			{Type: "tool", Pattern: "edit"},
		},
	}
}

// IsAllowed checks if a tool is allowed
func (r *Ruleset) IsAllowed(tool string) bool {
	// Check deny list first
	for _, rule := range r.DenyList {
		if rule.Type == "tool" && matchPattern(rule.Pattern, tool) {
			return false
		}
	}
	
	// Check allow list
	for _, rule := range r.AllowList {
		if rule.Type == "tool" && matchPattern(rule.Pattern, tool) {
			return true
		}
	}
	
	// Default based on mode
	switch r.Mode {
	case ModeAllow:
		return true
	case ModeDeny, ModeReadOnly:
		return false
	default:
		return false // Ask mode - not allowed by default
	}
}

// ShouldAsk checks if permission should be asked
func (r *Ruleset) ShouldAsk(tool string) bool {
	if r.Mode == ModeAsk {
		// Check if in ask list
		for _, rule := range r.AskList {
			if rule.Type == "tool" && matchPattern(rule.Pattern, tool) {
				return true
			}
		}
		// Not in allow list means ask
		for _, rule := range r.AllowList {
			if rule.Type == "tool" && matchPattern(rule.Pattern, tool) {
				return false
			}
		}
		return true
	}
	return false
}

func matchPattern(pattern, s string) bool {
	if pattern == "*" {
		return true
	}
	if pattern == s {
		return true
	}
	// Simple glob matching
	if len(pattern) > 0 && pattern[0] == '*' {
		suffix := pattern[1:]
		if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
			return true
		}
	}
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
