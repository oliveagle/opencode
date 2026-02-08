package types

// FileInfo represents metadata about a file
type FileInfo struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime int64  `json:"mod_time"`
}

// ProjectConfig represents project configuration
type ProjectConfig struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Root  string `json:"root"`
	Alias string `json:"alias"`
}

// RequestPayload represents a request from TUI to backend
type RequestPayload struct {
	Type    string            `json:"type"`
	Path    string            `json:"path"`
	Content string            `json:"content,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
}

// ResponsePayload represents a response from backend to TUI
type ResponsePayload struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SearchMatch represents a search match in file content
type SearchMatch struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Content string `json:"content"`
}
