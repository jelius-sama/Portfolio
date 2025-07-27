package handler

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

func VerifySudo(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	password := strings.TrimPrefix(authHeader, "Bearer ")

	expectedPassword := os.Getenv("SUDO_PASS")
	if expectedPassword == "" {
		http.Error(w, "Server misconfigured", http.StatusInternalServerError)
		return false
	}

	if subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}

func VerifySudoHandler(w http.ResponseWriter, r *http.Request) {
	if !VerifySudo(w, r) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
