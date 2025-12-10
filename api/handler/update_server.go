package handler

import (
    "bytes"
    "encoding/json"
    mailer "github.com/jelius-sama/libmailer/api"
    "github.com/jelius-sama/logger"
    "io"
    "net/http"
    "os"
    "os/exec"
    "sync"
    "time"
)

type logWriter struct {
    fn func(...any)
}

func (lw logWriter) Write(p []byte) (int, error) {
    lw.fn(string(p))
    return len(p), nil
}

var (
    tokenValue  *string
    tokenExpiry time.Time
    mu          sync.Mutex
)

// init starts a GC that clears the token when expired.
func init() {
    go func() {
        t := time.NewTicker(TokenValidity)
        defer t.Stop()

        for range t.C {
            mu.Lock()
            if tokenValue != nil && time.Now().After(tokenExpiry) {
                tokenValue = nil
            }
            mu.Unlock()
        }
    }()
}

func SetToken(value string, validity time.Duration) {
    mu.Lock()
    defer mu.Unlock()

    tokenValue = &value
    tokenExpiry = time.Now().Add(validity)
}

func GetToken() *string {
    mu.Lock()
    defer mu.Unlock()

    if tokenValue == nil {
        return nil
    }
    if time.Now().After(tokenExpiry) {
        tokenValue = nil
        return nil
    }
    return tokenValue
}

type Payload struct {
    Token string `json:"token"`
}

func UpdateServer(w http.ResponseWriter, r *http.Request) {
    if !VerifySudo(w, r) {
        return
    }

    body, _ := io.ReadAll(r.Body)
    trimmed := bytes.TrimSpace(body)

    var p Payload
    hasToken := false

    if len(trimmed) > 0 {
        // If JSON is malformed → reject request
        if err := json.Unmarshal(trimmed, &p); err != nil {
            logger.Error("Received malformed JSON in server update request.")
            http.Error(w, "You are banned!", http.StatusForbidden)
            return
        }

        hasToken = true
    }

    if hasToken {
        stored := GetToken()
        if stored != nil && p.Token == *stored {
            // Execute update script (in background)
            cmd := exec.Command("sh", "-c", "nohup update-prod >> /var/log/update-prod.log 2>&1 &")
            if err := cmd.Start(); err != nil {
                logger.Error("Failed to start update:", err)
                http.Error(w, "Failed to start update", http.StatusInternalServerError)
                return
            }
            logger.Info("Update script started, the server will be killed shortly...")
        } else {
            logger.Error("Initiated server update with an invalid admin token.")
            http.Error(w, "You are banned!", http.StatusForbidden)
            return
        }
    } else {
        // Generate a secure token for admin confirmation
        tok, err := GenerateSecureToken(16)
        if err != nil {
            logger.Error("Failed to generate admin token:", err)
            http.Error(w, "Failed to start update", http.StatusInternalServerError)
            return
        }

        // Store the token with the configured validity
        SetToken(tok, TokenValidity)

        config, err := mailer.LoadConfig()
        if err != nil {
            logger.Error("Failed to get mailer config:", err)
            http.Error(w, "Failed to start update", http.StatusInternalServerError)
            return
        }

        // Email contents
        subject := "Server Update Confirmation"

        bearer := os.Getenv("SUDO_PASS")

        body := "A server update was requested.\n\n" +
            "To approve it, run the following command on your machine:\n\n" +
            "curl -X POST " + os.Getenv("host") + "/api/update_server \\\n" +
            "     -H \"Authorization: Bearer " + bearer + "\" \\\n" +
            "     -d '{\"token\":\"" + tok + "\"}'\n\n" +
            "If you did not request this, please investigate immediately."

        err = mailer.SendMail(
            config.Host,
            config.Port,
            config.Username,
            config.Password,
            config.From,
            "personal@jelius.dev",
            subject,
            body,
            nil, nil, nil,
        )
        if err != nil {
            logger.Error("Failed to send update confirmation email to admin:", err)
            http.Error(w, "Failed to start update", http.StatusInternalServerError)
            return
        }
    }

    // INFO: If the update succeeds, the server will restart and cannot notify the user afterward.
    //       Therefore we send the response beforehand, even though the update might fail.
    w.WriteHeader(http.StatusNoContent)
}
