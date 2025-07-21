package util

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// AddrOf takes a value and returns its address as a pointer.
func AddrOf[T any](literal T) *T { return &literal }

func IsValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false // not a number
	}
	return port >= 1 && port <= 65535
}

func VerifySudo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	password := strings.TrimPrefix(authHeader, "Bearer ")

	expectedPassword := os.Getenv("SUDO_PASS")
	if expectedPassword == "" {
		http.Error(w, "Server misconfigured", http.StatusInternalServerError)
		return
	}

	if subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
