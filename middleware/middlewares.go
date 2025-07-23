package middleware

import (
	"KazuFolio/api"
	"KazuFolio/logger"
	"net/http"
	"strings"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.TimedInfo(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.TimedError("Encountered a panic, returning 500 to client and recovering the server!")
				if strings.HasPrefix(r.URL.Path, "/api/") {
					api.InternalErrorAPI(w, r, nil)
					return
				}

				if strings.HasPrefix(r.URL.Path, "/assets/") {
					api.InternalErrorAPI(w, r, nil)
					return
				}

				if strings.HasPrefix(r.URL.Path, "/src/") {
					api.InternalErrorAPI(w, r, nil)
					return
				}

				api.InternalErrorPage(w, r, nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
