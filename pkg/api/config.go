package api

import "net/http"

type (
	config struct {
		HasPersist bool `json:"hasPersist"`
	}
)

func (s Server) handleConfigRoute(w http.ResponseWriter, _ *http.Request) {
	s.respondJSON(w, http.StatusOK, nil, config{
		HasPersist: s.persistBackend != nil,
	})
}
