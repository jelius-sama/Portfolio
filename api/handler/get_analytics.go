package handler

import (
	"KazuFolio/db"
	"KazuFolio/types"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func GetAnalytics(w http.ResponseWriter, r *http.Request) {
	databaseConn := db.Conn

	// Parse query params
	timeRange := r.URL.Query().Get("time_range")
	eventType := r.URL.Query().Get("event_type")

	// Map time range to SQL interval
	var since time.Time
	now := time.Now()
	switch timeRange {
	case "24h":
		since = now.Add(-24 * time.Hour)
	case "7d", "":
		since = now.AddDate(0, 0, -7)
	case "30d":
		since = now.AddDate(0, 0, -30)
	case "90d":
		since = now.AddDate(0, 0, -90)
	default:
		http.Error(w, "Invalid time_range", http.StatusBadRequest)
		return
	}

	// Base query
	query := `
		SELECT
			session_id, event_type, event_timestamp, page_url, referrer_url,
			ip_address, country, region, city, user_agent, device_type,
			browser_name, browser_version, os_name, os_version,
			screen_width, screen_height, viewport_width, viewport_height,
			language, utm_source, utm_medium, utm_campaign, utm_term,
			utm_content, page_load_time_ms, time_on_page_sec,
			scroll_depth_pct, element_id, error_message, metadata
		FROM analytics
		WHERE event_timestamp >= ?
	`
	args := []any{since}

	// Optional event type filter
	if eventType != "" && eventType != "all" {
		query += " AND event_type = ?"
		args = append(args, eventType)
	}

	rows, err := databaseConn.Query(query, args...)
	if err != nil {
		http.Error(w, "Failed to query analytics data", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var events []types.AnalyticsEvent

	for rows.Next() {
		var evt types.AnalyticsEvent
		var metadataJSON sql.NullString
		var eventTime time.Time

		// Use sql.NullXXX for nullable columns
		var (
			referrerURL    sql.NullString
			ipAddress      sql.NullString
			country        sql.NullString
			region         sql.NullString
			city           sql.NullString
			userAgent      sql.NullString
			deviceType     sql.NullString
			browserName    sql.NullString
			browserVersion sql.NullString
			osName         sql.NullString
			osVersion      sql.NullString
			language       sql.NullString
			utmSource      sql.NullString
			utmMedium      sql.NullString
			utmCampaign    sql.NullString
			utmTerm        sql.NullString
			utmContent     sql.NullString
			elementID      sql.NullString
			errorMessage   sql.NullString
			screenWidth    sql.NullInt64
			screenHeight   sql.NullInt64
			viewportWidth  sql.NullInt64
			viewportHeight sql.NullInt64
			pageLoadTimeMs sql.NullInt64
			timeOnPageSec  sql.NullFloat64
			scrollDepthPct sql.NullFloat64
		)

		err := rows.Scan(
			&evt.SessionID,
			&evt.EventType,
			&eventTime,
			&evt.PageURL,
			&referrerURL,
			&ipAddress,
			&country,
			&region,
			&city,
			&userAgent,
			&deviceType,
			&browserName,
			&browserVersion,
			&osName,
			&osVersion,
			&screenWidth,
			&screenHeight,
			&viewportWidth,
			&viewportHeight,
			&language,
			&utmSource,
			&utmMedium,
			&utmCampaign,
			&utmTerm,
			&utmContent,
			&pageLoadTimeMs,
			&timeOnPageSec,
			&scrollDepthPct,
			&elementID,
			&errorMessage,
			&metadataJSON,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		evt.EventTimestamp = &eventTime

		// Assign optional fields if valid
		if referrerURL.Valid {
			evt.ReferrerURL = referrerURL.String
		}
		if ipAddress.Valid {
			evt.IPAddress = ipAddress.String
		}
		if country.Valid {
			evt.Country = country.String
		}
		if region.Valid {
			evt.Region = region.String
		}
		if city.Valid {
			evt.City = city.String
		}
		if userAgent.Valid {
			evt.UserAgent = userAgent.String
		}
		if deviceType.Valid {
			evt.DeviceType = deviceType.String
		}
		if browserName.Valid {
			evt.BrowserName = browserName.String
		}
		if browserVersion.Valid {
			evt.BrowserVersion = browserVersion.String
		}
		if osName.Valid {
			evt.OSName = osName.String
		}
		if osVersion.Valid {
			evt.OSVersion = osVersion.String
		}
		if language.Valid {
			evt.Language = language.String
		}
		if utmSource.Valid {
			evt.UTMSource = utmSource.String
		}
		if utmMedium.Valid {
			evt.UTMMedium = utmMedium.String
		}
		if utmCampaign.Valid {
			evt.UTMCampaign = utmCampaign.String
		}
		if utmTerm.Valid {
			evt.UTMTerm = utmTerm.String
		}
		if utmContent.Valid {
			evt.UTMContent = utmContent.String
		}
		if elementID.Valid {
			evt.ElementID = elementID.String
		}
		if errorMessage.Valid {
			evt.ErrorMessage = errorMessage.String
		}
		if screenWidth.Valid {
			val := int(screenWidth.Int64)
			evt.ScreenWidth = &val
		}
		if screenHeight.Valid {
			val := int(screenHeight.Int64)
			evt.ScreenHeight = &val
		}
		if viewportWidth.Valid {
			val := int(viewportWidth.Int64)
			evt.ViewportWidth = &val
		}
		if viewportHeight.Valid {
			val := int(viewportHeight.Int64)
			evt.ViewportHeight = &val
		}
		if pageLoadTimeMs.Valid {
			val := int(pageLoadTimeMs.Int64)
			evt.PageLoadTimeMs = &val
		}
		if timeOnPageSec.Valid {
			val := float64(timeOnPageSec.Float64)
			evt.TimeOnPageSec = &val
		}
		if scrollDepthPct.Valid {
			val := float64(scrollDepthPct.Float64)
			evt.ScrollDepthPct = &val
		}
		if metadataJSON.Valid {
			var metadata map[string]any
			if err := json.Unmarshal([]byte(metadataJSON.String), &metadata); err == nil {
				evt.Metadata = metadata
			}
		}

		events = append(events, evt)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Println("Encoding error:", err)
	}
}
