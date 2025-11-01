package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"KazuFolio/db"
	"KazuFolio/types"
	"errors"
	"github.com/jelius-sama/logger"
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

func GetBlogInternal(id string) (*types.Blog, int, error) {
	// Channel to receive blog data result
	type blogResult struct {
		blog   *types.Blog
		status int
		err    error
	}
	blogChan := make(chan blogResult, 1)

	// Channel to receive views result
	type viewsResult struct {
		views uint
		err   error
	}
	viewsChan := make(chan viewsResult, 1)

	// Goroutine to fetch blog data
	go func() {
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
			blogChan <- blogResult{nil, http.StatusNotFound, errors.New("Blog not found")}
			return
		} else if err != nil {
			blogChan <- blogResult{nil, http.StatusInternalServerError, errors.New("Database error: " + err.Error())}
			return
		}

		// Decode JSON-encoded parts list
		if err := json.Unmarshal([]byte(parts), &blog.Parts); err != nil {
			blogChan <- blogResult{nil, http.StatusInternalServerError, errors.New("Failed to parse parts list: " + err.Error())}
			return
		}

		blogChan <- blogResult{&blog, http.StatusOK, nil}
	}()

	// Goroutine to fetch blog views
	go func() {
		path := "/blog/" + id
		views, err := PageViewsInternal(path)
		viewsChan <- viewsResult{views, err}
	}()

	// Wait for blog data result
	blogRes := <-blogChan
	if blogRes.err != nil {
		return blogRes.blog, blogRes.status, blogRes.err
	}

	// Wait for views result
	viewsRes := <-viewsChan
	if viewsRes.err != nil {
		logger.TimedError("Error fetching views for blog", id, ":", viewsRes.err)
		// Views remains zero value on error
	}
	blogRes.blog.Views = viewsRes.views

	return blogRes.blog, blogRes.status, nil
}

func GetBlog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID for the blog was not provided!", http.StatusBadRequest)
		return
	}

	blog, status, err := GetBlogInternal(id)

	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		http.Error(w, "Failed to encode blog as JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
