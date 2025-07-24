package handler

import (
	"KazuFolio/logger"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"KazuFolio/util"
)

const tokenValidity = 30 * time.Minute

var (
	authTokens = struct {
		sync.RWMutex
		store map[string]time.Time
	}{
		store: make(map[string]time.Time),
	}
)

// init starts the token GC routine
func init() {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()
			authTokens.Lock()
			for token, expiry := range authTokens.store {
				if now.After(expiry) {
					delete(authTokens.store, token)
				}
			}
			authTokens.Unlock()
			logger.TimedInfo("Expired tokens cleaned by AuthToken GC.")
		}
	}()
}

// generateSecureToken returns a securely generated random token string
func generateSecureToken(nBytes int) (string, error) {
	bytes := make([]byte, nBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Authenticate creates and stores a token, sets it as a cookie, and sends it to the client
func Authenticate(w http.ResponseWriter, r *http.Request) {
	util.VerifySudo(w, r)

	token, err := generateSecureToken(32)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	expiry := time.Now().Add(tokenValidity)

	authTokens.Lock()
	authTokens.store[token] = expiry
	authTokens.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiry,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// VerifyAuthStatus validates token from cookie and query param, renews if valid
func VerifyAuthStatus(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "Missing auth token", http.StatusForbidden)
		return
	}
	cookieToken := cookie.Value

	queryToken := r.URL.Query().Get("token")
	if queryToken == "" || queryToken != cookieToken {
		http.Error(w, "Token mismatch", http.StatusForbidden)
		return
	}

	authTokens.RLock()
	expiry, exists := authTokens.store[cookieToken]
	authTokens.RUnlock()

	if !exists {
		http.Error(w, "Token not found", http.StatusForbidden)
		return
	}

	if time.Now().After(expiry) {
		authTokens.Lock()
		delete(authTokens.store, cookieToken)
		authTokens.Unlock()
		http.Error(w, "Token expired", 498)
		return
	}

	// Renew token
	newExpiry := time.Now().Add(tokenValidity)
	authTokens.Lock()
	authTokens.store[cookieToken] = newExpiry
	authTokens.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    cookieToken,
		Expires:  newExpiry,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authorized"))
}
