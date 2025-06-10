package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tremj/lbx/internal/core"
	"github.com/tremj/lbx/internal/storage"
)

func ListConfigsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := storage.GetKeys(ctx)
	if err != nil {
		http.Error(w, "error fetching all keys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := chi.URLParam(r, "name")
	data, err := storage.Get(ctx, name)
	if err != nil {
		http.Error(w, "Config not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func SaveConfigHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing config name", http.StatusBadGateway)
		return
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
