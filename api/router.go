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
	"strings"
)

var getOnlyRoute = util.AddrOf("Only Request with GET Method are allowed!")

func HandleRouting() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			Forbidden(w, r, getOnlyRoute)
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
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		w.Header().Set("Content-Type", mimeType)
		w.Write(content)
	})

	router.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			NotFound(w, r, nil)
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
			http.Error(w, "Error: Method not allowed!", http.StatusMethodNotAllowed)
			return
		}

		html, err := parser.GetHTML()
		if err != nil {
			http.Error(w, "Error: Internal Server Error!", http.StatusInternalServerError)
			return
		}

		metadata, err := parser.ParseMetadata()
		if err != nil {
			http.Error(w, "Error: Failed to parse metadata", http.StatusInternalServerError)
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
				Forbidden(w, r, getOnlyRoute)
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
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			w.Header().Set("Content-Type", mimeType)
			w.Write(content)
		})
	}

	return router
}
