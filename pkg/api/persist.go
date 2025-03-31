package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	persistCreateResponse struct {
		UID string `json:"uid"`
	}
)

func (s Server) handlePersistCreate(w http.ResponseWriter, r *http.Request) {
	templateJSON, err := io.ReadAll(r.Body)
	if err != nil {
		s.respondJSON(w, http.StatusBadRequest, fmt.Errorf("reading body: %w", err), nil)
		return
	}

	uid, err := s.persistBackend.Store(templateJSON)
	if err != nil {
		s.respondJSON(w, http.StatusInternalServerError, fmt.Errorf("storing template: %w", err), nil)
		return
	}

	s.respondJSON(w, http.StatusCreated, nil, persistCreateResponse{
		UID: uid,
	})
}

func (s Server) handlePersistGet(w http.ResponseWriter, r *http.Request) {
	templateJSON, err := s.persistBackend.Get(mux.Vars(r)["uid"])
	if err != nil {
		s.respondJSON(w, http.StatusInternalServerError, fmt.Errorf("fetching template: %w", err), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = io.Copy(w, bytes.NewReader(templateJSON)); err != nil {
		s.respondJSON(w, http.StatusInternalServerError, fmt.Errorf("writing template: %w", err), nil)
		return
	}
}
