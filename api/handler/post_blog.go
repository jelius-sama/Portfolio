package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"KazuFolio/db"
	"KazuFolio/types"

	"KazuFolio/util"
	"github.com/google/uuid"
)

func generateBlogID() string {
	timestamp := time.Now().UTC().Format("060102150405") // YYMMDDHHMMSS
	return timestamp + "-" + uuid.New().String()[:6]
}

// TODO: Send Email to admin for approval
func PostBlog(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Invalid multipart request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Hold parsed form values
	blog := types.Blog{}
	var contentFilePath string
	var parts []string

	// GitHub-style ID based on timestamp and UUID suffix
	blog.ID = generateBlogID()
	blog.CreatedAt = time.Now().String()
	blog.UpdatedAt = blog.CreatedAt

	// Parse multipart parts
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading multipart data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		name := part.FormName()
		if name == "" {
			continue
		}

		switch name {
		case "title":
			buf, _ := io.ReadAll(part)
			blog.Title = string(buf)

		case "summary":
			buf, _ := io.ReadAll(part)
			blog.Summary = string(buf)

		case "prequel_id":
			buf, _ := io.ReadAll(part)
			blog.PrequelID = util.AddrOf(string(buf))

		case "sequel_id":
			buf, _ := io.ReadAll(part)
			blog.SequelID = util.AddrOf(string(buf))

		case "parts":
			buf, _ := io.ReadAll(part)
			parts = strings.Split(string(buf), ",")

		case "content_file":
			// Save uploaded file to ~/Blogs/{id}.md
			home := os.Getenv("home")
			if home == "" {
				http.Error(w, "Missing 'home' environment variable", http.StatusInternalServerError)
				return
			}
			dir := filepath.Join(home, "Blogs")
			if err := os.MkdirAll(dir, 0755); err != nil {
				http.Error(w, "Failed to create Blogs directory: "+err.Error(), http.StatusInternalServerError)
				return
			}

			contentFilePath = filepath.Join(dir, blog.ID+".md")
			out, err := os.Create(contentFilePath)
			if err != nil {
				http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer out.Close()

			if _, err := io.Copy(out, part); err != nil {
				http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if contentFilePath == "" {
		http.Error(w, "Missing blog content file", http.StatusBadRequest)
		return
	}

	blog.MarkdownURL = contentFilePath
	blog.Parts = parts

	// Insert blog into the database
	stmt, err := db.Conn.Prepare(`
		INSERT INTO blogs (
			id, title, summary, markdown_url,
			created_at, updated_at, prequel_id, sequel_id, parts
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		http.Error(w, "Failed to prepare statement: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	partsJSON, err := json.Marshal(blog.Parts)
	if err != nil {
		http.Error(w, "Failed to encode parts list: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(
		blog.ID,
		blog.Title,
		blog.Summary,
		blog.MarkdownURL,
		blog.CreatedAt,
		blog.UpdatedAt,
		blog.PrequelID,
		blog.SequelID,
		string(partsJSON),
	)
	if err != nil {
		http.Error(w, "Failed to insert blog: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}
