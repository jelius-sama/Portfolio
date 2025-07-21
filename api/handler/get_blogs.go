package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"KazuFolio/db"
	"KazuFolio/types"
)

type BlogsResponse struct {
	Blogs      []types.Blog `json:"blogs"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	TotalPages int          `json:"totalPages"`
	Limit      int          `json:"limit"`
}

type CacheEntry struct {
	Data      BlogsResponse
	ExpiresAt time.Time
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
	// Track ongoing requests to prevent duplicate DB calls
	ongoing map[string]*sync.WaitGroup
}

var cache = &Cache{
	entries: make(map[string]CacheEntry),
	ongoing: make(map[string]*sync.WaitGroup),
}

const CACHE_DURATION = 10 * time.Minute

// generateCacheKey creates a unique key for the request parameters
func generateCacheKey(page, limit int, sortBy, order string) string {
	key := fmt.Sprintf("page:%d_limit:%d_sortBy:%s_order:%s", page, limit, sortBy, order)
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)
}

// cleanExpiredEntries removes expired cache entries
func (c *Cache) cleanExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.ExpiresAt) {
			delete(c.entries, key)
		}
	}
}

// getCachedResponse tries to get a cached response
func (c *Cache) getCachedResponse(key string) (BlogsResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return BlogsResponse{}, false
	}

	return entry.Data, true
}

// setCachedResponse stores a response in cache
func (c *Cache) setCachedResponse(key string, data BlogsResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(CACHE_DURATION),
	}
}

// getOrCreateWaitGroup manages ongoing requests to prevent duplicate DB calls
func (c *Cache) getOrCreateWaitGroup(key string) (*sync.WaitGroup, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if wg, exists := c.ongoing[key]; exists {
		return wg, false // Request already in progress
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	c.ongoing[key] = wg
	return wg, true // This is the first request for this key
}

// finishRequest cleans up the ongoing request tracking
func (c *Cache) finishRequest(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if wg, exists := c.ongoing[key]; exists {
		delete(c.ongoing, key)
		wg.Done()
	}
}

func fetchBlogsFromDB(page, limit int, sortBy, order string) (BlogsResponse, error) {
	// Get total count for pagination
	countQuery := "SELECT COUNT(*) FROM blogs"
	var total int
	if err := db.Conn.QueryRow(countQuery).Scan(&total); err != nil {
		return BlogsResponse{}, fmt.Errorf("failed to get total count: %w", err)
	}

	// Calculate pagination values
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	offset := (page - 1) * limit

	// Get blogs with pagination
	query := `
		SELECT id, title, summary, created_at, updated_at, prequel_id, sequel_id, parts
		FROM blogs
		ORDER BY ` + sortBy + ` ` + order + `
		LIMIT ? OFFSET ?`

	rows, err := db.Conn.Query(query, limit, offset)
	if err != nil {
		return BlogsResponse{}, fmt.Errorf("database query failed: %w", err)
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
			&blog.CreatedAt,
			&blog.UpdatedAt,
			&blog.PrequelID,
			&blog.SequelID,
			&parts); err != nil {
			return BlogsResponse{}, fmt.Errorf("failed to scan blog: %w", err)
		}

		if err := json.Unmarshal([]byte(parts), &blog.Parts); err != nil {
			return BlogsResponse{}, fmt.Errorf("failed to parse parts: %w", err)
		}

		blogs = append(blogs, blog)
	}

	if err := rows.Err(); err != nil {
		return BlogsResponse{}, fmt.Errorf("row iteration failed: %w", err)
	}

	return BlogsResponse{
		Blogs:      blogs,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
		Limit:      limit,
	}, nil
}

func GetBlogs(w http.ResponseWriter, r *http.Request) {
	// Clean expired cache entries periodically
	go cache.cleanExpiredEntries()

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

	// Generate cache key
	cacheKey := generateCacheKey(page, limit, sortBy, order)

	// Try to get from cache first
	if cachedResponse, found := cache.getCachedResponse(cacheKey); found {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(cachedResponse)
		return
	}

	// Check if request is already in progress
	wg, isFirst := cache.getOrCreateWaitGroup(cacheKey)

	if !isFirst {
		// Another request is already fetching this data, wait for it
		wg.Wait()
		// Try cache again after waiting
		if cachedResponse, found := cache.getCachedResponse(cacheKey); found {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT-AFTER-WAIT")
			json.NewEncoder(w).Encode(cachedResponse)
			return
		}
	}

	// This is the first request or cache miss after wait, fetch from DB
	response, err := fetchBlogsFromDB(page, limit, sortBy, order)

	// Clean up ongoing request tracking
	cache.finishRequest(cacheKey)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache the response
	cache.setCachedResponse(cacheKey, response)

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
