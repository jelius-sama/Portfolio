package handler

import (
	"encoding/json"
	"net/http"
	"os"
)

func GetVersion(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("version")
	if version == "" {
		http.Error(w, "version not set", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"version": version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
