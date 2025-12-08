package handler

import (
    "github.com/jelius-sama/logger"
    "net/http"
    "os/exec"
    "strings"
)

type logWriter struct {
    fn func(...any)
}

func (lw logWriter) Write(p []byte) (int, error) {
    lw.fn(string(p))
    return len(p), nil
}

func fuzzySignalTerminated(s string) bool {
    s = strings.ToLower(s)
    s = strings.ReplaceAll(s, "\t", " ")
    s = strings.ReplaceAll(s, "\n", " ")
    for strings.Contains(s, "  ") {
        s = strings.ReplaceAll(s, "  ", " ")
    }
    return strings.Contains(s, "signal") && strings.Contains(s, "terminated")
}

// TODO: Send Email to admin for approval
func UpdateServer(w http.ResponseWriter, r *http.Request) {
    if !VerifySudo(w, r) {
        return
    }

    // Execute update script (in background)
    cmd := exec.Command("update-prod")

    cmd.Stdout = logWriter{fn: logger.Info}
    cmd.Stderr = logWriter{fn: logger.Error}

    if err := cmd.Start(); err != nil {
        logger.Error("Failed to start update:", err)
        http.Error(w, "Failed to start update", http.StatusInternalServerError)
        return
    }
    logger.Info("Update script started, the server will be killed shortly...")

    go func() {
        // If Wait returns, it means update script exited instead of killing us.
        if err := cmd.Wait(); err != nil {
            if !fuzzySignalTerminated(err.Error()) {
                logger.Error("Update script ended unexpectedly:", err)
            }
        } else {
            logger.Error("Update script ended unexpectedly without killing process")
        }
    }()

    // INFO: If the update succeeds, the server will restart and cannot notify the user afterward.
    //       Therefore we send the response beforehand, even though the update might fail.
    w.WriteHeader(http.StatusNoContent)
}
