package tui

import (
	"regexp"
	"strings"
)

type SearchMatch struct {
	Row    int
	Col    int
	Text   string
	Length int
}

type SearchReplace struct {
	pattern    string
	isRegex    bool
	caseSensitive bool
	regex      *regexp.Regexp
}

func NewSearchReplace() *SearchReplace {
	return &SearchReplace{
		pattern:       "",
		isRegex:       false,
		caseSensitive: false,
		regex:         nil,
	}
}

func (sr *SearchReplace) SetPattern(pattern string, isRegex bool) error {
	sr.pattern = pattern

	if isRegex {
		flags := ""
		if !sr.caseSensitive {
			flags = "(?i)"
		}
		re, err := regexp.Compile(flags + pattern)
		if err != nil {
			return err
		}
		sr.regex = re
		sr.isRegex = true
	} else {
		sr.isRegex = false
		sr.regex = nil
	}

	return nil
}

func (sr *SearchReplace) SetCaseSensitive(sensitive bool) {
	sr.caseSensitive = sensitive
	if sr.isRegex && sr.regex != nil {
		// Recompile regex with new case sensitivity
		sr.SetPattern(sr.pattern, sr.isRegex)
	}
}

func (sr *SearchReplace) FindAll(text string) []SearchMatch {
	var matches []SearchMatch

	if sr.isRegex && sr.regex != nil {
		for _, match := range sr.regex.FindAllStringIndex(text, -1) {
			matchText := text[match[0]:match[1]]
			row := strings.Count(text[:match[0]], "\n")
			colStart := match[0]
			if idx := strings.LastIndex(text[:match[0]], "\n"); idx != -1 {
				colStart = match[0] - idx - 1
			}

			matches = append(matches, SearchMatch{
				Row:    row,
				Col:    colStart,
				Text:   matchText,
				Length: match[1] - match[0],
			})
		}
	} else {
		// Plain text search
		searchText := sr.pattern
		if !sr.caseSensitive {
			searchText = strings.ToLower(sr.pattern)
			text = strings.ToLower(text)
		}

		idx := 0
		for {
			pos := strings.Index(text[idx:], searchText)
			if pos == -1 {
				break
			}

			absPos := idx + pos
			row := strings.Count(text[:absPos], "\n")
			colStart := absPos
			if lastNewline := strings.LastIndex(text[:absPos], "\n"); lastNewline != -1 {
				colStart = absPos - lastNewline - 1
			}

			matches = append(matches, SearchMatch{
				Row:    row,
				Col:    colStart,
				Text:   sr.pattern,
				Length: len(sr.pattern),
			})

			idx = absPos + len(searchText)
		}
	}

	return matches
}

func (sr *SearchReplace) Replace(text, replacement string) string {
	if sr.isRegex && sr.regex != nil {
		return sr.regex.ReplaceAllString(text, replacement)
	}

	searchText := sr.pattern
	if !sr.caseSensitive {
		return replaceInsensitive(text, sr.pattern, replacement)
	}

	return strings.ReplaceAll(text, searchText, replacement)
}

func (sr *SearchReplace) ReplaceAll(text, replacement string) string {
	return sr.Replace(text, replacement)
}

func (sr *SearchReplace) GetPattern() string {
	return sr.pattern
}

func (sr *SearchReplace) IsRegex() bool {
	return sr.isRegex
}

func replaceInsensitive(text, search, replacement string) string {
	lower := strings.ToLower(text)
	searchLower := strings.ToLower(search)

	var result strings.Builder
	idx := 0

	for {
		pos := strings.Index(lower[idx:], searchLower)
		if pos == -1 {
			result.WriteString(text[idx:])
			break
		}

		absPos := idx + pos
		result.WriteString(text[idx:absPos])
		result.WriteString(replacement)
		idx = absPos + len(search)
	}

	return result.String()
}
