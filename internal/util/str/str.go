package str

import (
	"strings"
	"unicode"
)

// IsEmpty checks if string is empty
func IsEmpty(s string) bool {
	return s == ""
}

// IsBlank checks if string is blank (empty or whitespace only)
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty checks if string is not empty
func IsNotEmpty(s string) bool {
	return s != ""
}

// IsNotBlank checks if string is not blank
func IsNotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

// Trim trims whitespace
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// TrimLeft trims left whitespace
func TrimLeft(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

// TrimRight trims right whitespace
func TrimRight(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// TrimPrefix trims a prefix
func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// TrimSuffix trims a suffix
func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// HasPrefix checks if string has prefix
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix checks if string has suffix
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Contains checks if string contains substring
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ContainsAny checks if string contains any characters
func ContainsAny(s, chars string) bool {
	return strings.ContainsAny(s, chars)
}

// ContainsRune checks if string contains rune
func ContainsRune(s string, r rune) bool {
	return strings.ContainsRune(s, r)
}

// Index finds the index of substring
func Index(s, substr string) int {
	return strings.Index(s, substr)
}

// LastIndex finds the last index of substring
func LastIndex(s, substr string) int {
	return strings.LastIndex(s, substr)
}

// Count counts non-overlapping occurrences
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// ToUpper converts to uppercase
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower converts to lowercase
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToTitle converts to title case
func ToTitle(s string) string {
	return strings.ToTitle(s)
}

// Capitalize capitalizes the first character
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Uncapitalize uncapitalizes the first character
func Uncapitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Split splits by whitespace
func Split(s string) []string {
	return strings.Fields(s)
}

// SplitBy splits by separator
func SplitBy(s, sep string) []string {
	return strings.Split(s, sep)
}

// SplitN splits by separator with limit
func SplitN(s, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

// Join joins strings with separator
func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// Replace replaces occurrences
func Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

// ReplaceAll replaces all occurrences
func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// Repeat repeats string n times
func Repeat(s string, n int) string {
	return strings.Repeat(s, n)
}

// Reverse reverses a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// LeftPad pads string on the left
func LeftPad(s string, pad string, length int) string {
	if len(s) >= length {
		return s
	}
	padLen := length - len(s)
	padStr := Repeat(pad, (padLen/len(pad))+1)
	return padStr[:padLen] + s
}

// RightPad pads string on the right
func RightPad(s string, pad string, length int) string {
	if len(s) >= length {
		return s
	}
	padLen := length - len(s)
	padStr := Repeat(pad, (padLen/len(pad))+1)
	return s + padStr[:padLen]
}

// Truncate truncates string to length
func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}

// TruncateWithSuffix truncates string with suffix
func TruncateWithSuffix(s string, length int, suffix string) string {
	if len(s) <= length {
		return s
	}
	return s[:length-len(suffix)] + suffix
}

// Ellipsis truncates with ellipsis
func Ellipsis(s string, length int) string {
	return TruncateWithSuffix(s, length, "...")
}

// StartsWith checks if starts with (alias for HasPrefix)
func StartsWith(s, prefix string) bool {
	return HasPrefix(s, prefix)
}

// EndsWith checks if ends with (alias for HasSuffix)
func EndsWith(s, suffix string) bool {
	return HasSuffix(s, suffix)
}

// Equals checks equality
func Equals(a, b string) bool {
	return a == b
}

// EqualsIgnoreCase checks equality ignoring case
func EqualsIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}

// IsAlpha checks if string is alphabetic
func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsNumeric checks if string is numeric
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsAlphanumeric checks if string is alphanumeric
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsWhitespace checks if string is whitespace only
func IsWhitespace(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return len(s) > 0
}

// DefaultIfEmpty returns default if empty
func DefaultIfEmpty(s, defaultValue string) string {
	if s == "" {
		return defaultValue
	}
	return s
}

// DefaultIfBlank returns default if blank
func DefaultIfBlank(s, defaultValue string) string {
	if IsBlank(s) {
		return defaultValue
	}
	return s
}

// Substring extracts substring
func Substring(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start >= end {
		return ""
	}
	return s[start:end]
}

// Lines splits into lines
func Lines(s string) []string {
	return strings.Split(s, "\n")
}

// Words splits into words
func Words(s string) []string {
	return strings.Fields(s)
}

// Chars splits into characters
func Chars(s string) []string {
	result := make([]string, len(s))
	for i, r := range s {
		result[i] = string(r)
	}
	return result
}

// Length returns length
func Length(s string) int {
	return len(s)
}

// RuneCount returns rune count
func RuneCount(s string) int {
	return len([]rune(s))
}
