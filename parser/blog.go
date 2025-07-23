package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"KazuFolio/api/handler"
	"KazuFolio/types"
)

// ssrBlogPage fetches the blog data, generates metadata HTML + hydration script.
func SSRBlogPage(fullPath string) (string, error, int) {
	// Step 1: Extract blog ID from the path
	id := strings.TrimPrefix(fullPath, "/blog/")
	if id == "" {
		return "", fmt.Errorf("no blog ID found in path"), http.StatusBadRequest
	}

	// Step 2: Call the internal API to get the blog data
	blog, status, err := handler.GetBlogInternal(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch blog data: %w", err), status
	}

	if status != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", status), status
	}

	// Step 3: Decode the JSON into types.Blog
	resp_body, err := json.Marshal(blog)
	if err != nil {
		return "", fmt.Errorf("failed to read API response: %w", err), http.StatusInternalServerError
	}

	// Step 4: Create metadata JSON (same structure as the static one)
	metadata := types.RouteMetadata{
		Path:  fullPath,
		Title: fmt.Sprintf("%s | Jelius", blog.Title),
		Meta: []types.MetaTag{
			{
				Name:    "description",
				Content: blog.Summary,
			},
			{
				Name:    "application-name",
				Content: "Jelius Basumatary",
			},
			{
				Name:    "robots",
				Content: "index, follow",
			},
			{
				Name:    "format-detection",
				Content: "telephone=no",
			},
			{
				Name:    "apple-mobile-web-app-capable",
				Content: "yes",
			},
			{
				Name:    "apple-mobile-web-app-title",
				Content: "Jelius Basumatary",
			},
			{
				Name:    "theme-color",
				Content: "#000b11",
			},
			{
				Name:    "apple-mobile-web-app-status-bar-style",
				Content: "default",
			},
			{
				Property: "og:title",
				Content:  fmt.Sprintf("%s | Jelius", blog.Title),
			},
			{
				Property: "og:description",
				Content:  blog.Summary,
			},
			{
				Property: "og:url",
				Content:  "https://jelius.dev" + fullPath,
			},
			{
				Property: "og:site_name",
				Content:  "Jelius Basumatary",
			},
			{
				Property: "og:image",
				Content:  "/assets/jelius.jpg",
			},
			{
				Property: "og:type",
				Content:  "article",
			},
			{
				Name:    "twitter:card",
				Content: "summary",
			},
			{
				Name:    "twitter:site",
				Content: "@jelius_sama",
			},
			{
				Name:    "twitter:creator",
				Content: "@jelius_sama",
			},
			{
				Name:    "twitter:title",
				Content: fmt.Sprintf("%s | Jelius", blog.Title),
			},
			{
				Name:    "twitter:description",
				Content: blog.Summary,
			},
			{
				Name:    "twitter:image",
				Content: "/assets/jelius.jpg",
			},
		},
		Link: []types.LinkTag{
			{
				Rel:  "canonical",
				Href: "https://jelius.dev" + fullPath,
			},
		},
	}

	// Step 5: Marshal metadata into JSON for the parser
	jsonData, err := json.Marshal([]types.RouteMetadata{metadata})
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata JSON: %w", err), http.StatusInternalServerError
	}

	// Step 6: Parse metadata and generate HTML
	metaHTML, err := ParseMetadata(fullPath, &jsonData)

	if err != nil {
		return "", fmt.Errorf("failed to parse metadata: %w", err), http.StatusInternalServerError
	}

	// Step 7: Prepare final JSON for hydration script
	scriptPayload := map[string]any{
		"path":     fullPath,
		"metadata": metadata,
		"api_resp": json.RawMessage(resp_body),
	}
	scriptJSON, err := json.Marshal(scriptPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal script payload: %w", err), http.StatusInternalServerError
	}

	scriptHTML := `<script id="__SERVER_DATA__" type="application/json">` + htmlEscape(scriptJSON) + `</script>`

	// Step 8: Return combined HTML
	return scriptHTML + "\n" + metaHTML, nil, status
}

// htmlEscape makes sure the JSON inside the <script> doesn't break HTML parsing
func htmlEscape(data []byte) string {
	s := string(data)
	s = strings.ReplaceAll(s, "<", "\\u003c")
	s = strings.ReplaceAll(s, ">", "\\u003e")
	s = strings.ReplaceAll(s, "&", "\\u0026")
	return s
}
