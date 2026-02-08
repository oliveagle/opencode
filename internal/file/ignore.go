package file

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GitIgnore implements ignore matching based on .gitignore rules
type GitIgnore struct {
	patterns []ignorePattern
	rootDir  string
}

type ignorePattern struct {
	pattern  string
	negate   bool
	isDir    bool
	rootOnly bool
}

// NewGitIgnore creates a new GitIgnore from a directory
func NewGitIgnore(rootDir string) *GitIgnore {
	gi := &GitIgnore{
		rootDir: rootDir,
	}
	gi.loadGitignore()
	return gi
}

func (gi *GitIgnore) loadGitignore() {
	// Load .gitignore from root
	gitignorePath := filepath.Join(gi.rootDir, ".gitignore")
	if data, err := os.ReadFile(gitignorePath); err == nil {
		gi.parseGitignore(string(data))
	}

	// Add default ignores
	defaultIgnores := []string{
		".git",
		".svn",
		".hg",
		"node_modules",
		"vendor",
		"*.swp",
		"*.swo",
		"*~",
		".DS_Store",
		"Thumbs.db",
		"*.pyc",
		"__pycache__",
		".idea",
		".vscode",
		"*.iml",
		"target",
		"dist",
		"build",
		"out",
		".opencode",
	}
	for _, pattern := range defaultIgnores {
		gi.parseLine(pattern)
	}
}

func (gi *GitIgnore) parseGitignore(content string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		gi.parseLine(line)
	}
}

func (gi *GitIgnore) parseLine(line string) {
	// Skip empty lines and comments
	if line == "" || strings.HasPrefix(line, "#") {
		return
	}

	pattern := ignorePattern{}

	// Check for negation
	if strings.HasPrefix(line, "!") {
		pattern.negate = true
		line = line[1:]
	}

	// Check for directory-only pattern
	if strings.HasSuffix(line, "/") {
		pattern.isDir = true
		line = strings.TrimSuffix(line, "/")
	}

	// Check for root-only pattern
	if strings.HasPrefix(line, "/") {
		pattern.rootOnly = true
		line = strings.TrimPrefix(line, "/")
	}

	pattern.pattern = line
	gi.patterns = append(gi.patterns, pattern)
}

// Match checks if a path should be ignored
func (gi *GitIgnore) Match(path string, isDir bool) bool {
	// Normalize path
	path = filepath.ToSlash(path)
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	matched := false
	for _, p := range gi.patterns {
		if p.isDir && !isDir {
			continue
		}

		if p.rootOnly {
			// Match only from root
			if gi.matchPattern(p.pattern, path, isDir) {
				matched = !p.negate
			}
		} else {
			// Match from any directory
			if gi.matchPattern(p.pattern, path, isDir) || gi.matchPatternAnyDir(p.pattern, path, isDir) {
				matched = !p.negate
			}
		}
	}

	return matched
}

func (gi *GitIgnore) matchPattern(pattern, path string, isDir bool) bool {
	// Simple pattern matching
	if pattern == path {
		return true
	}

	// Handle ** patterns
	if strings.Contains(pattern, "**") {
		parts := strings.Split(pattern, "**")
		if len(parts) == 2 {
			prefix := parts[0]
			suffix := parts[1]
			if prefix == "" && suffix == "" {
				return true
			}
			if prefix != "" && !strings.HasPrefix(path, prefix) {
				return false
			}
			if suffix != "" && !strings.HasSuffix(path, suffix) {
				return false
			}
			return true
		}
	}

	// Handle * patterns
	if strings.Contains(pattern, "*") {
		return gi.matchGlob(pattern, path)
	}

	// Handle extension patterns
	if strings.HasPrefix(pattern, "*.") {
		ext := pattern[1:] // *.go -> .go
		return strings.HasSuffix(path, ext)
	}

	// Exact match or prefix for directories
	if isDir && strings.HasPrefix(path, pattern+"/") {
		return true
	}

	return false
}

func (gi *GitIgnore) matchPatternAnyDir(pattern, path string, isDir bool) bool {
	// Try matching from any subdirectory
	parts := strings.Split(path, "/")
	for i := 0; i < len(parts); i++ {
		subPath := strings.Join(parts[i:], "/")
		if gi.matchPattern(pattern, subPath, isDir) {
			return true
		}
	}
	return false
}

func (gi *GitIgnore) matchGlob(pattern, path string) bool {
	// Simple glob matching
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	return gi.matchGlobParts(patternParts, pathParts)
}

func (gi *GitIgnore) matchGlobParts(patternParts, pathParts []string) bool {
	if len(patternParts) == 0 {
		return len(pathParts) == 0
	}
	if len(pathParts) == 0 {
		return false
	}

	pattern := patternParts[0]
	path := pathParts[0]

	if pattern == "*" {
		return gi.matchGlobParts(patternParts[1:], pathParts[1:])
	}

	if pattern == "**" {
		// ** matches zero or more directories
		if len(patternParts) == 1 {
			return true
		}
		for i := 0; i <= len(pathParts); i++ {
			if gi.matchGlobParts(patternParts[1:], pathParts[i:]) {
				return true
			}
		}
		return false
	}

	if gi.matchSimpleGlob(pattern, path) {
		return gi.matchGlobParts(patternParts[1:], pathParts[1:])
	}

	return false
}

func (gi *GitIgnore) matchSimpleGlob(pattern, s string) bool {
	if pattern == "*" {
		return true
	}

	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		middle := pattern[1 : len(pattern)-1]
		return strings.Contains(s, middle)
	}

	if strings.HasPrefix(pattern, "*") {
		suffix := pattern[1:]
		return strings.HasSuffix(s, suffix)
	}

	if strings.HasSuffix(pattern, "*") {
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(s, prefix)
	}

	return pattern == s
}
