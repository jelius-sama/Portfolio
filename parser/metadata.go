package parser

import (
	"KazuFolio/types"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

func ParseMetadata(path string) (string, error) {
	data, err := os.ReadFile(filepath.Join(os.Getenv("ROOT_PATH"), "config", "static.route.json"))
	if err != nil {
		return "", fmt.Errorf("failed to read metadata file: %w", err)
	}

	var routes []types.RouteMetadata
	if err := json.Unmarshal(data, &routes); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	var defaultMeta, currentMeta, notFoundMeta *types.RouteMetadata

	// Index the metadata by path
	for i := range routes {
		switch routes[i].Path {
		case "*":
			defaultMeta = &routes[i]
		case path:
			currentMeta = &routes[i]
		case "#not_found":
			notFoundMeta = &routes[i]
		}
	}

	// Use #not_found if no match
	if currentMeta == nil {
		currentMeta = notFoundMeta
	}

	if currentMeta == nil {
		// Still no metadata? Nothing to render
		return "", nil
	}

	// Merge defaultMeta + currentMeta
	var mergedMeta = types.RouteMetadata{
		Path:  currentMeta.Path,
		Title: currentMeta.Title,
		Meta:  append(defaultMeta.Meta, currentMeta.Meta...),
		Link:  append(defaultMeta.Link, currentMeta.Link...),
	}

	var builder strings.Builder

	if mergedMeta.Title != "" {
		builder.WriteString(fmt.Sprintf(`<title id="__SERVER_PROPS__">%s</title>`+"\n", html.EscapeString(mergedMeta.Title)))
	}

	for _, meta := range mergedMeta.Meta {
		switch {
		case meta.Charset != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" charset="%s">`+"\n", html.EscapeString(meta.Charset)))
		case meta.Name != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" name="%s" content="%s">`+"\n", html.EscapeString(meta.Name), html.EscapeString(meta.Content)))
		case meta.Property != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" property="%s" content="%s">`+"\n", html.EscapeString(meta.Property), html.EscapeString(meta.Content)))
		case meta.HTTPEquiv != "":
			builder.WriteString(fmt.Sprintf(`<meta id="__SERVER_PROPS__" http-equiv="%s" content="%s">`+"\n", html.EscapeString(meta.HTTPEquiv), html.EscapeString(meta.Content)))
		}
	}

	for _, link := range mergedMeta.Link {
		builder.WriteString(`<link id="__SERVER_PROPS__"`)
		builder.WriteString(fmt.Sprintf(` rel="%s"`, html.EscapeString(link.Rel)))
		builder.WriteString(fmt.Sprintf(` href="%s"`, html.EscapeString(link.Href)))

		if link.Type != "" {
			builder.WriteString(fmt.Sprintf(` type="%s"`, html.EscapeString(link.Type)))
		}
		if link.Crossorigin != "" {
			builder.WriteString(fmt.Sprintf(` crossorigin="%s"`, html.EscapeString(link.Crossorigin)))
		}
		if link.Media != "" {
			builder.WriteString(fmt.Sprintf(` media="%s"`, html.EscapeString(link.Media)))
		}
		if link.Sizes != "" {
			builder.WriteString(fmt.Sprintf(` sizes="%s"`, html.EscapeString(link.Sizes)))
		}
		if link.As != "" {
			builder.WriteString(fmt.Sprintf(` as="%s"`, html.EscapeString(link.As)))
		}
		if link.Referrer != "" {
			builder.WriteString(fmt.Sprintf(` referrerpolicy="%s"`, html.EscapeString(link.Referrer)))
		}
		if link.Title != "" {
			builder.WriteString(fmt.Sprintf(` title="%s"`, html.EscapeString(link.Title)))
		}

		builder.WriteString(`>` + "\n")
	}

	return builder.String(), nil
}
