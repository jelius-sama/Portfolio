package db

const CreateBlogsTable = `
CREATE TABLE IF NOT EXISTS blogs (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT,

    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),

    prequel_id TEXT,
    sequel_id TEXT,

    parts TEXT, -- JSON-encoded list of blog IDs

    FOREIGN KEY (prequel_id) REFERENCES blogs(id),
    FOREIGN KEY (sequel_id) REFERENCES blogs(id)
);
`
