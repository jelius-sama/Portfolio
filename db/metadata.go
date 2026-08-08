package db

func createMetadataTable() error {
    var schema = `
        CREATE TABLE IF NOT EXISTS metadata (
            path TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            description TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS m_links (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            metadata_path TEXT NOT NULL,
            rel TEXT NOT NULL,
            href TEXT NOT NULL,
            media TEXT,
            FOREIGN KEY (metadata_path) REFERENCES metadata(path) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS m_meta (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            metadata_path TEXT NOT NULL,
            name TEXT,
            property TEXT,
            content TEXT NOT NULL,
            FOREIGN KEY (metadata_path) REFERENCES metadata(path) ON DELETE CASCADE
        );

        CREATE INDEX IF NOT EXISTS idx_m_links_metadata_path ON m_links(metadata_path);
        CREATE INDEX IF NOT EXISTS idx_m_meta_metadata_path ON m_meta(metadata_path);
    `

    if _, err := DB.Exec(schema); err != nil {
        return err
    }

    return nil
}

