package db

const CreateAnalyticsTable = `
CREATE TABLE IF NOT EXISTS analytics (
    id                INTEGER  PRIMARY KEY AUTOINCREMENT,
    session_id        TEXT     NOT NULL,
    event_type        TEXT     NOT NULL,
    event_timestamp   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    page_url          TEXT     NOT NULL,
    referrer_url      TEXT,
    ip_address        TEXT,
    country           TEXT,
    region            TEXT,
    city              TEXT,
    user_agent        TEXT,
    device_type       TEXT,
    browser_name      TEXT,
    browser_version   TEXT,
    os_name           TEXT,
    os_version        TEXT,
    screen_width      INTEGER,
    screen_height     INTEGER,
    viewport_width    INTEGER,
    viewport_height   INTEGER,
    language          TEXT,
    utm_source        TEXT,
    utm_medium        TEXT,
    utm_campaign      TEXT,
    utm_term          TEXT,
    utm_content       TEXT,
    page_load_time_ms INTEGER,
    time_on_page_sec  REAL,
    scroll_depth_pct  REAL,
    element_id        TEXT,
    error_message     TEXT,
    metadata          JSON
);
`
