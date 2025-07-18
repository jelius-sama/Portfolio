package handler

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

func VerifySudo(w http.ResponseWriter, r *http.Request) {
	// Extract Bearer token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	password := strings.TrimPrefix(authHeader, "Bearer ")

	// Load sudo password from systemd credential file
	credPath := "/run/cred/Portfolio.service/sudo_pass"
	expectedPasswordBytes, err := os.ReadFile(credPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	expectedPassword := strings.TrimSpace(string(expectedPasswordBytes))

	// Compare passwords
	if subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
