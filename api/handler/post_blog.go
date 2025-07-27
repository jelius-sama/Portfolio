package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"KazuFolio/db"
	"KazuFolio/logger"
	"KazuFolio/types"
	"KazuFolio/util"
)

func generateBlogID() string {
	// Create a unique string combining timestamp and some randomness
	timestamp := time.Now().UTC().UnixNano()
	unique := fmt.Sprintf("%d-%d", timestamp, time.Now().UTC().UnixMilli())

	// Generate SHA1 hash and take first 7 characters (like GitHub)
	hash := sha1.Sum([]byte(unique))
	return hex.EncodeToString(hash[:])[:7]
}

// TODO: Send Email to admin for approval
func PostBlog(w http.ResponseWriter, r *http.Request) {
	// Verify sudo access
	if !VerifySudo(w, r) {
		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Invalid multipart request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Hold parsed form values
	blog := types.Blog{}
	var parts []string

	// GitHub-style ID based on hash
	blog.ID = generateBlogID()
	now := time.Now().UTC().Format(time.RFC3339)
	blog.CreatedAt = now
	blog.UpdatedAt = now

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
			buf, err := io.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading title: "+err.Error(), http.StatusInternalServerError)
				return
			}
			blog.Title = string(buf)

		case "summary":
			buf, err := io.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading summary: "+err.Error(), http.StatusInternalServerError)
				return
			}
			blog.Summary = string(buf)

		case "prequel_id":
			buf, err := io.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading prequel_id: "+err.Error(), http.StatusInternalServerError)
				return
			}
			s := strings.TrimSpace(string(buf))
			if s != "" {
				blog.PrequelID = util.AddrOf(s)
			}

		case "sequel_id":
			buf, err := io.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading sequel_id: "+err.Error(), http.StatusInternalServerError)
				return
			}
			s := strings.TrimSpace(string(buf))
			if s != "" {
				blog.SequelID = util.AddrOf(s)
			}

		case "parts":
			buf, err := io.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading parts: "+err.Error(), http.StatusInternalServerError)
				return
			}
			partsStr := strings.TrimSpace(string(buf))
			if partsStr != "" {
				parts = strings.Split(partsStr, ",")
				// Trim whitespace from each part
				for i, part := range parts {
					parts[i] = strings.TrimSpace(part)
				}
			}

		case "content_file":
			// Save uploaded file to ~/Blogs/{id}.md
			home := os.Getenv("home")
			if len(home) == 0 {
				http.Error(w, "Something went wrong!", http.StatusInternalServerError)
				logger.TimedPanic("Missing 'home' environment variable")
				return
			}
			dir := filepath.Join(home, "Blogs")
			if err := os.MkdirAll(dir, 0755); err != nil {
				http.Error(w, "Failed to create Blogs directory: "+err.Error(), http.StatusInternalServerError)
				logger.TimedPanic("Failed to create Blogs directory: " + err.Error())
				return
			}

			contentFilePath := filepath.Join(dir, blog.ID+".md")
			out, err := os.Create(contentFilePath)
			if err != nil {
				http.Error(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError)
				logger.TimedError("Failed to create file: " + err.Error())
				return
			}
			defer out.Close()

			if _, err := io.Copy(out, part); err != nil {
				http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
				logger.TimedError("Failed to save file: " + err.Error())
				return
			}
		}

		// Close the part to free resources
		part.Close()
	}

	// Validation: Check required fields
	if blog.Title == "" {
		http.Error(w, "Blog title is required", http.StatusBadRequest)
		return
	}

	blog.Parts = parts

	// Insert blog into the database
	stmt, err := db.Conn.Prepare(`
		INSERT INTO blogs (
			id, title, summary,
			created_at, updated_at, prequel_id, sequel_id, parts
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		http.Error(w, "Failed to prepare statement: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var partsJSON []byte
	if len(blog.Parts) > 0 {
		partsJSON, err = json.Marshal(blog.Parts)
		if err != nil {
			http.Error(w, "Failed to encode parts list: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		partsJSON = []byte("[]") // Empty JSON array for null parts
	}

	_, err = stmt.Exec(
		blog.ID,
		blog.Title,
		blog.Summary,
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

	util.WriteJSON(w, http.StatusCreated, blog)
}
