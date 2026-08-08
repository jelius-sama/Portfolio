package renderer

import (
    "database/sql"
    "fmt"
    "net/url"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func normalizeURL(href string) string {
    // Skip obviously non-URL content
    if href == "" || strings.ContainsAny(href, " \t\n") || !looksLikeURL(href) {
        return href
    }

    if l, err := url.Parse(href); err == nil && !l.IsAbs() {
        return fmt.Sprintf("%s%s", types.EVAssetCDNHostname.Get().Value, href)
    }
    return href
}

func looksLikeURL(href string) bool {
    // Must start with / (relative path), http://, https://, or a common scheme
    return strings.HasPrefix(href, "/") ||
        strings.HasPrefix(href, "http://") ||
        strings.HasPrefix(href, "https://") ||
        strings.HasPrefix(href, "data:") ||
        strings.HasPrefix(href, "blob:") ||
        (strings.Contains(href, "://") && !strings.ContainsAny(href, " \t\n="))
}

func preProcessMetadata(metadata *types.Metadata) *types.Metadata {
    for i := range metadata.Links {
        metadata.Links[i].Href = normalizeURL(metadata.Links[i].Href)
    }

    for i := range metadata.Meta {
        metadata.Meta[i].Content = normalizeURL(metadata.Meta[i].Content)
    }

    return metadata
}

func GetMetadata(c fiber.Ctx) (*types.Metadata, error) {
    var metadata types.Metadata
    var ctx = c.RequestCtx()
    var path string

    if pp, ok := c.Locals("pseudo_path").(string); ok && len(pp) != 0 {
        path = pp
    } else {
        path = c.Path()
    }

    metadata.Path = path

    if path != "*" {
        if err := db.DB.QueryRowContext(ctx,
            `SELECT title, description FROM metadata WHERE path = ?`,
            path,
        ).Scan(&metadata.Title, &metadata.Description); err != nil {
            if err == sql.ErrNoRows {
                return nil, fmt.Errorf("metadata not found for path %q", path)
            }
            return nil, err
        }
    }

    if linkRows, err := db.DB.QueryContext(ctx,
        `SELECT rel, href, media FROM m_links WHERE metadata_path = ? OR metadata_path = '*'`,
        path,
    ); err != nil {
        return nil, err
    } else {
        defer linkRows.Close()

        for linkRows.Next() {
            var link types.MLink
            var media sql.NullString

            if err := linkRows.Scan(&link.Rel, &link.Href, &media); err != nil {
                return nil, err
            }

            if media.Valid {
                link.Media = &media.String
            }

            metadata.Links = append(metadata.Links, link)
        }
        if err = linkRows.Err(); err != nil {
            return nil, err
        }
    }

    if metaRows, err := db.DB.QueryContext(ctx,
        `SELECT name, property, content FROM m_meta WHERE metadata_path = ? OR metadata_path = '*'`,
        path,
    ); err != nil {
        return nil, err
    } else {
        defer metaRows.Close()

        for metaRows.Next() {
            var m types.MMeta
            var name, property sql.NullString

            if err := metaRows.Scan(&name, &property, &m.Content); err != nil {
                return nil, err
            }

            if name.Valid {
                m.Name = &name.String
            }
            if property.Valid {
                m.Property = &property.String
            }

            metadata.Meta = append(metadata.Meta, m)
        }
        if err = metaRows.Err(); err != nil {
            return nil, err
        }
    }

    return preProcessMetadata(&metadata), nil
}

func GetDynamicRouteMetadata(c fiber.Ctx, metadata *types.Metadata) {
    var host, err = url.Parse(types.EVHostname.Get().Value)
    if err != nil {
        logger.Error(c.Path(), err)
        // TODO: If for some reason `url.Parse` fails (I doubt it would), just assume a static domain name and manually generate a URL object
        // type URL struct {
        //     Scheme   string
        //     Opaque   string    // encoded opaque data
        //     User     *Userinfo // username and password information
        //     Host     string    // "host" or "host:port" (see Hostname and Port methods)
        //     Path     string    // path (relative paths may omit leading slash)
        //     Fragment string    // fragment for references (without '#')
        //
        //     // RawQuery contains the encoded query values, without the initial '?'.
        //     // Use URL.Query to decode the query.
        //     RawQuery string
        //
        //     // RawPath is an optional field containing an encoded path hint.
        //     // See the EscapedPath method for more details.
        //     //
        //     // In general, code should call EscapedPath instead of reading RawPath.
        //     RawPath string
        //
        //     // RawFragment is an optional field containing an encoded fragment hint.
        //     // See the EscapedFragment method for more details.
        //     //
        //     // In general, code should call EscapedFragment instead of reading RawFragment.
        //     RawFragment string
        //
        //     // ForceQuery indicates whether the original URL contained a query ('?') character.
        //     // When set, the String method will include a trailing '?', even when RawQuery is empty.
        //     ForceQuery bool
        //
        //     // OmitHost indicates the URL has an empty host (authority).
        //     // When set, the String method will not include the host when it is empty.
        //     OmitHost bool
        // }
    }

    metadata.Title = c.Locals("title").(string)
    metadata.Description = c.Locals("description").(string)
    metadata.Meta = append(metadata.Meta,
        types.MMeta{Name: new("description"), Content: c.Locals("description").(string)},
        types.MMeta{Name: new("application-name"), Content: "Jelius Basumatary"},
        types.MMeta{Name: new("robots"), Content: "index, follow"},
        types.MMeta{Name: new("format-detection"), Content: "telephone=no"},
        types.MMeta{Name: new("apple-mobile-web-app-capable"), Content: "yes"},
        types.MMeta{Name: new("apple-mobile-web-app-title"), Content: "Jelius Basumatary"},
        types.MMeta{Name: new("theme-color"), Content: "#1e1e2e"},
        types.MMeta{Name: new("apple-mobile-web-app-status-bar-style"), Content: "default"},

        types.MMeta{Property: new("og:title"), Content: c.Locals("title").(string)},
        types.MMeta{Property: new("og:description"), Content: c.Locals("description").(string)},
        types.MMeta{Property: new("og:url"), Content: host.JoinPath(c.Path()).String()},
        types.MMeta{Property: new("og:site_name"), Content: "Jelius Basumatary"},
        types.MMeta{Property: new("og:image"), Content: "/assets/compressed/jelius.webp"},
        types.MMeta{Property: new("og:type"), Content: "article"},

        types.MMeta{Name: new("twitter:card"), Content: "summary"},
        types.MMeta{Name: new("twitter:site"), Content: "@jelius_sama"},
        types.MMeta{Name: new("twitter:creator"), Content: "@jelius_sama"},
        types.MMeta{Name: new("twitter:title"), Content: c.Locals("title").(string)},
        types.MMeta{Name: new("twitter:description"), Content: c.Locals("description").(string)},
        types.MMeta{Name: new("twitter:image"), Content: "/assets/compressed/jelius.webp"},
    )

    metadata.Links = append(metadata.Links,
        types.MLink{Rel: "canonical", Href: host.JoinPath(c.Path()).String()},
    )

    preProcessMetadata(metadata)
}

