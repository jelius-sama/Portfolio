package handler

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "os"
)

type captureWriter struct {
    header http.Header
    body   bytes.Buffer
    status int
}

func newCaptureWriter() *captureWriter {
    return &captureWriter{
        header: make(http.Header),
        status: http.StatusOK,
    }
}

func (cw *captureWriter) Header() http.Header {
    return cw.header
}

func (cw *captureWriter) Write(b []byte) (int, error) {
    return cw.body.Write(b)
}

func (cw *captureWriter) WriteHeader(status int) {
    cw.status = status
}

func SpecializedTask(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()

    task := q.Get("task")
    token := q.Get("token")

    switch task {
    case "update_server":

        stored := GetToken()
        if stored == nil || token != *stored {
            http.Error(w, "Invalid or expired token", http.StatusForbidden)
            return
        }

        payload := Payload{Token: token}
        data, _ := json.Marshal(payload)

        req, _ := http.NewRequest("POST", "/api/update_server", io.NopCloser(bytes.NewReader(data)))
        req.ContentLength = int64(len(data))
        req.Header.Set("Content-Type", "application/json")

        sudoPass := os.Getenv("SUDO_PASS")
        req.Header.Set("Authorization", "Bearer "+sudoPass)

        // Clean remote addr
        req.RemoteAddr = r.RemoteAddr

        cw := newCaptureWriter()
        UpdateServer(cw, req)

        if cw.status == http.StatusNoContent {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("Update request accepted. Server may restart shortly."))
            return
        }

        for k, v := range cw.header {
            for _, vv := range v {
                w.Header().Add(k, vv)
            }
        }
        w.WriteHeader(cw.status)
        w.Write(cw.body.Bytes())

    default:
        http.Error(w, "Invalid task", http.StatusBadRequest)
    }
}
