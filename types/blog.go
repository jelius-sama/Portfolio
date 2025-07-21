package types

type Blog struct {
	ID      string `json:"id"`      // e.g., "52a081c"
	Title   string `json:"title"`   // Blog post title
	Summary string `json:"summary"` // Short excerpt or summary

	CreatedAt string `json:"createdAt"` // ISO 8601 format (or use time.Time)
	UpdatedAt string `json:"updatedAt"` // Optional

	PrequelID *string `json:"prequelId"` // Blog ID of prequel (nullable)
	SequelID  *string `json:"sequelId"`  // Blog ID of sequel (nullable)

	Parts []string `json:"parts"` // List of Blog IDs if this is part of a multipart series
}
