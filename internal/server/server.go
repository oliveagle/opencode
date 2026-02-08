package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/anomalyco/opencode/internal/util/log"
)

// Server represents the HTTP server
type Server struct {
	port    int
	server  *http.Server
	router  *Router
	clients map[*Client]bool
	mu      sync.RWMutex
	events  chan Event
}

// Router handles HTTP routing
type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

// Handle registers a handler for a method and path
func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]http.HandlerFunc)
	}
	r.routes[path][method] = handler
}

// Get registers a GET handler
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.Handle("GET", path, handler)
}

// Post registers a POST handler
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.Handle("POST", path, handler)
}

// Put registers a PUT handler
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.Handle("PUT", path, handler)
}

// Delete registers a DELETE handler
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.Handle("DELETE", path, handler)
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	if handlers, ok := r.routes[path]; ok {
		if handler, ok := handlers[method]; ok {
			handler(w, req)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Try to match with trailing slash
	if handlers, ok := r.routes[path+"/"]; ok {
		if handler, ok := handlers[method]; ok {
			handler(w, req)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

// New creates a new server
func New(port int) *Server {
	s := &Server{
		port:    port,
		router:  NewRouter(),
		clients: make(map[*Client]bool),
		events:  make(chan Event, 1000),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// API routes
	s.router.Get("/api/session", s.handleSessionList)
	s.router.Post("/api/session", s.handleSessionCreate)
	s.router.Get("/api/session/:id", s.handleSessionGet)
	s.router.Delete("/api/session/:id", s.handleSessionDelete)

	s.router.Get("/api/config", s.handleConfigGet)
	s.router.Post("/api/config", s.handleConfigUpdate)

	s.router.Get("/api/models", s.handleModelsList)
	s.router.Get("/api/providers", s.handleProvidersList)

	s.router.Get("/api/file", s.handleFileList)
	s.router.Get("/api/file/*", s.handleFileRead)
	s.router.Post("/api/file/*", s.handleFileWrite)

	s.router.Post("/api/chat", s.handleChat)
	s.router.Post("/api/stream", s.handleStream)

	// SSE endpoint
	s.router.Get("/events", s.handleSSE)
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	// Start event broadcaster
	go s.broadcastEvents(ctx)

	log.Default.Info("Server starting", map[string]interface{}{
		"port": s.port,
	})

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

// Emit emits an event to all connected clients
func (s *Server) Emit(event Event) {
	select {
	case s.events <- event:
	default:
		log.Default.Warn("Event channel full, dropping event", nil)
	}
}

func (s *Server) broadcastEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.events:
			s.mu.RLock()
			for client := range s.clients {
				select {
				case client.events <- event:
				default:
					// Client buffer full, skip
				}
			}
			s.mu.RUnlock()
		}
	}
}

// Event represents a server-sent event
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Client represents a connected SSE client
type Client struct {
	events chan Event
}

func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Register client
	client := &Client{events: make(chan Event, 100)}
	s.mu.Lock()
	s.clients[client] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, client)
		s.mu.Unlock()
	}()

	// Send events
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case event := <-client.events:
			data, _ := json.Marshal(event)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

// Placeholder handlers - to be implemented
func (s *Server) handleSessionList(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	json.NewEncoder(w).Encode([]interface{}{})
}

func (s *Server) handleSessionCreate(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleSessionGet(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSessionDelete(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleConfigGet(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func (s *Server) handleConfigUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleModelsList(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	json.NewEncoder(w).Encode([]interface{}{})
}

func (s *Server) handleProvidersList(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	json.NewEncoder(w).Encode([]interface{}{})
}

func (s *Server) handleFileList(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	json.NewEncoder(w).Encode([]interface{}{})
}

func (s *Server) handleFileRead(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleFileWrite(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement streaming response
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send a placeholder response
	fmt.Fprintf(w, "data: {\"type\":\"done\"}\n\n")
	flusher.Flush()
}

// Port returns the server port
func (s *Server) Port() int {
	return s.port
}

// Router returns the router for custom route registration
func (s *Server) Router() *Router {
	return s.router
}

// WaitForReady waits for the server to be ready
func (s *Server) WaitForReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", s.port))
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server not ready after %v", timeout)
}
