package handler

import (
	"KazuFolio/db"
	"KazuFolio/types"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"KazuFolio/logger"
)

func SaveAnalytics(w http.ResponseWriter, r *http.Request) {
	var evt types.AnalyticsEvent
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	if evt.IPAddress == "" {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			evt.IPAddress = strings.TrimSpace(parts[0])
		} else if xr := r.Header.Get("X-Real-IP"); xr != "" {
			evt.IPAddress = xr
		} else {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				evt.IPAddress = r.RemoteAddr
			} else {
				evt.IPAddress = host
			}
		}
	}

	if evt.EventTimestamp == nil {
		now := time.Now().UTC()
		evt.EventTimestamp = &now
	}

	query := `
INSERT INTO analytics (
    session_id, event_type, event_timestamp,
    page_url, referrer_url,
    ip_address, country, region, city,
    user_agent, device_type, browser_name, browser_version,
    os_name, os_version, screen_width, screen_height,
    viewport_width, viewport_height, language,
    utm_source, utm_medium, utm_campaign, utm_term, utm_content,
    page_load_time_ms, time_on_page_sec, scroll_depth_pct,
    element_id, error_message, metadata
) VALUES (
    ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
`

	args := []any{
		evt.SessionID,
		evt.EventType,
		*evt.EventTimestamp,
		evt.PageURL,
		nullString(evt.ReferrerURL),
		nullString(evt.IPAddress),
		nullString(evt.Country),
		nullString(evt.Region),
		nullString(evt.City),
		nullString(evt.UserAgent),
		nullString(evt.DeviceType),
		nullString(evt.BrowserName),
		nullString(evt.BrowserVersion),
		nullString(evt.OSName),
		nullString(evt.OSVersion),
		evt.ScreenWidth,
		evt.ScreenHeight,
		evt.ViewportWidth,
		evt.ViewportHeight,
		nullString(evt.Language),
		nullString(evt.UTMSource),
		nullString(evt.UTMMedium),
		nullString(evt.UTMCampaign),
		nullString(evt.UTMTerm),
		nullString(evt.UTMContent),
		evt.PageLoadTimeMs,
		evt.TimeOnPageSec,
		evt.ScrollDepthPct,
		nullString(evt.ElementID),
		nullString(evt.ErrorMessage),
		toJSON(evt.Metadata),
	}

	if _, err := db.Conn.Exec(query, args...); err != nil {
		logger.Error("Failed to save analytics:", err)
		http.Error(w, "failed to save analytics", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func toJSON(m map[string]any) any {
	if len(m) == 0 {
		return nil
	}
	b, _ := json.Marshal(m)
	return string(b)
}
