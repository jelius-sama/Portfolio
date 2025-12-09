package handler

import (
    "github.com/jelius-sama/logger"
    "net/http"
    "os/exec"
    "syscall"
)

type logWriter struct {
    fn func(...any)
}

func (lw logWriter) Write(p []byte) (int, error) {
    lw.fn(string(p))
    return len(p), nil
}

// TODO: Send Email to admin for approval
func UpdateServer(w http.ResponseWriter, r *http.Request) {
    if !VerifySudo(w, r) {
        return
    }

    // Execute update script (in background)
    cmd := exec.Command("update-prod")
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Setsid: true,
    }

    cmd.Stdout = logWriter{fn: logger.Info}
    cmd.Stderr = logWriter{fn: logger.Error}

    if err := cmd.Start(); err != nil {
        logger.Error("Failed to start update:", err)
        http.Error(w, "Failed to start update", http.StatusInternalServerError)
        return
    }
    logger.Info("Update script started, the server will be killed shortly...")

    // INFO: If the update succeeds, the server will restart and cannot notify the user afterward.
    //       Therefore we send the response beforehand, even though the update might fail.
    w.WriteHeader(http.StatusNoContent)
}
