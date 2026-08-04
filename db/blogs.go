package db

// createBlogsTable creates the blogs table with support for posts, relationships, and soft deletes
func createBlogsTable() error {
    schema := `
    CREATE TABLE IF NOT EXISTS blogs (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        description TEXT,
        published_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        deleted_at DATETIME,
        prequel_id INTEGER,
        sequel_id INTEGER,
        FOREIGN KEY (prequel_id) REFERENCES blogs(id) ON DELETE SET NULL,
        FOREIGN KEY (sequel_id) REFERENCES blogs(id) ON DELETE SET NULL
    );

    CREATE INDEX IF NOT EXISTS idx_published_at ON blogs(published_at);
    CREATE INDEX IF NOT EXISTS idx_deleted_at ON blogs(deleted_at);
    CREATE INDEX IF NOT EXISTS idx_prequel_id ON blogs(prequel_id);
    CREATE INDEX IF NOT EXISTS idx_sequel_id ON blogs(sequel_id);
    CREATE INDEX IF NOT EXISTS idx_title ON blogs(title);
    `

    _, err := DB.Exec(schema)
    if err != nil {
        return err
    }

    return nil
}

