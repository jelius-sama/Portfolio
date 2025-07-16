package main

import (
	"KazuFolio/api"
	"KazuFolio/logger"
	m "KazuFolio/middleware"
	"errors"
	"net/http"
	"os"
	"path/filepath"
)

// TODO: Cleanup Server Side Code (eg.: http error response re-using function)
// TODO: Cleanup Client Side Code
// TODO: Write Home page with new UI design
// TODO: Work on Metadata for links page
// TODO: Perfect SEO & Metadata for all pages
// TODO: Add PWA Support
// TODO: Test for unexpected behaviors & test the code
// TODO: Check if the server sent data is accurate (eg.: metadata parsing)

var Environment = "development"
var Port = "6969"
var Version string

func init() {
	exePath, err := os.Executable()
	if err != nil {
		logger.Panic("could not get executable path:", err)
	}

	os.Setenv("ROOT_PATH", filepath.Dir(filepath.Dir(exePath)))
	os.Setenv("version", Version)
	os.Setenv("env", Environment)
	os.Setenv("port", Port)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true // File exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false // File does not exist
	}
	// Handle other potential errors (e.g., permissions)
	logger.Panic("Error checking file `"+filePath+"`\n    ", err)
	return false
}

func main() {
	startServer := func() (error, string) {
		portToStart := Port
		routeHandler := m.Chain(m.MiddlewareChain{
			Handler: api.HandleRouting(),
			Middlewares: []m.Middleware{
				m.RecoveryMiddleware,
				m.LoggingMiddleware,
			},
		})

		if Environment == "development" {
			portToStart = Port
			os.Setenv("port", portToStart)

			logger.Info("Server started on port :" + portToStart)
			return http.ListenAndServe(":"+portToStart, routeHandler), portToStart
		} else {
			fullchain := "/etc/letsencrypt/live/jelius.dev/fullchain.pem"
			privkey := "/etc/letsencrypt/live/jelius.dev/privkey.pem"
			exists1 := fileExists(fullchain)
			exists2 := fileExists(privkey)

			if exists1 && exists2 {
				portToStart = "443"
				os.Setenv("port", portToStart)

				logger.Info("Server started on port :" + portToStart)
				return http.ListenAndServeTLS(":"+portToStart, fullchain, privkey, routeHandler), portToStart
			} else {
				portToStart = Port
				os.Setenv("port", portToStart)

				logger.Warning("Production server was started without SSL certificate falling back to http")

				logger.Info("Server started on port :" + portToStart)
				return http.ListenAndServe(":"+portToStart, routeHandler), portToStart
			}
		}
	}

	if err, port := startServer(); err != nil {
		logger.Panic("Could not start the server on port :"+port, "\n", err)
	}
}
