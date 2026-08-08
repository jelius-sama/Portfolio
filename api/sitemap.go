package api

import (
    "net/url"
    "time"

    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
)

func GenerateSitemap(c fiber.Ctx) error {
    var host, err = url.Parse(types.EVHostname.Get().Value)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }

    var now = time.Now().Format("2006-01-02")

    // Static routes
    var urls = []types.SiteMapURLEntry{
        {Loc: host.JoinPath("/").String(), LastMod: now, ChangeFreq: "daily", Priority: "1.0"},
        {Loc: host.JoinPath("/achievements").String(), LastMod: now, ChangeFreq: "daily", Priority: "1"},
        {Loc: host.JoinPath("/blogs").String(), LastMod: now, ChangeFreq: "daily", Priority: "0.8"},
        {Loc: host.JoinPath("/links").String(), LastMod: now, ChangeFreq: "monthly", Priority: "0.6"},
    }

    // Fetch all blogs directly
    if rows, err := db.DB.Query(`
        SELECT id, updated_at
        FROM blogs
        ORDER BY updated_at DESC
    `); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    } else {
        defer rows.Close()

        for rows.Next() {
            var id, updatedAt string
            if err := rows.Scan(&id, &updatedAt); err != nil {
                continue // skip malformed rows
            }

            if t, err := time.Parse(time.RFC3339, updatedAt); err != nil {
                continue
            } else {
                updatedAt = t.Format("2006-01-02")
            }

            urls = append(urls, types.SiteMapURLEntry{
                Loc:        host.JoinPath("blog", id).String(),
                LastMod:    updatedAt,
                ChangeFreq: "monthly",
                Priority:   "0.7",
            })
        }
    }

    return c.XML(types.SiteMapURLSet{
        Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
        URLs:  urls,
    })
}

