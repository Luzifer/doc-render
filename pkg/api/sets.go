package api

import (
	"fmt"
	"net/http"

	"github.com/Luzifer/doc-render/pkg/latex"
)

func (s Server) handleSourceSetRoute(w http.ResponseWriter, _ *http.Request) {
	sets, err := latex.GetSourceSets(s.sourceSetDir)
	if err != nil {
		s.respondJSON(w, http.StatusInternalServerError, fmt.Errorf("getting source sets: %w", err), nil)
		return
	}

	s.respondJSON(w, http.StatusOK, nil, sets)
}
