package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const owner = "jelius-sama"

func LatestCommit(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	repo := r.URL.Query().Get("repo")
	branch := r.URL.Query().Get("branch")

	if repo == "" || branch == "" {
		http.Error(w, "Missing required query parameters: repo and branch", http.StatusBadRequest)
		return
	}

	// Build GitHub API URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, branch)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to reach GitHub API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("GitHub API error: %s", resp.Status), resp.StatusCode)
		return
	}

	var commit struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Name  string `json:"name"`
				Email string `json:"email"`
				Date  string `json:"date"`
			} `json:"author"`
		} `json:"commit"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		http.Error(w, "Failed to parse GitHub response", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commit)
}
