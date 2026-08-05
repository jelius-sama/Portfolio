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

type CreateBlogPost struct {
    Title     string  `json:"title"`
    Excerpt   string  `json:"excerpt"`
    PrequelID *string `json:"prequel_id"`
    SequelID  *string `json:"sequel_id"`
}

type BlogResponse struct {
    ID          string        `json:"id"`
    PublishedAt time.Time     `json:"published_at"`
    UpdatedAt   time.Time     `json:"updated_at"`
    DeletedAt   *time.Time    `json:"-"`
    Prequel     *BlogResponse `json:"prequel"`
    Sequel      *BlogResponse `json:"sequel"`
    Title       string        `json:"title"`
    Excerpt     string        `json:"excerpt"`
    Views       uint          `json:"views"`
}

type PaginatedBlogsResponse struct {
    Data      []BlogPost `json:"data"`
    Page      int        `json:"page"`
    Limit     int        `json:"limit"`
    HasMore   bool       `json:"has_more"`
    TotalRows int        `json:"total_rows"`
    Sort      string     `json:"sort"`
}

// SampleBlogPosts is hardcoded placeholder data — swap the handler that
// calls BlogPostsList/BlogInfo over to a real DB-backed source later;
// nothing in these components needs to change when you do.
var SampleBlogPosts = []BlogPost{
    {
        ID:          "caveats-workarounds-openmediacloud",
        Title:       "Caveats & Workarounds — OpenMediaCloud",
        Excerpt:     "Architectural Workarounds for Compute-Dependent Features",
        PublishedAt: time.Date(2026, 5, 11, 13, 14, 0, 0, time.UTC),
        UpdatedAt:   time.Date(2026, 5, 11, 7, 47, 0, 0, time.UTC),
        Views:       13,
    },
    {
        ID:          "setup-guide-openmediacloud",
        Title:       "Setup Guide — OpenMediaCloud",
        Excerpt:     "This guide walks you through setting up OpenMediaCloud on a Linux server",
        PublishedAt: time.Date(2026, 5, 11, 13, 2, 0, 0, time.UTC),
        UpdatedAt:   time.Date(2026, 5, 11, 7, 47, 0, 0, time.UTC),
        Views:       12,
    },
    {
        ID:          "proxy-cuts-jellyfin-hosting-costs",
        Title:       "A Proxy that cuts my Jellyfin Hosting costs by 80%",
        Excerpt:     "Say goodbye to egress costs",
        PublishedAt: time.Date(2026, 4, 3, 20, 32, 0, 0, time.UTC),
        UpdatedAt:   time.Date(2026, 5, 11, 7, 47, 0, 0, time.UTC),
        Views:       55,
    },
    {
        ID:          "safari-is-a-dumpster-fire",
        Title:       "Safari Is a Dumpster Fire and Apple Knows It",
        Excerpt:     "Safari is not \"privacy-focused brilliance.\" — It's a broken, inconsistent, developer-hostile piece of shit.",
        PublishedAt: time.Date(2026, 2, 13, 14, 33, 0, 0, time.UTC),
        Views:       18,
    },
    {
        ID:          "convert-cbz",
        Title:       "Convert Folders to CBZ Archives with \"convert-cbz\" — Fast, Cross-Platform Comic Converter",
        Excerpt:     "A powerful, cross-platform tool for converting directory of folders of images into CBZ comic book archives. \"convert-cbz\" helps you to quickly create clean, organized CBZ files for your entire comic or manga library.",
        PublishedAt: time.Date(2025, 9, 22, 12, 59, 0, 0, time.UTC),
        UpdatedAt:   time.Date(2025, 11, 16, 8, 59, 0, 0, time.UTC),
        Views:       601,
    },
    {
        ID:          "email-done-right-part-3",
        Title:       "Email Done Right: My End-to-End Journey from Gmail to AWS SES — Part 3",
        Excerpt:     "Getting SES production access, cleaning up configs, and finally sending emails that look truly professional.",
        PublishedAt: time.Date(2025, 8, 23, 20, 9, 0, 0, time.UTC),
        UpdatedAt:   time.Date(2025, 8, 23, 14, 51, 0, 0, time.UTC),
        Views:       24,
    },
    {
        ID:          "email-done-right-part-2",
        Title:       "Email Done Right: My End-to-End Journey from Gmail to AWS SES — Part 2",
        Excerpt:     "Setting up AWS SES, DNS records, and testing email flow while dealing with sandbox restrictions.",
        PublishedAt: time.Date(2025, 8, 23, 20, 8, 0, 0, time.UTC),
        UpdatedAt:   time.Date(2025, 8, 23, 14, 51, 0, 0, time.UTC),
        Views:       22,
    },
    {
        ID:          "email-done-right-part-1",
        Title:       "Email Done Right: My End-to-End Journey from Gmail to AWS SES — Part 1",
        Excerpt:     "Why Gmail's \"sent via\" headers and limits pushed me to find a professional email-sending solution.",
        PublishedAt: time.Date(2025, 8, 23, 20, 4, 0, 0, time.UTC),
        UpdatedAt:   time.Date(2025, 8, 23, 14, 51, 0, 0, time.UTC),
        Views:       30,
    },
    {
        ID:          "hello-world-im-jelius",
        Title:       "Hello World — I'm Jelius",
        Excerpt:     "A quick introduction to who I am, what I do, and why I build. Meet the mind behind jelius.dev.",
        PublishedAt: time.Date(2025, 7, 21, 20, 24, 0, 0, time.UTC),
        Views:       73,
    },
}

