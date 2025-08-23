package main

import (
	vars "KazuFolio"
	"KazuFolio/api"
	"KazuFolio/db"
	"KazuFolio/logger"
	"KazuFolio/types"
	"KazuFolio/util"
	"database/sql"
	"errors"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"path/filepath"
)

var (
	Environment = "development"
	DevPort     = "6969"
	Version     string
	Home        string
	Host        = "http://localhost"

	ReverseProxy string
	ProxyPort    string
)

// NOTE: Be sure to explicitly set the home path in the `../config/server.config.json` file.
//   - This is important because if the server runs directly without a reverse proxy,
//     it may require root privileges, in which case `os.UserHomeDir()` may not resolve
//     to the intended user's home directory.
//   - Conversely, if the server does not run with root privileges, it will be unable
//     to bind to ports 443 or 80. This may force users to specify a port number
//     in the URL when accessing your domain.
func init() {
	exePath, err := os.Executable()
	if err != nil {
		logger.Panic("could not get executable path:", err)
	}

	if Environment == types.ENV.Dev && Home == "" {
		if h, err := os.UserHomeDir(); err == nil {
			Home = h
		}
	}

	if ReverseProxy == "true" {
		if util.IsValidPort(ProxyPort) == false {
			logger.Panic("supplied port for reverse proxy is invalid.")
		}

		vars.ReverseProxy = types.BehindReverseProxy{
			StatementValid: true,
			Port:           ProxyPort,
		}
	} else {
		vars.ReverseProxy = types.BehindReverseProxy{
			StatementValid: false,
		}
	}

	os.Setenv("ROOT_PATH", filepath.Dir(filepath.Dir(exePath)))
	os.Setenv("version", Version)
	os.Setenv("env", Environment)
	os.Setenv("home", Home)
	os.Setenv("host", Host)
	os.Setenv("db_file", filepath.Join(Home, "/portfolio.db"))

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

type ServerResp struct {
	Http  error
	Https error
}

// TODO: Implement automatic cache purge when a new build is detected.
func main() {
	defer db.Conn.Close()

	// INFO:: startServer checks the current environment configuration.
	//         - In development mode, it starts the server on the DevPort.
	//         - In production mode:
	//           - If behind a reverse proxy, it starts the server on ReverseProxy.Port.
	//           - Otherwise, it attempts to start the server with HTTPS on port 443 and HTTP on port 80.
	//             The HTTP server only forwards requests to HTTPS.
	//           - If SSL certificates are not found for port 443, the server defaults to starting on port 80.
	startServer := func() (ServerResp, string) {
		respChan := make(chan ServerResp, 2)

		routeHandler := util.Chain(util.MiddlewareChain{
			Handler: api.HandleRouting(),
			Middlewares: []types.Middleware{
				api.RecoveryMiddleware,
				api.LoggingMiddleware,
			},
		})

		// INFO: Development Server
		if Environment == types.ENV.Dev {
			if util.IsValidPort(DevPort) == false {
				logger.Panic("supplied port for dev server is invalid, falling back to :6969")
				DevPort = "6969"
			}

			os.Setenv("port", DevPort)
			logger.Info("Development server started on port :" + DevPort)

			go func() {
				err := http.ListenAndServe(":"+DevPort, routeHandler)
				respChan <- ServerResp{Http: err}
			}()

			return <-respChan, DevPort
		}

		// TODO: Support secure connection over https
		if vars.ReverseProxy.StatementValid == true {
			os.Setenv("port", vars.ReverseProxy.Port)
			logger.Info("Production server started behind reverse proxy on port :" + vars.ReverseProxy.Port)

			go func() {
				err := http.ListenAndServe(":"+vars.ReverseProxy.Port, routeHandler)
				respChan <- ServerResp{Http: err}
			}()

			return <-respChan, vars.ReverseProxy.Port
		}

		// INFO: Production server
		fullchain := os.Getenv("FULL_CHAIN")
		privkey := os.Getenv("PRIV_KEY")

		exists1 := fileExists(fullchain)
		exists2 := fileExists(privkey)

		if exists1 && exists2 {
			portToStart := "443"
			os.Setenv("port", portToStart)
			logger.Info("Production server started on port :" + portToStart)

			go func() {
				err := http.ListenAndServeTLS(":"+portToStart, fullchain, privkey, routeHandler)
				respChan <- ServerResp{Https: err}
			}()

			go func() {
				err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					target := "https://" + r.Host + r.URL.RequestURI()
					http.Redirect(w, r, target, http.StatusMovedPermanently)
				}))
				respChan <- ServerResp{Http: err}
			}()

			var combined ServerResp
			for i := 0; i < 2; i++ {
				resp := <-respChan
				if resp.Http != nil {
					combined.Http = resp.Http
				}
				if resp.Https != nil {
					combined.Https = resp.Https
				}
			}

			return combined, portToStart
		}

		// INFO: If no SSL certificates were found, fallback to HTTP only mode
		fallbackPort := "80"
		os.Setenv("port", fallbackPort)
		logger.Warning("Production server was started without SSL certificate falling back to http only mode.")
		logger.Info("Production server started on port :" + fallbackPort)

		go func() {
			err := http.ListenAndServe(":"+fallbackPort, routeHandler)
			respChan <- ServerResp{Http: err}
		}()

		return <-respChan, fallbackPort
	}

	if serverResp, port := startServer(); serverResp.Http != nil || serverResp.Https != nil {
		if serverResp.Http != nil {
			if vars.ReverseProxy.StatementValid == true {
				logger.Error("Server failed behind reverse proxy on port " + port + ": " + serverResp.Http.Error())
			} else {
				logger.Error("HTTP server failed on port " + port + ":\n" + serverResp.Http.Error())
			}
		}
		if serverResp.Https != nil {
			logger.Error("HTTPS server failed on port " + port + ":\n" + serverResp.Https.Error())
		}
		logger.Panic("One or more servers failed to start.")
	}
}
