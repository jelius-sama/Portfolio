package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"KazuFolio/db"
	"KazuFolio/types"
)

func GetBlogs(w http.ResponseWriter, r *http.Request) {
	// Default pagination values
	page := 1
	limit := 10
	sortBy := "created_at"
	order := "desc"

	// Parse query parameters
	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}
	if s := r.URL.Query().Get("sortBy"); s != "" {
		sortBy = s
	}
	if o := r.URL.Query().Get("order"); o != "" {
		o = strings.ToLower(o)
		if o == "asc" || o == "desc" {
			order = o
		}
	}

	// Validate sortBy to prevent SQL injection
	validSortFields := map[string]bool{
		"id": true, "title": true, "created_at": true, "updated_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	offset := (page - 1) * limit

	query := `
		SELECT id, title, summary, markdown_url, created_at, updated_at, prequel_id, sequel_id, parts
		FROM blogs
		ORDER BY ` + sortBy + ` ` + order + `
		LIMIT ? OFFSET ?`

	rows, err := db.Conn.Query(query, limit, offset)
	if err != nil {
		http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var blogs []types.Blog

	for rows.Next() {
		var blog types.Blog
		var parts string

		if err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Summary,
			&blog.MarkdownURL,
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.PrequelID,
			&blog.SequelID,
			&parts); err != nil {
			http.Error(w, "Failed to scan blog: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal([]byte(parts), &blog.Parts); err != nil {
			http.Error(w, "Failed to parse parts: "+err.Error(), http.StatusInternalServerError)
			return
		}

		blogs = append(blogs, blog)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Row iteration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blogs); err != nil {
		http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
