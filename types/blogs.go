package types

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
    Title     string
    Excerpt   string
    Slug      string
    Published string
    Updated   string
    Views     int
}

// SampleBlogPosts is hardcoded placeholder data — swap the handler that
// calls BlogPostsList/BlogInfo over to a real DB-backed source later;
// nothing in these components needs to change when you do.
var SampleBlogPosts = []BlogPost{
    {Title: "Caveats & Workarounds — OpenMediaCloud", Excerpt: "Architectural Workarounds for Compute-Dependent Features", Slug: "caveats-workarounds-openmediacloud", Published: "May 11, 2026 at 01:14 PM", Updated: "May 11, 2026 at 07:47 AM", Views: 13},
    {Title: "Setup Guide — OpenMediaCloud", Excerpt: "This guide walks you through setting up OpenMediaCloud on a Linux server", Slug: "setup-guide-openmediacloud", Published: "May 11, 2026 at 01:02 PM", Updated: "May 11, 2026 at 07:47 AM", Views: 12},
    {Title: "A Proxy that cuts my Jellyfin Hosting costs by 80%", Excerpt: "Say goodbye to egress costs", Slug: "proxy-cuts-jellyfin-hosting-costs", Published: "April 3, 2026 at 08:32 PM", Updated: "May 11, 2026 at 07:47 AM", Views: 55},
    {Title: "Safari Is a Dumpster Fire and Apple Knows It", Excerpt: "Safari is not \"privacy-focused brilliance.\" — It's a broken, inconsistent, developer-hostile piece of shit.", Slug: "safari-is-a-dumpster-fire", Published: "February 13, 2026 at 02:33 PM", Updated: "", Views: 18},
    {Title: "Convert Folders to CBZ Archives with \"convert-cbz\" — Fast, Cross-Platform Comic Converter", Excerpt: "A powerful, cross-platform tool for converting directory of folders of images into CBZ comic book archives. \"convert-cbz\" helps you to quickly create clean, organized CBZ files for your entire comic or manga library.", Slug: "convert-cbz", Published: "September 22, 2025 at 12:59 PM", Updated: "November 16, 2025 at 08:59 AM", Views: 601},
    {Title: "Email Done Right: My End-to-End Journey from Gmail to AWS SES — Part 3", Excerpt: "Getting SES production access, cleaning up configs, and finally sending emails that look truly professional.", Slug: "email-done-right-part-3", Published: "August 23, 2025 at 08:09 PM", Updated: "August 23, 2025 at 02:51 PM", Views: 24},
    {Title: "Email Done Right: My End-to-End Journey from Gmail to AWS SES — Part 2", Excerpt: "Setting up AWS SES, DNS records, and testing email flow while dealing with sandbox restrictions.", Slug: "email-done-right-part-2", Published: "August 23, 2025 at 08:08 PM", Updated: "August 23, 2025 at 02:51 PM", Views: 22},
    {Title: "Email Done Right: My End-to-End Journey from Gmail to AWS SES — Part 1", Excerpt: "Why Gmail's \"sent via\" headers and limits pushed me to find a professional email-sending solution.", Slug: "email-done-right-part-1", Published: "August 23, 2025 at 08:04 PM", Updated: "August 23, 2025 at 02:51 PM", Views: 30},
    {Title: "Hello World — I'm Jelius", Excerpt: "A quick introduction to who I am, what I do, and why I build. Meet the mind behind jelius.dev.", Slug: "hello-world-im-jelius", Published: "July 21, 2025 at 08:24 PM", Updated: "", Views: 73},
}

