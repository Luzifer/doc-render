// Package api contains the API definition for the project
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Luzifer/doc-render/pkg/latex"
	"github.com/Luzifer/doc-render/pkg/recipientcsv"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	// Option confiures options for the API server
	Option func(*Server)

	// Server represents the API server holding the methods for the routes
	Server struct {
		sourceSetDir string
		texAPIJobURL string
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
	sr.HandleFunc("/render/{sourceset}", s.handleRenderRoute).Methods(http.MethodPost)
	sr.HandleFunc("/sets", s.handleSourceSetRoute).Methods(http.MethodGet)
}

func (s Server) handleRenderRoute(w http.ResponseWriter, r *http.Request) {
	var (
		addrTo    = []recipientcsv.Person{{}}
		err       error
		payload   renderRequest
		sourceSet = mux.Vars(r)["sourceset"]
	)

	if ct, _, _ := strings.Cut(r.Header.Get("Content-Type"), ";"); ct != "application/json" {
		s.respondJSON(w, http.StatusBadRequest, fmt.Errorf("invalid payload type %q", ct), nil)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.respondJSON(w, http.StatusBadRequest, fmt.Errorf("parsing request payload: %w", err), nil)
		return
	}

	if payload.FoxCSV != nil {
		if addrTo, err = recipientcsv.Parse(strings.NewReader(*payload.FoxCSV)); err != nil {
			s.respondJSON(w, http.StatusBadRequest, fmt.Errorf("parsing FoxCSV: %w", err), nil)
			return
		}
	}

	// Generate document
	pdf, err := latex.Render(r.Context(), latex.RenderOpts{
		TexAPIURL: s.texAPIJobURL,

		SourceBaseFolder: s.sourceSetDir,
		SourceSet:        sourceSet,

		Recipients: addrTo,
		Values:     payload.Values,
	})
	if err != nil {
		s.respondJSON(w, http.StatusInternalServerError, fmt.Errorf("rendering PDF: %w", err), nil)
		return
	}
	defer func() {
		if err := pdf.Close(); err != nil {
			logrus.WithError(err).Error("closing PDF reader")
		}
	}()

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Cache-Control", "no-cache")

	if _, err = io.Copy(w, pdf); err != nil {
		logrus.WithError(err).Error("copying PDF to remote browser")
	}
}

func (s Server) handleSourceSetRoute(w http.ResponseWriter, _ *http.Request) {
	sets, err := latex.GetSourceSets(s.sourceSetDir)
	if err != nil {
		s.respondJSON(w, http.StatusInternalServerError, fmt.Errorf("getting source sets: %w", err), nil)
		return
	}

	s.respondJSON(w, http.StatusOK, nil, sets)
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

	if err = json.NewEncoder(w).Encode(data); err != nil {
		logger.WithError(err).Error("encoding response")
	}
}
