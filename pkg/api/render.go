package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Luzifer/doc-render/pkg/latex"
	"github.com/Luzifer/doc-render/pkg/recipientcsv"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

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
