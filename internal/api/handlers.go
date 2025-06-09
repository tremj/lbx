package api

import (
	"io"
	"net/http"

	"github.com/tremj/lbx/internal/core"
)

func ListConfigsHandler(w http.ResponseWriter, r *http.Request) {}

func SaveConfigHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing config name", http.StatusBadGateway)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = core.SaveConfig(ctx, name, data)
	if err != nil {
		http.Error(w, "save failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteConfigHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing config name", http.StatusBadGateway)
		return
	}

	err := core.DeleteConfig(ctx, name)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
