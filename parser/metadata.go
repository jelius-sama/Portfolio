package parser

import (
	"KazuFolio/types"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ParseMetadata() (string, error) {
	data, err := os.ReadFile(filepath.Join(os.Getenv("ROOT_PATH"), "config", "static.route.json"))
	if err != nil {
		return "", fmt.Errorf("failed to read metadata file: %w", err)
	}

	var routes []types.RouteMetadata
	if err := json.Unmarshal(data, &routes); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	var defaultMeta *types.RouteMetadata
	var builder strings.Builder

	// First, find the default "*" route
	for _, route := range routes {
		if route.Path == "*" {
			defaultMeta = &route
			break
		}
	}

	// Now iterate again, skipping "*", and merge
	for _, route := range routes {
		if route.Path == "*" {
			continue
		}

		// Merge default with specific
		combinedMeta := types.RouteMetadata{
			Path:  route.Path,
			Title: route.Title,
			Meta:  append([]types.MetaTag{}, append(defaultMeta.Meta, route.Meta...)...),
			Link:  append([]types.LinkTag{}, append(defaultMeta.Link, route.Link...)...),
		}

		// Render HTML
		if combinedMeta.Title != "" {
			builder.WriteString(fmt.Sprintf(`<title id="__SERVER_PROPS__">%s</title>`+"\n", htmlEscape(combinedMeta.Title)))
		}

		for _, meta := range combinedMeta.Meta {
			switch {
			case meta.Charset != "":
				builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" charset="%s">`+"\n", htmlEscape(meta.Charset)))
			case meta.Name != "":
				builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" name="%s" content="%s">`+"\n", htmlEscape(meta.Name), htmlEscape(meta.Content)))
			case meta.Property != "":
				builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" property="%s" content="%s">`+"\n", htmlEscape(meta.Property), htmlEscape(meta.Content)))
			case meta.HTTPEquiv != "":
				builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" http-equiv="%s" content="%s">`+"\n", htmlEscape(meta.HTTPEquiv), htmlEscape(meta.Content)))
			}
		}

		for _, link := range combinedMeta.Link {
			builder.WriteString(fmt.Sprintf(`<link id="__SERVER_PROPS__" rel="%s" href="%s">`+"\n", htmlEscape(link.Rel), htmlEscape(link.Href)))
		}
	}

	return builder.String(), nil
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		`&`, "&amp;",
		`<`, "&lt;",
		`>`, "&gt;",
		`"`, "&quot;",
	)
	return replacer.Replace(s)
}
