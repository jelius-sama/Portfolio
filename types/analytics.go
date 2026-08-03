package types

type TopCountryResponse struct {
    CountryCode string  `json:"country_code"`
    VisitCount  int     `json:"visit_count"`
    Percentage  float64 `json:"percentage"`
}

type PaginatedTopCountriesResponse struct {
    Data      []TopCountryResponse `json:"data"`
    Page      int                  `json:"page"`
    Limit     int                  `json:"limit"`
    HasMore   bool                 `json:"has_more"`
    TotalRows int                  `json:"total_rows"`
}

type AvgVisitsResponse struct {
    TimeWindow  string  `json:"time_window"`
    CountryCode string  `json:"country_code"`
    PagePath    string  `json:"page_path"`
    AvgVisits   float64 `json:"avg_visits"`
    Date        string  `json:"date"`
}

type PaginatedAvgVisitsResponse struct {
    Data      []AvgVisitsResponse `json:"data"`
    Page      int                 `json:"page"`
    Limit     int                 `json:"limit"`
    HasMore   bool                `json:"has_more"`
    TotalRows int                 `json:"total_rows"`
}

type AnalyticsEventResponse struct {
    EventID      int64  `json:"event_id"`
    CountryCode  string `json:"country_code"`
    PagePath     string `json:"page_path"`
    TimestampUTC string `json:"timestamp"`
}

type PaginatedEventsResponse struct {
    Data      []AnalyticsEventResponse `json:"data"`
    Page      int                      `json:"page"`
    Limit     int                      `json:"limit"`
    HasMore   bool                     `json:"has_more"`
    TotalRows int                      `json:"total_rows"`
}

type TopPageResponse struct {
    PagePath   string  `json:"page_path"`
    VisitCount int     `json:"visit_count"`
    Percentage float64 `json:"percentage"`
}

type PaginatedTopPagesResponse struct {
    Data      []TopPageResponse `json:"data"`
    Page      int               `json:"page"`
    Limit     int               `json:"limit"`
    HasMore   bool              `json:"has_more"`
    TotalRows int               `json:"total_rows"`
}

type TrackAnalyticsRequest struct {
    CountryCode string `json:"country_code"`
    PagePath    string `json:"page_path"`
}

type SortOrder uint8

const (
    SOAsc SortOrder = iota
    SODesc
)

