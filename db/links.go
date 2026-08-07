package db

func createLinksTable() error {
    var schema = `
        CREATE TABLE IF NOT EXISTS links_page (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            handle TEXT NOT NULL,
            tag_line TEXT NOT NULL,
            image_path TEXT NOT NULL,
            image_alt_text TEXT,
            who_am_i TEXT NOT NULL,
            qr_image_path TEXT NOT NULL,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS link_entries (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            links_page_id INTEGER NOT NULL,
            icon TEXT NOT NULL,
            title TEXT NOT NULL,
            subtitle TEXT NOT NULL,
            href TEXT NOT NULL,
            position INTEGER NOT NULL,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (links_page_id) REFERENCES links_page(id) ON DELETE CASCADE
        );

        CREATE INDEX IF NOT EXISTS idx_links_page_id ON link_entries(links_page_id);
    `

    if _, err := DB.Exec(schema); err != nil {
        return err
    }

    return nil
}

