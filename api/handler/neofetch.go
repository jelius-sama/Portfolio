package handler

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type InfoItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

var (
	neofetchCacheData     []InfoItem
	neofetchCacheTime     time.Time
	neofetchCacheDuration = 12 * time.Hour
	neofetchCacheMutex    sync.Mutex
)

func NeofetchInfo(w http.ResponseWriter, r *http.Request) {
	neofetchCacheMutex.Lock()
	defer neofetchCacheMutex.Unlock()

	// Serve from cache if still valid
	if time.Since(neofetchCacheTime) < neofetchCacheDuration && neofetchCacheData != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(neofetchCacheData)
		return
	}

	cmd := exec.Command("neofetch", "--stdout")
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to run neofetch", http.StatusInternalServerError)
		return
	}

	var infoList []InfoItem
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, ":") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key != "" && value != "" {
			infoList = append(infoList, InfoItem{
				Label: key,
				Value: value,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "Failed to parse neofetch output", http.StatusInternalServerError)
		return
	}

	// Cache the result
	neofetchCacheData = infoList
	neofetchCacheTime = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(infoList)
}
