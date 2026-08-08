package db

func createAnalyticsTables() error {
    var schema = `
    CREATE TABLE IF NOT EXISTS analytics_events (
        event_id INTEGER PRIMARY KEY AUTOINCREMENT,
        country_code TEXT NOT NULL,
        page_path TEXT NOT NULL,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX IF NOT EXISTS idx_country_code ON analytics_events(country_code);
    CREATE INDEX IF NOT EXISTS idx_page_path ON analytics_events(page_path);
    CREATE INDEX IF NOT EXISTS idx_timestamp ON analytics_events(timestamp);
    `

    if _, err := DB.Exec(schema); err != nil {
        return err
    }

    return nil
}

