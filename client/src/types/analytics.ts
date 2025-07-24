export interface AnalyticsRecord {
    id: number
    session_id: string
    event_type: string
    event_timestamp: string
    page_url: string
    referrer_url?: string
    ip_address?: string
    country?: string
    region?: string
    city?: string
    user_agent?: string
    device_type?: string
    browser_name?: string
    browser_version?: string
    os_name?: string
    os_version?: string
    screen_width?: number
    screen_height?: number
    viewport_width?: number
    viewport_height?: number
    language?: string
    utm_source?: string
    utm_medium?: string
    utm_campaign?: string
    utm_term?: string
    utm_content?: string
    page_load_time_ms?: number
    time_on_page_sec?: number
    scroll_depth_pct?: number
    element_id?: string
    error_message?: string
    metadata?: any
}

export interface AnalyticsSummary {
    totalPageViews: number
    uniqueVisitors: number
    avgTimeOnPage: number
    avgPageLoadTime: number
    topPages: Array<{ page: string; views: number }>
    topCountries: Array<{ country: string; visitors: number }>
    deviceTypes: Array<{ type: string; count: number }>
    browsers: Array<{ browser: string; count: number }>
    recentEvents: AnalyticsRecord[]
}
