package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"path/filepath"
)

func GetAchievementsFile(w http.ResponseWriter, r *http.Request) {
	// Construct the path
	achievementsFile := filepath.Join(os.Getenv("home"), "Achievements.md")

	// Open and serve the file
	file, err := os.Open(achievementsFile)
	if err != nil {
		http.Error(w, "Could not find the blog file.", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set MIME type
	w.Header().Set("Content-Type", "text/markdown")
	http.ServeFile(w, r, achievementsFile)
}

type AchievementsStats struct {
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
	Views     uint  `json:"views"`
}

// Hardcoded fallback: 2025-09-13 00:00:00 UTC
const fallbackDate int64 = 1757731200

func GetAchievementsStatsInternal() AchievementsStats {
	views, _ := PageViewsInternal("/achievements")
	achievementsFile := filepath.Join(os.Getenv("home"), "Achievements.md")

	stats := AchievementsStats{
		CreatedAt: fallbackDate,
		UpdatedAt: fallbackDate,
		Views:     views,
	}

	if fi, err := os.Stat(achievementsFile); err == nil {
		// Always trust mtime for UpdatedAt
		stats.UpdatedAt = fi.ModTime().Unix()

		// NOTE: Does not work reliably on every Linux Distribution.
		// Try to extract a "creation-ish" time
		// if stat, ok := fi.Sys().(*syscall.Stat_t); ok && stat.Ctim.Sec != 0 {
		// 	stats.CreatedAt = stat.Ctim.Sec
		// }
	}

	return stats
}

func GetAchievementsStats(w http.ResponseWriter, r *http.Request) {
	stats := GetAchievementsStatsInternal()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}
