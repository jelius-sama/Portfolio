package handler

import (
	"KazuFolio/db"
	"encoding/xml"
	"net/http"
	"os"
	"time"
)

type urlEntry struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

type urlSet struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	URLs    []urlEntry `xml:"url"`
}

func GenerateSitemap(w http.ResponseWriter, r *http.Request) {
	baseURL := os.Getenv("host")
	now := time.Now().Format("2006-01-02")

	// Static routes
	urls := []urlEntry{
		{Loc: baseURL + "/", LastMod: now, ChangeFreq: "daily", Priority: "1.0"},
		{Loc: baseURL + "/blogs", LastMod: now, ChangeFreq: "daily", Priority: "0.8"},
		{Loc: baseURL + "/links", LastMod: now, ChangeFreq: "monthly", Priority: "0.6"},
	}

	// Fetch all blogs directly
	rows, err := db.Conn.Query(`
		SELECT id, created_at, updated_at
		FROM blogs
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, createdAt, updatedAt string
		if err := rows.Scan(&id, &createdAt, &updatedAt); err != nil {
			continue // skip malformed rows
		}

		lastMod := updatedAt
		if lastMod == "" {
			lastMod = createdAt
		}

		t, err := time.Parse(time.RFC3339, lastMod)
		if err != nil {
			lastMod = now
		} else {
			lastMod = t.Format("2006-01-02")
		}

		urls = append(urls, urlEntry{
			Loc:        baseURL + "/blog/" + id,
			LastMod:    lastMod,
			ChangeFreq: "monthly",
			Priority:   "0.7",
		})
	}

	w.Header().Set("Content-Type", "application/xml")

	_ = xml.NewEncoder(w).Encode(urlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	})
}
