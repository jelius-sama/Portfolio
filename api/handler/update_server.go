package handler

import (
	"crypto/subtle"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func UpdateServer(w http.ResponseWriter, r *http.Request) {
	// Extract Bearer token
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	password := strings.TrimPrefix(authHeader, "Bearer ")

	// Load expected password from systemd credentials
	credPath := "/run/cred/Portfolio.service/sudo_pass"
	expectedPasswordBytes, err := os.ReadFile(credPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	expectedPassword := strings.TrimSpace(string(expectedPasswordBytes))

	// Secure compare
	if subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) != 1 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Execute update script
	cmd := exec.Command("bash", "-c", "~/update_prod.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		http.Error(w, "Failed to start update", http.StatusInternalServerError)
		return
	}

	// Detach the update process (don't wait for completion)
	go func() {
		_ = cmd.Wait()
	}()

	w.WriteHeader(http.StatusNoContent)
}
