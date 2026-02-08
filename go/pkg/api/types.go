package api

// FileInfo represents metadata about a file
type FileInfo struct {
	Path    string
	Name    string
	Size    int64
	IsDir   bool
	ModTime int64
}

// ProjectConfig represents project configuration
type ProjectConfig struct {
	Name  string
	Path  string
	Root  string
	Alias string
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
