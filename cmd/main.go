package main

import (
	"KazuFolio/api"
	"KazuFolio/db"
	"KazuFolio/logger"
	m "KazuFolio/middleware"
	"database/sql"
	"errors"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var Environment = "development"
var Port = "6969"
var Version string
var Home = "/home/ec2-user"

func init() {
	exePath, err := os.Executable()
	if err != nil {
		logger.Panic("could not get executable path:", err)
	}

	os.Setenv("ROOT_PATH", filepath.Dir(filepath.Dir(exePath)))
	os.Setenv("version", Version)
	os.Setenv("env", Environment)
	os.Setenv("port", Port)
	os.Setenv("home", Home)
	os.Setenv("db_file", path.Join(Home, "/portfolio.db"))

	db.Conn, err = sql.Open("sqlite", os.Getenv("db_file"))
	if err != nil {
		logger.Panic(err)
	}

	if err = db.InitializeSchema(db.Conn); err != nil {
		logger.Panic(err)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
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

	defer db.Conn.Close()

	if err, port := startServer(); err != nil {
		logger.Panic("Could not start the server on port :"+port, "\n", err)
	}
}
