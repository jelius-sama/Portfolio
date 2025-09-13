package handler

import (
	"database/sql"
	"fmt"

	"KazuFolio/db"
)

// PageViewsInternal returns the number of visits for a given relative path.
func PageViewsInternal(path string) (uint, error) {
	databaseConn := db.Conn

	query := `
		SELECT COUNT(*) FROM analytics
		WHERE page_url LIKE ?
	`
	likePattern := fmt.Sprintf("%%%s%%", path)

	var count uint
	err := databaseConn.QueryRow(query, likePattern).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return count, nil
}
