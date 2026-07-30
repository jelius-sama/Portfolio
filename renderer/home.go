package renderer

import (
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

