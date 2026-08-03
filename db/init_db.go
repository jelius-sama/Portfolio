package db

import (
    "database/sql"
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
    return errors
}

// InitDB initializes the SQLite3 database connection and creates the schema
func InitDB(dbPath string) error {
    mu.Lock()
    defer mu.Unlock()

    var err error
    DB, err = sql.Open("sqlite3", dbPath)
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

