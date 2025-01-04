// Package frontend contains the frontend assets and a server to
// register routes
package frontend

import (
	"bytes"
	"embed"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	// Server represents the fRontend server holding the methods for the routes
	Server struct{}
)

//go:embed assets/**
var assets embed.FS

// New creates a new Server
func New() *Server {
	return &Server{}
}

// Register adds the routes to the router
func (s Server) Register(r *mux.Router) {
	r.HandleFunc("/", s.handleIndexRoute).Methods(http.MethodGet)
	r.PathPrefix("/assets").Handler(http.FileServer(http.FS(assets)))
}

func (Server) handleIndexRoute(w http.ResponseWriter, _ *http.Request) {
	index, err := assets.ReadFile("assets/index.html")
	if err != nil {
		http.Error(w, "index not found", http.StatusNotFound)
		return
	}

	if _, err = io.Copy(w, bytes.NewReader(index)); err != nil {
		logrus.WithError(err).Error("copying index source to browser")
	}
}
