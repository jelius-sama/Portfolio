package db

import (
    "database/sql"
    "fmt"
    "sync"

    "github.com/jelius-sama/logger"
    _ "modernc.org/sqlite"
)

var (
    DB *sql.DB
    mu sync.Mutex
)

func createTables() []error {
    var errors []error
    errors = append(errors, createAnalyticsTables())
    errors = append(errors, createBlogsTable())
    errors = append(errors, createLinksTable())
    errors = append(errors, createHomeTables())
    errors = append(errors, createMetadataTable())
    return errors
}

// applyPragmas sets connection-level pragmas that determine SQLite's
// concurrency behavior. These must be set on every connection in the
// pool (sql.DB doesn't guarantee pragma state survives across pooled
// conns for all drivers), so they're issued once here right after open
// and modernc.org/sqlite persists WAL mode at the file level anyway.
func applyPragmas(db *sql.DB) error {
    pragmas := []string{
        "PRAGMA journal_mode = WAL;",   // readers don't block writers, writers don't block readers
        "PRAGMA busy_timeout = 5000;",  // retry up to 5s on lock contention instead of erroring immediately
        "PRAGMA synchronous = NORMAL;", // safe with WAL, much cheaper than FULL
        "PRAGMA foreign_keys = ON;",    // off by default per-connection in SQLite
        "PRAGMA cache_size = -20000;",  // ~20MB page cache (negative = KB, not pages)
        "PRAGMA temp_store = MEMORY;",  // keep temp b-trees/sorts off disk
    }

    for _, p := range pragmas {
        if _, err := db.Exec(p); err != nil {
            return fmt.Errorf("failed to apply pragma %q: %w", p, err)
        }
    }
    return nil
}

// InitDB initializes the SQLite3 database connection and creates the schema
func InitDB(dbPath string) error {
    mu.Lock()
    defer mu.Unlock()

    var err error
    DB, err = sql.Open("sqlite", dbPath)
    if err != nil {
        logger.Fatal("Failed to open database:", err.Error())
        return err
    }

    // Test the connection
    err = DB.Ping()
    if err != nil {
        logger.Fatal("Failed to ping database:", err.Error())
        return err
    }

    if err := applyPragmas(DB); err != nil {
        logger.Fatal("Failed to apply pragmas:", err.Error())
        return err
    }

    // WAL still serializes writers, so an unbounded pool just means more
    // goroutines queued on busy_timeout rather than more real throughput.
    // Cap it modestly; readers benefit from a few concurrent conns, writes
    // don't benefit from more than one at a time regardless of pool size.
    DB.SetMaxOpenConns(8)
    DB.SetMaxIdleConns(8)

    logger.Info("Database connection established")

    // Create all tables if they don't exist
    var errs = createTables()
    for i := range errs {
        if errs[i] != nil {
            logger.Error("Failed to create tables:", errs[i].Error())
            if i == len(errs) {
                logger.Fatal("Exiting due to one or more errors.")
            }
        }
    }

    return nil
}

// CloseDB closes the database connection
func CloseDB() error {
    if DB != nil {
        return DB.Close()
    }
    return nil
}

