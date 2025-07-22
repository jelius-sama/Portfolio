package api

import (
	"KazuFolio/logger"
	"KazuFolio/parser"
	"KazuFolio/util"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error   string  `json:"error"`
	Message *string `json:"message"`
}

func writeError(w http.ResponseWriter, statusCode int, errorText string, msg *string) {
	resp := ErrorResponse{
		Error:   errorText,
		Message: msg,
	}

	util.WriteJSON(w, statusCode, resp)
}

func BadRequest(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusBadRequest, "400 Bad Request", msg)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusUnauthorized, "401 Unauthorized", msg)
}

func Forbidden(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusForbidden, "403 Forbidden", msg)
}

func NotFound(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusNotFound, "404 Not Found", msg)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusMethodNotAllowed, "405 Method Not Allowed", msg)
}

func RequestTimeout(w http.ResponseWriter, r *http.Request, msg *string) {
	writeError(w, http.StatusRequestTimeout, "408 Request Timeout", msg)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, msg *string) {
	// Attempt to get the HTML shell
	html, err := parser.GetHTML()
	if err != nil {
		// If even the shell fails, fall back to basic response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Error("failed to get html shell in error handler:\n    " + err.Error())
		return
	}

	// Create the SSR data as JSON
	payload := map[string]interface{}{
		"status":  500,
		"message": "Internal Server Error",
	}
	if msg != nil {
		payload["message"] = *msg
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Error("failed to marshal SSR error JSON:\n    " + err.Error())
		return
	}

	// Build the <script> tag with ID __SERVER_DATA__
	ssrScript := fmt.Sprintf(`<script id="__SERVER_DATA__" type="application/json">%s</script>`, jsonData)

	// Inject into the HTML shell at the SSR Data marker
	finalHTML := bytes.Replace(html, []byte("<!-- SSR Data -->"), []byte(ssrScript), 1)

	// Write the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(finalHTML)
}
