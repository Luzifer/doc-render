// Package api contains the API definition for the project
package api

import (
	"encoding/json"
	"net/http"

	"github.com/Luzifer/doc-render/pkg/persist"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	// Option confiures options for the API server
	Option func(*Server)

	// Server represents the API server holding the methods for the routes
	Server struct {
		persistBackend persist.Backend
		sourceSetDir   string
		texAPIJobURL   string
	}

	renderRequest struct {
		FoxCSV *string        `json:"foxCSV,omitempty"`
		Values map[string]any `json:"values"`
	}
)

// New creates a new Server
func New(opts ...Option) *Server {
	s := &Server{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithPersistBackend configures a backend to persist templates in
func WithPersistBackend(backend persist.Backend) Option {
	return func(s *Server) { s.persistBackend = backend }
}

// WithSourceSetDir configures the base-path of the source-set directory
func WithSourceSetDir(dir string) Option {
	return func(s *Server) { s.sourceSetDir = dir }
}

// WithTexAPIJobURL configures the URL of the TeX-API `/job` endpoint
func WithTexAPIJobURL(url string) Option {
	return func(s *Server) { s.texAPIJobURL = url }
}

// Register adds the routes to the router using a sub-router on the `/api` prefix
func (s Server) Register(r *mux.Router) {
	sr := r.PathPrefix("/api").Subrouter()
	sr.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })

	sr.HandleFunc("/config", s.handleConfigRoute).Methods(http.MethodGet)

	sr.HandleFunc("/persist", s.handlePersistCreate).Methods(http.MethodPost)
	sr.HandleFunc("/persist/{uid}", s.handlePersistGet).Methods(http.MethodGet)

	sr.HandleFunc("/render/{sourceset}", s.handleRenderRoute).Methods(http.MethodPost)

	sr.HandleFunc("/sets", s.handleSourceSetRoute).Methods(http.MethodGet)
}

func (Server) respondJSON(w http.ResponseWriter, status int, err error, data any) {
	var (
		reqID  = uuid.New().String()
		logger = logrus.WithField("req_id", reqID)
	)

	if err != nil {
		logger.WithError(err).Error("handling http request")
		data = map[string]any{
			"success":   false,
			"requestId": reqID,
		}

		if status == http.StatusOK {
			status = http.StatusInternalServerError
		}
	}

	if data == nil && status == http.StatusOK {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(data); err != nil {
		logger.WithError(err).Error("encoding response")
	}
}
