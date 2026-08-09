// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package api

import (
    "fmt"
    "net/url"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
)

func GenerateRobots(c fiber.Ctx) error {
    var host, err = url.Parse(types.EVHostname.Get().Value)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }

    var config types.RobotsConfig = types.RobotsConfig{
        Rules: []types.RobotsRule{
            {
                UserAgent: "*",
                Allow:     []string{"/"},
                Disallow:  []string{"/analytics", "/music"},
            },
        },
        Host: types.EVHostname.Get().Value,
        Sitemaps: []string{
            host.JoinPath("sitemap.xml").String(),
        },
    }

    var builder strings.Builder

    for _, rule := range config.Rules {
        fmt.Fprintf(&builder, "User-agent: %s\n", rule.UserAgent)
        for _, allow := range rule.Allow {
            fmt.Fprintf(&builder, "Allow: %s\n", allow)
        }
        for _, disallow := range rule.Disallow {
            fmt.Fprintf(&builder, "Disallow: %s\n", disallow)
        }
        if rule.CrawlDelay != nil {
            fmt.Fprintf(&builder, "Crawl-delay: %d\n", *rule.CrawlDelay)
        }
        builder.WriteString("\n")
    }

    if len(config.Host) != 0 {
        fmt.Fprintf(&builder, "Host: %s\n\n", config.Host)
    }

    for _, sitemap := range config.Sitemaps {
        fmt.Fprintf(&builder, "Sitemap: %s\n", sitemap)
    }

    return c.SendString(builder.String())
}

