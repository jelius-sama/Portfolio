package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"KazuFolio/db"
	"KazuFolio/types"
)

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
		SELECT id, title, summary, markdown_url, created_at, updated_at, prequel_id, sequel_id, parts
		FROM blogs WHERE id = ?`
	err := db.Conn.QueryRow(query, id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Summary,
		&blog.MarkdownURL,
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
