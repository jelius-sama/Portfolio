package renderer

import (
    "context"
    "database/sql"
    "fmt"

    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/template/components"
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

// GetLinksPage retrieves a links page by handle and returns pages.LinksArg
func getLinksPageData(ctx context.Context) (pages.LinksArg, error) {
    var linksArg pages.LinksArg
    var footerWhoAmI, qrImagePath sql.NullString

    // Fetch main links page data
    var pageQuery = `
        SELECT 
            handle, tag_line, image_path, image_alt_text, who_am_i, qr_image_path
        FROM links_page
        ORDER BY updated_at DESC
        LIMIT 1
    `

    if err := db.DB.QueryRowContext(ctx, pageQuery).Scan(
        &linksArg.Handle,
        &linksArg.TagLine,
        &linksArg.ImagePath,
        &linksArg.ImageAltText,
        &footerWhoAmI,
        &qrImagePath,
    ); err != nil {
        if err == sql.ErrNoRows {
            return pages.LinksArg{}, fmt.Errorf("links page not found")
        }
        return pages.LinksArg{}, err
    }

    // Initialize footer
    linksArg.Footer = components.LinksFooter{
        WhoAmI:      footerWhoAmI.String,
        QRImagePath: qrImagePath.String,
    }

    // Fetch all link entries ordered by position
    var entriesQuery = `
        SELECT 
            icon, title, subtitle, href
        FROM link_entries
        WHERE links_page_id = (SELECT id FROM links_page ORDER BY updated_at DESC LIMIT 1)
        ORDER BY position ASC
    `

    var rows, err = db.DB.QueryContext(ctx, entriesQuery)
    if err != nil {
        return pages.LinksArg{}, err
    }
    defer rows.Close()

    var links []components.LinkEntry
    for rows.Next() {
        var entry components.LinkEntry
        var subtitle sql.NullString

        if err := rows.Scan(
            &entry.Icon,
            &entry.Title,
            &subtitle,
            &entry.Href,
        ); err != nil {
            return pages.LinksArg{}, err
        }

        entry.Subtitle = subtitle.String
        links = append(links, entry)
    }

    if err = rows.Err(); err != nil {
        return pages.LinksArg{}, err
    }

    linksArg.Links = links
    return linksArg, nil
}

func (v *ViewManager) RenderLinks(c fiber.Ctx) error {
    var metadata types.Metadata = types.Metadata{
        Title:       "Links",
        Description: "My Social Links",
    }

    if data, err := getLinksPageData(c.RequestCtx()); err == nil {
        return Renderer(c, metadata, pages.Links(data))
    } else {
        logger.Error(c.Path(), err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }
}

