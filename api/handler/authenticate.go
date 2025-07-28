package handler

import (
	"KazuFolio/logger"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"
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

// Authenticate creates and stores a token, sets it as a cookie
func Authenticate(w http.ResponseWriter, r *http.Request) {
	if !VerifySudo(w, r) {
		return
	}

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authenticated"))
}

// Internal function to verify auth token
func VerifyAuthToken(token string) error {
	authTokens.RLock()
	expiry, exists := authTokens.store[token]
	authTokens.RUnlock()

	if !exists {
		return fmt.Errorf("token not found")
	}

	if time.Now().After(expiry) {
		authTokens.Lock()
		delete(authTokens.store, token)
		authTokens.Unlock()
		return fmt.Errorf("token expired")
	}

	// Renew token
	newExpiry := time.Now().Add(tokenValidity)
	authTokens.Lock()
	authTokens.store[token] = newExpiry
	authTokens.Unlock()

	return nil
}

// HTTP handler that uses the internal function
func VerifyAuthStatus(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "Missing auth token", http.StatusForbidden)
		return
	}
	token := cookie.Value

	err = VerifyAuthToken(token)
	if err != nil {
		if err.Error() == "token expired" {
			http.Error(w, err.Error(), 498)
		} else {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
		return
	}

	// Renewed cookie
	newExpiry := time.Now().Add(tokenValidity)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  newExpiry,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authorized"))
}
