package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"KazuFolio/db"
	"KazuFolio/types"
	"path/filepath"
	"regexp"
)

// Regex to allow only alphanumeric characters, hyphens, and underscores
var safeID = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func GetBlogFile(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID for the blog was not provided!", http.StatusBadRequest)
		return
	}

	// Sanitize the ID: reject anything with unsafe characters
	if !safeID.MatchString(id) {
		http.Error(w, "Invalid blog ID format.", http.StatusBadRequest)
		return
	}

	// Construct the path
	blogsDir := filepath.Join(os.Getenv("home"), "Blogs")
	blogFile := filepath.Join(blogsDir, id+".md")

	// Open and serve the file
	file, err := os.Open(blogFile)
	if err != nil {
		http.Error(w, "Could not find the blog file.", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set MIME type
	w.Header().Set("Content-Type", "text/markdown")
	http.ServeFile(w, r, blogFile)
}

func GetBlog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID for the blog was not provided!", http.StatusBadRequest)
		return
	}

	var (
		blog  types.Blog
		parts string // Will be parsed from JSON string into blog.Parts
	)

	query := `
		SELECT id, title, summary, created_at, updated_at, prequel_id, sequel_id, parts
		FROM blogs WHERE id = ?`
	err := db.Conn.QueryRow(query, id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Summary,
		&blog.CreatedAt,
		&blog.UpdatedAt,
		&blog.PrequelID,
		&blog.SequelID,
		&parts,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode JSON-encoded parts list
	if err := json.Unmarshal([]byte(parts), &blog.Parts); err != nil {
		http.Error(w, "Failed to parse parts list: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		http.Error(w, "Failed to encode blog as JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
