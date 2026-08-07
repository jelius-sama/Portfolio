package renderer

import (
    "database/sql"
    "fmt"
    "net/url"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/gofiber/fiber/v3"
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

// TODO: Handle special case paths such as `#not_found`
func GetMetadata(c fiber.Ctx) (*types.Metadata, error) {
    var metadata types.Metadata
    var path = c.Path()
    var ctx = c.RequestCtx()

    metadata.Path = path

    if err := db.DB.QueryRowContext(ctx,
        `SELECT title, description FROM metadata WHERE path = ?`,
        path,
    ).Scan(&metadata.Title, &metadata.Description); err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("metadata not found for path %q", path)
        }
        return nil, err
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

