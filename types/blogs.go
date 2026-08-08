package types

import "time"

type BlogsSortOrder uint8

const PostPerPage = 5

const (
    BSONew BlogsSortOrder = iota
    BSOOld
    BSOPopular
)

func (bso BlogsSortOrder) String() string {
    switch bso {
    case BSONew:
        return "Newest"
    case BSOOld:
        return "Oldest"
    case BSOPopular:
        return "Most Viewed"
    default:
        return ""
    }
}

// BlogPost is one entry in the blog list.
type BlogPost struct {
    ID          string     `json:"id"`
    PublishedAt time.Time  `json:"published_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"-"`
    PrequelID   *string    `json:"prequel_id"`
    SequelID    *string    `json:"sequel_id"`
    Title       string     `json:"title"`
    Excerpt     string     `json:"excerpt"`
    Views       uint       `json:"views"`
}

type BlogResponse struct {
    ID          string     `json:"id"`
    PublishedAt time.Time  `json:"published_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"-"`
    // TODO: Implement the fixme.
    Prequel *BlogResponse `json:"prequel"` // FIXME: Maybe a better solution would be to use `BlogPost` here
    Sequel  *BlogResponse `json:"sequel"`  // FIXME: Maybe a better solution would be to use `BlogPost` here
    Title   string        `json:"title"`
    Excerpt string        `json:"excerpt"`
    Views   uint          `json:"views"`
}

type CreateBlogPost struct {
    Title     string  `json:"title"`
    Excerpt   string  `json:"excerpt"`
    PrequelID *string `json:"prequel_id"`
    SequelID  *string `json:"sequel_id"`
}

type PaginatedBlogsResponse struct {
    Data      []BlogPost     `json:"data"`
    Page      int            `json:"page"`
    Limit     int            `json:"limit"`
    HasMore   bool           `json:"has_more"`
    TotalRows int            `json:"total_rows"`
    Sort      BlogsSortOrder `json:"sort"`
}

