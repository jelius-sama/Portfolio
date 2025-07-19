package types

import "time"

type AnalyticsEvent struct {
	SessionID      string         `json:"session_id"`
	EventType      string         `json:"event_type"`
	EventTimestamp *time.Time     `json:"event_timestamp,omitempty"`
	PageURL        string         `json:"page_url"`
	ReferrerURL    string         `json:"referrer_url,omitempty"`
	IPAddress      string         `json:"ip_address,omitempty"`
	Country        string         `json:"country,omitempty"`
	Region         string         `json:"region,omitempty"`
	City           string         `json:"city,omitempty"`
	UserAgent      string         `json:"user_agent,omitempty"`
	DeviceType     string         `json:"device_type,omitempty"`
	BrowserName    string         `json:"browser_name,omitempty"`
	BrowserVersion string         `json:"browser_version,omitempty"`
	OSName         string         `json:"os_name,omitempty"`
	OSVersion      string         `json:"os_version,omitempty"`
	ScreenWidth    *int           `json:"screen_width,omitempty"`
	ScreenHeight   *int           `json:"screen_height,omitempty"`
	ViewportWidth  *int           `json:"viewport_width,omitempty"`
	ViewportHeight *int           `json:"viewport_height,omitempty"`
	Language       string         `json:"language,omitempty"`
	UTMSource      string         `json:"utm_source,omitempty"`
	UTMMedium      string         `json:"utm_medium,omitempty"`
	UTMCampaign    string         `json:"utm_campaign,omitempty"`
	UTMTerm        string         `json:"utm_term,omitempty"`
	UTMContent     string         `json:"utm_content,omitempty"`
	PageLoadTimeMs *int           `json:"page_load_time_ms,omitempty"`
	TimeOnPageSec  *float64       `json:"time_on_page_sec,omitempty"`
	ScrollDepthPct *float64       `json:"scroll_depth_pct,omitempty"`
	ElementID      string         `json:"element_id,omitempty"`
	ErrorMessage   string         `json:"error_message,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}
