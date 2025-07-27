package handler

import (
	"KazuFolio/types"
	"fmt"
	"net/http"
	"strings"
)

var config types.RobotsConfig = types.RobotsConfig{
	Rules: []types.RobotsRule{
		{
			UserAgent: "*",
			Allow:     []string{"/"},
			Disallow:  []string{"/analytics"},
		},
	},
	Host: "https://jelius.dev",
	Sitemaps: []string{
		"https://jelius.dev/sitemap.xml",
	},
}

// INFO: GenerateRobots generates and writes robots.txt content to the response
func GenerateRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	var builder strings.Builder

	// INFO: Write user-agent groups
	for _, rule := range config.Rules {
		builder.WriteString(fmt.Sprintf("User-agent: %s\n", rule.UserAgent))
		for _, allow := range rule.Allow {
			builder.WriteString(fmt.Sprintf("Allow: %s\n", allow))
		}
		for _, disallow := range rule.Disallow {
			builder.WriteString(fmt.Sprintf("Disallow: %s\n", disallow))
		}
		if rule.CrawlDelay != nil {
			builder.WriteString(fmt.Sprintf("Crawl-delay: %d\n", *rule.CrawlDelay))
		}
		builder.WriteString("\n")
	}

	if config.Host != "" {
		builder.WriteString(fmt.Sprintf("Host: %s\n\n", config.Host))
	}

	// Sitemap directives
	for _, sitemap := range config.Sitemaps {
		builder.WriteString(fmt.Sprintf("Sitemap: %s\n", sitemap))
	}

	_, _ = w.Write([]byte(builder.String()))
}
