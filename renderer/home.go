package renderer

import (
    "fmt"

    "git.jelius.dev/jelius-sama/Portfolio/template/components"
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/gofiber/fiber/v3"
)

func (v *ViewManager) RenderHome(c fiber.Ctx) error {
    var metadata types.Metadata = types.Metadata{
        Title:       "Home",
        Description: "Portfolio",
    }
    // TODO: Remove hardcoded data to be loaded from a config file
    return Renderer(c, metadata, pages.Home(pages.HomeData{
        HeroSectionData: components.HeroData{
            FirstName:         "Jelius",
            LastName:          "Basumatary",
            EducationLevel:    components.EDLevelUG,
            StudyDomain:       "Computer Science",
            WhoAmI:            "An undergraduate CS student who loves building things that work well and are easy to use.",
            SpecializedDomain: []string{"Backend Developer", "Systems Programmer"},
        },
        AboutSectionMessage: []string{
            "I'm an undergraduate CS student who loves building things that work well and are easy to use. I enjoy solving problems, making code simple, and learning new things every day.",
            "Right now, I'm working toward my bachelor's degree. I enjoy developing software that makes a real difference in everyday life, and I'm always curious how things work under the hood.",
            "When I'm not coding, you'll usually find me exploring new technologies or sharing what I've learned with the developer community.",
        },
        SkillSectionData: components.SkillData{
            Languages:  []string{"Rust", "Swift", "C", "TypeScript", "Go"},
            Frameworks: []string{"Node.js", "SolidJS", "React", "HTMX", "Templ"},
            Others:     []string{"Redis", "PostgreSQL", "AWS", "Cloudflare", "Git"},
        },
        ProjectSectionData: []components.ProjectData{
            {
                Title:       "Pixelle",
                Description: "An anime image gallery application with collections of illustrations and photography.",
                Tags:        []string{"Next.js", "TypeScript", "PostgreSQL"},
                Links: []components.ProjectLink{
                    {Label: "Code", Href: "#", Kind: "code"},
                    {Label: "Live Demo", Href: "#", Kind: "demo"},
                },
                Thumbnail: fmt.Sprintf("%s/assets/compressed/project-pixelle.webp", types.EVAssetCDNHostname.Get().Value),
            },
            {
                Title:       "VPS Watch Dog",
                Description: "A watch dog program that monitors your VPS's system usage and alerts via mail if usage is running high.",
                Tags:        []string{"Go", "SMTP", "POSIX API"},
                Links: []components.ProjectLink{
                    {Label: "Code", Href: "#", Kind: "code"},
                },
                Thumbnail: fmt.Sprintf("%s/assets/compressed/VPSWatchDog.webp", types.EVAssetCDNHostname.Get().Value),
            },
            {
                Title:       "Storage Watch Dog",
                Description: "A program that watches a specific directory and if the storage space is running low alerts via mails.",
                Tags:        []string{"Go", "SMTP"},
                Links: []components.ProjectLink{
                    {Label: "Code", Href: "#", Kind: "code"},
                },
                Thumbnail: fmt.Sprintf("%s/assets/compressed/StorageWatchDog.webp", types.EVAssetCDNHostname.Get().Value),
            },
            {
                Title:       "AWS Mail Parser",
                Description: "Polls your AWS SQS and when a new mail event is detected fetches that mail from S3 bucket and saves it to your maildir after parsing it.",
                Tags:        []string{"Go", "AWS", "File Parsing", "IMAP"},
                Links: []components.ProjectLink{
                    {Label: "Code", Href: "#", Kind: "code"},
                },
                Thumbnail: fmt.Sprintf("%s/assets/compressed/AWSMailParser.webp", types.EVAssetCDNHostname.Get().Value),
            },
            {
                Title:       "Convert CBZ",
                Description: "A high-performance, concurrent tool for converting folders containing images into CBZ (Comic Book Archive) files. Built in Go for speed and reliability.",
                Tags:        []string{"Go", "Archive", "Concurrency"},
                Links: []components.ProjectLink{
                    {Label: "Code", Href: "#", Kind: "code"},
                    {Label: "Code", Href: "#", Kind: "code"},
                    {Label: "Blog Post", Href: "#", Kind: "blog"},
                },
                Thumbnail: fmt.Sprintf("%s/assets/compressed/convert_cbz.webp", types.EVAssetCDNHostname.Get().Value),
            },
        },
        ExperienceSectionData: []components.ExperienceEntry{
            {
                Name:        "Birth",
                Description: "ST. Augustine Hospital • New born",
                DateRange:   "November 2007",
                IsActive:    false,
            },
            {
                Name:        "Student",
                Description: "Holy Child • Primary School",
                DateRange:   "Jan 2009 – Dec 2013",
                IsActive:    false,
            },
            {
                Name:        "Student",
                Description: "ST. John's • High School",
                DateRange:   "Jan 2014 – April 2025",
                IsActive:    false,
            },
            {
                Name:        "Undergraduate",
                Description: "Assam Don Bosco • University",
                DateRange:   "April 2025 – Present",
                IsActive:    true,
            },
        },
        ContactSectionData: components.ContactData{
            Email: "contact@jelius.dev",
            QuickLinks: []components.QuickContactLinks{
                {
                    Platform: "X",
                    URL:      "#",
                    Icon:     "X",
                },
                {
                    Platform: "LinkedIn",
                    URL:      "#",
                    Icon:     "in",
                },
                {
                    Platform: "GitHub",
                    URL:      "#",
                    Icon:     "gh",
                },
                {
                    Platform: "Linktree",
                    URL:      "/links",
                    Icon:     "lt",
                },
            },
        },
    }))
}

