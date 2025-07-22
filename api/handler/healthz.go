package handler

import (
	"KazuFolio/db"
	"KazuFolio/util"
	"net/http"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type ComponentStatus string

const (
	StatusOK        ComponentStatus = "ok"
	StatusLocked    ComponentStatus = "locked"
	StatusHeavyLoad ComponentStatus = "heavy_load"
	StatusDegraded  ComponentStatus = "degraded"
)

type HealthReport struct {
	Status     ComponentStatus            `json:"status"`
	Timestamp  time.Time                  `json:"timestamp"`
	Components map[string]ComponentStatus `json:"components"`
}

var (
	healthCacheMutex sync.RWMutex
	lastHealth       *HealthReport
	processingLock   sync.Mutex
)

// Healthz checks the health of the application and responds with the report
func Healthz(w http.ResponseWriter, r *http.Request) {
	// Only allow one request to generate a new health check
	if !processingLock.TryLock() {
		// If already being processed, serve cached
		serveCachedHealth(w)
		return
	}
	defer processingLock.Unlock()

	report := &HealthReport{
		Status:     StatusOK,
		Timestamp:  time.Now().UTC(),
		Components: make(map[string]ComponentStatus),
	}

	var wg sync.WaitGroup
	var dbStatus ComponentStatus
	var cpuStatus ComponentStatus

	wg.Add(2)

	// Database check
	go func() {
		defer wg.Done()
		err := db.Conn.Ping()
		if err != nil {
			dbStatus = StatusLocked
		} else {
			dbStatus = StatusOK
		}
	}()

	// CPU check
	go func() {
		defer wg.Done()
		percentages, err := cpu.Percent(0, false)
		if err != nil || len(percentages) == 0 {
			cpuStatus = StatusDegraded
			return
		}
		if percentages[0] > 85 {
			cpuStatus = StatusHeavyLoad
		} else {
			cpuStatus = StatusOK
		}
	}()

	wg.Wait()

	report.Components["database"] = dbStatus
	report.Components["load"] = cpuStatus
	report.Components["webserver"] = StatusOK
	report.Components["api"] = StatusOK

	// Overall status decision
	if dbStatus != StatusOK || cpuStatus != StatusOK {
		report.Status = StatusDegraded
	}

	// Cache the latest health report
	healthCacheMutex.Lock()
	lastHealth = report
	healthCacheMutex.Unlock()

	statusCode := http.StatusOK
	if report.Status != StatusOK {
		statusCode = http.StatusServiceUnavailable
	}
	util.WriteJSON(w, statusCode, report)
}

// Serve cached health report with info that it's not current
func serveCachedHealth(w http.ResponseWriter) {
	healthCacheMutex.RLock()
	defer healthCacheMutex.RUnlock()

	if lastHealth == nil {
		http.Error(w, "Health check not ready", http.StatusServiceUnavailable)
		return
	}

	statusCode := http.StatusOK
	if lastHealth.Status != StatusOK {
		statusCode = http.StatusServiceUnavailable
	}
	util.WriteJSON(w, statusCode, lastHealth)
}
