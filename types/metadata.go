// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package types

import "encoding/xml"

type MLink struct {
    Rel   string  `json:"rel"`
    Href  string  `json:"href"`
    Media *string `json:"media"`
}

type MMeta struct {
    Name     *string `json:"name"`
    Property *string `json:"property"`
    Content  string  `json:"content"`
}

type Metadata struct {
    Path        string  `json:"path"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Links       []MLink `json:"links"`
    Meta        []MMeta `json:"meta"`
}

type SiteMapURLEntry struct {
    Loc        string `xml:"loc"`
    LastMod    string `xml:"lastmod,omitempty"`
    ChangeFreq string `xml:"changefreq,omitempty"`
    Priority   string `xml:"priority,omitempty"`
}

type SiteMapURLSet struct {
    XMLName xml.Name          `xml:"urlset"`
    Xmlns   string            `xml:"xmlns,attr"`
    URLs    []SiteMapURLEntry `xml:"url"`
}

type RobotsRule struct {
    UserAgent  string
    Disallow   []string
    Allow      []string
    CrawlDelay *int
}

type RobotsConfig struct {
    Rules    []RobotsRule
    Host     string
    Sitemaps []string
}

