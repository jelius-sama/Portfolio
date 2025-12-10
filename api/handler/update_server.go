package handler

import (
    "bytes"
    "encoding/json"
    "fmt"
    mailer "github.com/jelius-sama/libmailer/api"
    "github.com/jelius-sama/logger"
    "io"
    "net/http"
    "net/url"
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
        tok, err := GenerateSecureToken(64)
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

        subject := "Server Update Confirmation"
        u, _ := url.Parse(os.Getenv("host"))
        hostname := u.Hostname()
        if hostname == "" {
            hostname = "Portfolio"
        }
        body := fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Server Update Confirmation</title>
<style>
  :root{
    --bg:#ffffff;
    --card:#f6f7fb;
    --text:#0f1724;
    --muted:#6b7280;
    --accent:#0b69ff;
    --accent-contrast:#ffffff;
    --code-bg:#f3f4f6;
  }
  @media (prefers-color-scheme: dark){
    :root{
      --bg:#0b0f14;
      --card:#0f1720;
      --text:#e6eef8;
      --muted:#9aa6b2;
      --accent:#58a6ff;
      --accent-contrast:#031022;
      --code-bg:#071018;
    }
  }

  html,body{height:100%%;margin:0;background:var(--bg);font-family:Inter,-apple-system,system-ui,"Segoe UI",Roboto,"Helvetica Neue",Arial;}
  .wrap{min-height:100%%;display:flex;align-items:center;justify-content:center;padding:28px;}
  .card{width:100%%;max-width:680px;background:linear-gradient(180deg, rgba(255,255,255,0.02), transparent), var(--card);border-radius:12px;padding:28px;box-shadow:0 6px 30px rgba(2,6,23,0.12);border:1px solid rgba(0,0,0,0.04);}
  .logo{font-weight:700;color:var(--text);letter-spacing:0.4px;margin-bottom:12px}
  h1{margin:0 0 8px 0;color:var(--text);font-size:20px}
  p.lead{margin:0 0 20px 0;color:var(--muted);line-height:1.45}
  .cta{display:inline-block;padding:12px 18px;border-radius:8px;background:var(--accent);color:var(--accent-contrast);text-decoration:none;font-weight:600;margin:12px 0;box-shadow:0 6px 18px rgba(11,105,255,0.12);}
  .token{display:inline-block;padding:10px 12px;background:var(--code-bg);border-radius:8px;color:var(--text);font-family:ui-monospace,SFMono-Regular,Menlo,Monaco,monospace;word-break:break-all;margin-bottom:12px}
  footer{margin-top:18px;color:var(--muted);font-size:13px}
</style>
</head>
<body>
  <div class="wrap">
    <div class="card">
      <div class="logo">%s</div>
      <h1>Server update requested</h1>
      <p class="lead">
        An update of your portfolio server was requested. If you trust this request, approve it using the button below.
      </p>

      <a class="cta" href="%s" target="_blank" rel="noopener noreferrer">Approve update</a>

      <footer>If you did not request this, no action is required — but please investigate any unexpected activity.</footer>
    </div>
  </div>
</body>
</html>`, hostname, os.Getenv("host")+"/api/specialized_task?task=update_server&token="+tok)

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
