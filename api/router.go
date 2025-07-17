package api

import (
	embed "KazuFolio"
	"KazuFolio/logger"
	"KazuFolio/parser"
	"KazuFolio/types"
	"KazuFolio/util"
	"bytes"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var getOnlyRoute = util.AddrOf("Only Request with GET Method are allowed!")

func HandleRouting() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			MethodNotAllowed(w, r, getOnlyRoute)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		content, err := fs.ReadFile(embed.AssetsFS, path)
		if err != nil {
			NotFound(w, r, util.AddrOf("Requested Asset was not found"))
			return
		}

		ext := filepath.Ext(path)
		mimeType := mime.TypeByExtension(ext)
		if len(mimeType) == 0 {
			mimeType = "application/octet-stream"
		}

		w.Header().Set("Content-Type", mimeType)
		w.Write(content)
	})

	router.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			_, file, line, ok := runtime.Caller(1)
			NotFound(w, r, util.AddrOf("Unreachable code reached!"))
			if ok {
				logger.Error("Reached unreachable code in `", file, "` at line: ", line)
			} else {
				logger.Error("Reached unreachable code (unknown location)")
			}
			return
		}

		method := r.Method
		path := strings.TrimPrefix(r.URL.Path, "/api/")
		lookupKey := method + " /" + path

		if handler, exists := ApiRoutes[lookupKey]; exists {
			handler(w, r)
			return
		}

		NotFound(w, r, util.AddrOf("API Route Not Found!"))
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			MethodNotAllowed(w, r, getOnlyRoute)
			return
		}

		html, err := parser.GetHTML()
		if err != nil {
			InternalServerError(w, r, util.AddrOf("Something went wrong when getting html from FS!"))
			logger.Error("failed to get html shell:\n    " + err.Error())
			return
		}

		metadata, err := parser.ParseMetadata(r.URL.Path)
		if err != nil {
			InternalServerError(w, r, util.AddrOf("Failed to parse metadata of the page!"))
			logger.Error("metadata parsing failed:\n    " + err.Error())
			return
		}

		// Replace marker in HTML
		html = bytes.Replace(html, []byte("<!-- Server Props -->"), []byte(metadata), 1)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(html)
	})

	if os.Getenv("env") == types.ENV.Prod {
		router.HandleFunc("/src/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				MethodNotAllowed(w, r, getOnlyRoute)
				return
			}

			path := strings.TrimPrefix(r.URL.Path, "/")
			content, err := fs.ReadFile(embed.ViteFS, "client/dist/"+path)
			if err != nil {
				NotFound(w, r, util.AddrOf("Requested Source file was not found!"))
				return
			}

			ext := filepath.Ext(path)
			mimeType := mime.TypeByExtension(ext)
			if len(mimeType) == 0 {
				mimeType = "application/octet-stream"
			}

			w.Header().Set("Content-Type", mimeType)
			w.Write(content)
		})
	}

	return router
}
