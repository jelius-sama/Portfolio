package handler

import (
	"crypto/subtle"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

// TODO: Send Email to admin for approval
func UpdateServer(w http.ResponseWriter, r *http.Request) {
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

	// Execute update script (in background)
	cmd := exec.Command("bash", path.Join(os.Getenv("home"), "/update_prod.sh"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		http.Error(w, "Failed to start update", http.StatusInternalServerError)
		return
	}

	go func() {
		_ = cmd.Wait()
	}()

	w.WriteHeader(http.StatusNoContent)
}
