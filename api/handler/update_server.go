package handler

import (
	"KazuFolio/util"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// TODO: Send Email to admin for approval
func UpdateServer(w http.ResponseWriter, r *http.Request) {
	util.VerifySudo(w, r)

	// Execute update script (in background)
	cmd := exec.Command("bash", filepath.Join(os.Getenv("home"), "/update_prod.sh"))
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
