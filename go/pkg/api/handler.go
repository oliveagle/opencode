package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anomalyco/opencode/pkg/core"
	"github.com/anomalyco/opencode/pkg/types"
)

type Handler struct {
	files    *core.FileService
	projects *core.ProjectService
}

func NewHandler(files *core.FileService, projects *core.ProjectService) *Handler {
	return &Handler{
		files:    files,
		projects: projects,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// File operations
	mux.HandleFunc("/api/files/read", h.readFile)
	mux.HandleFunc("/api/files/write", h.writeFile)
	mux.HandleFunc("/api/files/list", h.listDir)
	mux.HandleFunc("/api/files/info", h.getFileInfo)
	mux.HandleFunc("/api/files/delete", h.deleteFile)
	mux.HandleFunc("/api/files/create-dir", h.createDir)

	// Project operations
	mux.HandleFunc("/api/projects/list", h.listProjects)
	mux.HandleFunc("/api/projects/get", h.getProject)
	mux.HandleFunc("/api/projects/create", h.createProject)
	mux.HandleFunc("/api/projects/delete", h.deleteProject)

	// Health check
	mux.HandleFunc("/health", h.health)
}

func (h *Handler) readFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		h.respond(w, http.StatusBadRequest, nil, "path required")
		return
	}

	content, err := h.files.ReadFile(path)
	if err != nil {
		h.respond(w, http.StatusNotFound, nil, err.Error())
		return
	}

	h.respond(w, http.StatusOK, map[string]string{"content": content}, "")
}

func (h *Handler) writeFile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := h.files.WriteFile(req.Path, req.Content); err != nil {
		h.respond(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	h.respond(w, http.StatusOK, nil, "")
}

func (h *Handler) listDir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}

	files, err := h.files.ListDir(path)
	if err != nil {
		h.respond(w, http.StatusNotFound, nil, err.Error())
		return
	}

	h.respond(w, http.StatusOK, map[string]interface{}{"files": files}, "")
}

func (h *Handler) getFileInfo(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		h.respond(w, http.StatusBadRequest, nil, "path required")
		return
	}

	info, err := h.files.GetFileInfo(path)
	if err != nil {
		h.respond(w, http.StatusNotFound, nil, err.Error())
		return
	}

	h.respond(w, http.StatusOK, info, "")
}

func (h *Handler) deleteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		h.respond(w, http.StatusBadRequest, nil, "path required")
		return
	}

	if err := h.files.DeleteFile(path); err != nil {
		h.respond(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	h.respond(w, http.StatusOK, nil, "")
}

func (h *Handler) createDir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		h.respond(w, http.StatusBadRequest, nil, "path required")
		return
	}

	if err := h.files.CreateDir(path); err != nil {
		h.respond(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	h.respond(w, http.StatusOK, nil, "")
}

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	projects := h.projects.ListProjects()
	h.respond(w, http.StatusOK, map[string]interface{}{"projects": projects}, "")
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		h.respond(w, http.StatusBadRequest, nil, "name required")
		return
	}

	proj := h.projects.GetProject(name)
	if proj == nil {
		h.respond(w, http.StatusNotFound, nil, "project not found")
		return
	}

	h.respond(w, http.StatusOK, proj, "")
}

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var proj types.ProjectConfig
	if err := json.NewDecoder(r.Body).Decode(&proj); err != nil {
		h.respond(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	cfg := &types.ProjectConfig{
		Name:  proj.Name,
		Path:  proj.Path,
		Root:  proj.Root,
		Alias: proj.Alias,
	}

	if err := h.projects.CreateProject(cfg); err != nil {
		h.respond(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	h.respond(w, http.StatusOK, cfg, "")
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		h.respond(w, http.StatusBadRequest, nil, "name required")
		return
	}

	if err := h.projects.DeleteProject(name); err != nil {
		h.respond(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	h.respond(w, http.StatusOK, nil, "")
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	h.respond(w, http.StatusOK, map[string]string{"status": "healthy"}, "")
}

func (h *Handler) respond(w http.ResponseWriter, code int, data interface{}, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := types.ResponsePayload{
		Success: code >= 200 && code < 300,
		Data:    data,
		Error:   errMsg,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
