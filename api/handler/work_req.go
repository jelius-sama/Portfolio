package handler

import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "html/template"
    "io"
    "net"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
    "time"

    libmailer "github.com/jelius-sama/libmailer/api"
    "github.com/jelius-sama/logger"
)

const (
    maxRequestsPerIdentity = 6
    minIntervalBetweenReqs = 60 * time.Second
    banDuration            = 3 * 365 * 24 * time.Hour // 3 years
    maxFileSize            = 10 << 20                 // 10 MB per file
    maxFormMemory          = 32 << 20                 // 32 MB total in-memory for multipart
)

// identityRecord tracks request history for one hashed identity.
// The identity is SHA-256(email + ":" + canonicalIP) — the raw values
// are never stored. Checking an incoming request hashes the same way
// and compares.
type identityRecord struct {
    count       int       // total accepted requests so far
    lastSeen    time.Time // time of the most recent accepted request
    bannedUntil time.Time // zero value means not banned
}

// rateLimiter holds the in-memory store, safe for concurrent use.
type rateLimiter struct {
    mu      sync.RWMutex
    entries map[[32]byte]*identityRecord // key: SHA-256 digest
}

var limiter = &rateLimiter{
    entries: make(map[[32]byte]*identityRecord),
}

// WorkReq is the struct we populate from the multipart form.
type WorkReq struct {
    Name         string
    Email        string
    Phone        string
    Service      string
    Title        string
    Description  string
    TechPref     string
    TimelineDays string
    DailyRate    string
    Notes        string
    Attachments  []AttachmentMeta
}

// AttachmentMeta holds file metadata only — we do not buffer the full
// file contents in the struct; callers that need the bytes should read
// from the multipart.File directly before this handler returns.
type AttachmentMeta struct {
    Filename string
    Size     int64
    MIMEType string
    TempPath string
}

// emailData is the template data bag passed to emailHTMLTemplate.
type emailData struct {
    Name         string
    Email        string
    Phone        string
    Service      string
    Title        string
    Description  string
    TechPref     string
    TimelineDays string
    DailyRate    string
    Notes        string
    Attachments  []AttachmentMeta
    ReceivedAt   string
}

var unsafeChars = regexp.MustCompile(`[^a-zA-Z0-9._\-]`)

// safeFilename strips characters that are unsafe in Unix file paths so
// the original filename can be used directly on disk without shell-escaping.
func safeFilename(name string) string {
    base := filepath.Base(name) // drop any path component sent by the client
    safe := unsafeChars.ReplaceAllString(base, "_")
    if safe == "" || safe == "." {
        return "_attachment"
    }
    return safe
}

// tempDirForIdentity creates (if necessary) a subdirectory under os.TempDir()
// whose name is the hex-encoded SHA-256 of the identity key — the same digest
// already used for rate-limiting, so no new raw data is introduced.
// The directory name is therefore safe for all Unix filesystems.
//
// The caller is responsible for calling os.RemoveAll on the returned path
// once the email has been sent.
func tempDirForIdentity(key [32]byte) (string, error) {
    dirName := hex.EncodeToString(key[:]) // 64 lowercase hex chars, filesystem-safe
    dir := filepath.Join(os.TempDir(), "workreq_"+dirName)
    if err := os.MkdirAll(dir, 0o700); err != nil {
        return "", fmt.Errorf("create temp dir: %w", err)
    }
    return dir, nil
}

// identityHash derives a 32-byte SHA-256 digest from the normalized
// email and IP. This is the only form in which the identity is stored.
// To check a future request, call the same function and compare digests.
func identityHash(email, ip string) [32]byte {
    normalized := strings.ToLower(strings.TrimSpace(email)) + ":" + ip
    return sha256.Sum256([]byte(normalized))
}

// canonicalIP returns the real client IP, unwrapping common proxy headers.
// It strips the port and normalizes IPv6 so the hash is stable.
func canonicalIP(r *http.Request) string {
    // Trust X-Forwarded-For only for the first (leftmost) hop, which is
    // the actual client when sitting behind a single trusted reverse proxy.
    if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
        if first := strings.SplitN(fwd, ",", 2)[0]; first != "" {
            if ip := net.ParseIP(strings.TrimSpace(first)); ip != nil {
                return ip.String()
            }
        }
    }
    if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
        if ip := net.ParseIP(strings.TrimSpace(realIP)); ip != nil {
            return ip.String()
        }
    }
    host, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        // RemoteAddr had no port (shouldn't happen in practice)
        if ip := net.ParseIP(r.RemoteAddr); ip != nil {
            return ip.String()
        }
        return r.RemoteAddr
    }
    if ip := net.ParseIP(host); ip != nil {
        return ip.String()
    }
    return host
}

type denyReason int

const (
    allowed       denyReason = iota
    deniedBanned             // lifetime ban (exceeded max requests)
    deniedTooFast            // within the 60-second sliding window
)

// check returns whether the identity is allowed to make a new request.
// It does NOT mutate state — call admit() only after you've verified
// the rest of the request is valid.
func (rl *rateLimiter) check(key [32]byte) denyReason {
    rl.mu.RLock()
    rec, exists := rl.entries[key]
    rl.mu.RUnlock()

    if !exists {
        return allowed
    }

    now := time.Now()

    // Banned indefinitely (well, for 3 years from the ban timestamp).
    if !rec.bannedUntil.IsZero() && now.Before(rec.bannedUntil) {
        return deniedBanned
    }

    // Within the 60-second sliding window from the last accepted request.
    if now.Sub(rec.lastSeen) < minIntervalBetweenReqs {
        return deniedTooFast
    }

    return allowed
}

// admit records a successful request for the given key.
// Must be called only after check() returned allowed AND the request
// has been fully validated — do not admit on parse errors.
func (rl *rateLimiter) admit(key [32]byte) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    rec, exists := rl.entries[key]
    if !exists {
        rec = &identityRecord{}
        rl.entries[key] = rec
    }

    rec.count++
    rec.lastSeen = time.Now()

    // If this request pushes them to the max, ban from now.
    if rec.count >= maxRequestsPerIdentity {
        rec.bannedUntil = time.Now().Add(banDuration)
    }
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(payload)
}

func deny(w http.ResponseWriter, reason string) {
    writeJSON(w, http.StatusForbidden, map[string]string{"error": reason})
}

var emailHTMLTemplate = template.Must(template.New("workreq").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1"/>
<meta name="color-scheme" content="dark light"/>
<title>Work Request — {{.Title}}</title>
<style>
  :root{color-scheme:dark light}
  *{box-sizing:border-box;margin:0;padding:0}
  body{background:#000b11;color:#e2e8f0;font-family:'Inter',system-ui,Arial,sans-serif;font-size:15px;line-height:1.6;padding:32px 16px}
  .wrapper{max-width:620px;margin:0 auto}
  .eyebrow{font-family:'JetBrains Mono',monospace;font-size:11px;letter-spacing:.1em;text-transform:uppercase;color:#38bdf8;margin-bottom:8px}
  h1{font-size:22px;font-weight:600;letter-spacing:-.02em;margin-bottom:4px;color:#e2e8f0}
  .subtitle{font-size:13px;color:#7a9ab0;margin-bottom:24px}
  .card{background:#001520;border:1px solid #0f3050;border-radius:10px;overflow:hidden;margin-bottom:16px}
  .card-header{padding:14px 20px;border-bottom:1px solid #0f3050}
  .card-header-title{font-size:11px;font-weight:600;letter-spacing:.08em;text-transform:uppercase;color:#7a9ab0}
  .card-body{padding:20px}
  .row{display:flex;gap:16px;margin-bottom:14px}
  .row:last-child{margin-bottom:0}
  .col{flex:1;min-width:0}
  .label{display:block;font-size:11px;font-weight:600;letter-spacing:.06em;text-transform:uppercase;color:#3d5f75;margin-bottom:3px}
  .value{font-size:14px;color:#e2e8f0;word-break:break-word}
  .value.mono{font-family:'JetBrains Mono',monospace;font-size:13px}
  .tag{display:inline-block;background:#0c2030;color:#38bdf8;font-size:12px;font-weight:500;padding:2px 10px;border-radius:100px;border:1px solid #0f3050}
  .divider{border:none;border-top:1px solid #0f3050;margin:14px 0}
  .prose{font-size:14px;color:#e2e8f0;white-space:pre-wrap;word-break:break-word}
  .att-row{display:flex;align-items:center;gap:12px;background:#0c2030;border:1px solid #0f3050;border-radius:6px;padding:10px 14px;margin-bottom:8px}
  .att-row:last-child{margin-bottom:0}
  .att-icon{font-size:18px;flex-shrink:0}
  .att-name{font-size:13px;font-weight:500;color:#e2e8f0;word-break:break-all}
  .att-meta{font-size:11px;color:#7a9ab0;margin-top:1px}
  .budget-highlight{background:#0c2030;border:1px solid #0f3050;border-radius:8px;padding:14px 18px;display:flex;gap:32px}
  .budget-block .bval{font-family:'JetBrains Mono',monospace;font-size:20px;font-weight:600;color:#38bdf8}
  .budget-block .blabel{font-size:11px;color:#7a9ab0;text-transform:uppercase;letter-spacing:.06em;margin-top:2px}
  .total-note{font-size:12px;color:#3d5f75;margin-top:10px}
  .footer{font-size:11px;color:#3d5f75;text-align:center;margin-top:24px}
  @media(max-width:480px){.row{flex-direction:column;gap:10px}.budget-highlight{flex-direction:column;gap:14px}}
</style>
</head>
<body>
<div class="wrapper">

  <div class="eyebrow">jelius.dev</div>
  <h1>New work request</h1>
  <p class="subtitle">Received {{.ReceivedAt}}</p>

  <!-- Contact -->
  <div class="card">
    <div class="card-header"><span class="card-header-title">Contact</span></div>
    <div class="card-body">
      <div class="row">
        <div class="col">
          <span class="label">Name</span>
          <span class="value">{{.Name}}</span>
        </div>
        <div class="col">
          <span class="label">Email</span>
          <span class="value mono">{{.Email}}</span>
        </div>
      </div>
      {{if .Phone}}
      <hr class="divider"/>
      <div class="row">
        <div class="col">
          <span class="label">Phone</span>
          <span class="value mono">{{.Phone}}</span>
        </div>
      </div>
      {{end}}
    </div>
  </div>

  <!-- Project -->
  <div class="card">
    <div class="card-header"><span class="card-header-title">Project</span></div>
    <div class="card-body">
      <div class="row">
        <div class="col">
          <span class="label">Service</span>
          <span class="tag">{{.Service}}</span>
        </div>
        <div class="col">
          <span class="label">Title</span>
          <span class="value">{{.Title}}</span>
        </div>
      </div>
      <hr class="divider"/>
      <span class="label">Description</span>
      <p class="prose">{{.Description}}</p>
      {{if .TechPref}}
      <hr class="divider"/>
      <span class="label">Tech preferences</span>
      <p class="value mono">{{.TechPref}}</p>
      {{end}}
    </div>
  </div>

  <!-- Budget & timeline -->
  <div class="card">
    <div class="card-header"><span class="card-header-title">Timeline &amp; budget</span></div>
    <div class="card-body">
      <div class="budget-highlight">
        <div class="budget-block">
          <div class="bval">{{.TimelineDays}}</div>
          <div class="blabel">days requested</div>
        </div>
        <div class="budget-block">
          <div class="bval">₹{{.DailyRate}}</div>
          <div class="blabel">per day offered</div>
        </div>
      </div>
      <p class="total-note">Note: the timeline is the client's estimate and may differ from actual delivery time.</p>
    </div>
  </div>

  {{if .Attachments}}
  <!-- Attachments -->
  <div class="card">
    <div class="card-header"><span class="card-header-title">Attachments ({{len .Attachments}})</span></div>
    <div class="card-body">
      {{range .Attachments}}
      <div class="att-row">
        <span class="att-icon">📎</span>
        <div>
          <div class="att-name">{{.Filename}}</div>
          <div class="att-meta">{{.MIMEType}} &nbsp;·&nbsp; {{.Size}} bytes</div>
        </div>
      </div>
      {{end}}
    </div>
  </div>
  {{end}}

  {{if .Notes}}
  <!-- Notes -->
  <div class="card">
    <div class="card-header"><span class="card-header-title">Additional notes</span></div>
    <div class="card-body">
      <p class="prose">{{.Notes}}</p>
    </div>
  </div>
  {{end}}

  <p class="footer">Sent via jelius.dev/api/work_req</p>

</div>
</body>
</html>`))

// buildEmailHTML renders the HTML body for the work-request email.
func buildEmailHTML(req WorkReq) (string, error) {
    data := emailData{
        Name:         req.Name,
        Email:        req.Email,
        Phone:        req.Phone,
        Service:      req.Service,
        Title:        req.Title,
        Description:  req.Description,
        TechPref:     req.TechPref,
        TimelineDays: req.TimelineDays,
        DailyRate:    req.DailyRate,
        Notes:        req.Notes,
        Attachments:  req.Attachments,
        ReceivedAt:   time.Now().UTC().Add(5*time.Hour + 30*time.Minute).Format("2 Jan 2006, 15:04 IST"),
    }
    var buf bytes.Buffer
    if err := emailHTMLTemplate.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("render email template: %w", err)
    }
    return buf.String(), nil
}

func WorkReqHandler(w http.ResponseWriter, r *http.Request) {
    origin := r.Header.Get("Origin")
    if origin == "https://work.jelius.dev" || origin == "https://freelance.jelius.dev" {
        w.Header().Set("Access-Control-Allow-Origin", origin)
    }
    w.Header().Set("Vary", "Origin")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    if requestedHeaders := r.Header.Get("Access-Control-Request-Headers"); len(requestedHeaders) != 0 {
        w.Header().Set("Access-Control-Allow-Headers", requestedHeaders)
    }

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent)
        return
    }

    get := func(field string) string {
        return strings.TrimSpace(r.FormValue(field))
    }

    email := get("email")
    if email == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{
            "error": "email is required",
        })
        return
    }

    // Rate-limit check (before doing any more work)
    ip := canonicalIP(r)
    key := identityHash(email, ip)

    switch limiter.check(key) {
    case deniedBanned:
        deny(w, "This email address has been permanently blocked from sending work requests.")
        return
    case deniedTooFast:
        deny(w, "Please wait at least 60 seconds between requests.")
        return
    }

    // Parse the multipart form, keeping at most maxFormMemory in RAM.
    if err := r.ParseMultipartForm(maxFormMemory); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{
            "error": "invalid multipart form",
        })
        return
    }
    defer func() {
        if r.MultipartForm != nil {
            _ = r.MultipartForm.RemoveAll()
        }
    }()

    req := WorkReq{
        Name:         get("name"),
        Email:        email,
        Phone:        get("phone"),
        Service:      get("service"),
        Title:        get("title"),
        Description:  get("description"),
        TechPref:     get("tech_pref"),
        TimelineDays: get("timeline_days"),
        DailyRate:    get("daily_rate"),
        Notes:        get("notes"),
    }

    // Required fields — return 400 before admitting so a bad form submit
    // does not consume one of the 6 allowed requests.
    missing := []string{}
    if req.Name == "" {
        missing = append(missing, "name")
    }
    if req.Service == "" {
        missing = append(missing, "service")
    }
    if req.Title == "" {
        missing = append(missing, "title")
    }
    if req.Description == "" {
        missing = append(missing, "description")
    }
    if req.TimelineDays == "" {
        missing = append(missing, "timeline_days")
    }
    if req.DailyRate == "" {
        missing = append(missing, "daily_rate")
    }

    if len(missing) > 0 {
        writeJSON(w, http.StatusBadRequest, map[string]any{
            "error":   "missing required fields",
            "missing": missing,
        })
        return
    }

    // Attachments: write to per-identity temp dir
    //
    // We write each file to os.TempDir()/workreq_<hex-identity>/<safeFilename>
    // so libmailer can attach them by path. The directory is removed after
    // the email send regardless of outcome, keeping temp storage clean.

    var tempDir string // set only when there is at least one attachment

    if r.MultipartForm != nil && r.MultipartForm.File != nil && len(r.MultipartForm.File["attachments"]) > 0 {
        var err error
        tempDir, err = tempDirForIdentity(key)
        if err != nil {
            writeJSON(w, http.StatusInternalServerError, map[string]string{
                "error": "could not prepare upload directory",
            })
            return
        }

        for _, headers := range r.MultipartForm.File["attachments"] {
            if headers.Size > maxFileSize {
                // Clean up any files already written before returning.
                _ = os.RemoveAll(tempDir)
                writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{
                    "error": fmt.Sprintf("file %q exceeds the 10 MB limit", headers.Filename),
                })
                return
            }

            src, err := headers.Open()
            if err != nil {
                _ = os.RemoveAll(tempDir)
                writeJSON(w, http.StatusInternalServerError, map[string]string{
                    "error": "could not read attachment",
                })
                return
            }

            // Sniff MIME type from the first 512 bytes, then seek back to
            // the start so the full file is written to disk.
            sniff := make([]byte, 512)
            n, _ := src.Read(sniff)
            mimeType := http.DetectContentType(sniff[:n])

            destPath := filepath.Join(tempDir, safeFilename(headers.Filename))

            // If two files share the same safe name, differentiate them so
            // neither is silently overwritten.
            if _, statErr := os.Stat(destPath); statErr == nil {
                ext := filepath.Ext(destPath)
                base := strings.TrimSuffix(destPath, ext)
                destPath = fmt.Sprintf("%s_%d%s", base, time.Now().UnixNano(), ext)
            }

            dst, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
            if err != nil {
                _ = src.Close()
                _ = os.RemoveAll(tempDir)
                writeJSON(w, http.StatusInternalServerError, map[string]string{
                    "error": "could not write attachment to disk",
                })
                return
            }

            // Write the sniff buffer first, then stream the remainder.
            if _, err = dst.Write(sniff[:n]); err != nil {
                _ = src.Close()
                _ = dst.Close()
                _ = os.RemoveAll(tempDir)
                writeJSON(w, http.StatusInternalServerError, map[string]string{
                    "error": "could not write attachment to disk",
                })
                return
            }
            if _, err = io.Copy(dst, src); err != nil {
                _ = src.Close()
                _ = dst.Close()
                _ = os.RemoveAll(tempDir)
                writeJSON(w, http.StatusInternalServerError, map[string]string{
                    "error": "could not write attachment to disk",
                })
                return
            }
            _ = src.Close()
            _ = dst.Close()

            req.Attachments = append(req.Attachments, AttachmentMeta{
                Filename: headers.Filename,
                Size:     headers.Size,
                MIMEType: mimeType,
                TempPath: destPath,
            })
        }
    }

    // Admit (consume one slot)
    // Only reached if all validation passed and files are safely on disk.
    limiter.admit(key)

    // Ensure temp files are always removed, whether the send succeeds or not.
    if tempDir != "" {
        defer os.RemoveAll(tempDir)
    }

    // Build email HTML

    htmlBody, err := buildEmailHTML(req)
    if err != nil {
        logger.Error(err)
        http.Error(w, "Something went wrong!", http.StatusInternalServerError)
        return
    }

    cfg, err := libmailer.LoadConfig()
    personalEmail := os.Getenv("PERSONAL_EMAIL")
    workEmail := os.Getenv("WORK_EMAIL")

    if (err != nil) || (len(personalEmail) == 0 && len(workEmail) == 0) {
        logger.Error(err)
        http.Error(w, "Something went wrong!", http.StatusInternalServerError)
        return
    }

    // Collect temp paths for libmailer — the files already exist on disk.
    attachPaths := make([]string, 0, len(req.Attachments))
    for _, a := range req.Attachments {
        attachPaths = append(attachPaths, a.TempPath)
    }

    err = libmailer.SendMail(
        cfg,
        cfg.From,
        workEmail,
        fmt.Sprintf("Work Req: %s", req.Title),
        htmlBody,
        []string{personalEmail},
        nil,
        attachPaths,
    )
    if err != nil {
        logger.Error(err)
        http.Error(w, "Something went wrong!", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusAccepted)
}

