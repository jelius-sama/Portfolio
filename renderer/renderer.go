package renderer

import (
    "git.jelius.dev/jelius-sama/Portfolio/template"
    "git.jelius.dev/jelius-sama/Portfolio/template/components"
    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/a-h/templ"
    "github.com/gofiber/fiber/v3"
)

func Renderer(c fiber.Ctx, metadata types.Metadata, bodyContent templ.Component) error {
    var page PageContext = PageContext{}
    page.DetermineRenderMode(c)
    c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

    if (page.IsPartial && page.TargetPart != TPBody) || !page.IsPartial {
        page.Metadata = metadata
    }

    if page.IsPartial {
        switch page.TargetPart {
        case TPHead:
            return template.Metadata(page.Metadata).Render(c.RequestCtx(), c.Response().BodyWriter())
        case TPBody:
            return bodyContent.Render(c.RequestCtx(), c.Response().BodyWriter())
        case TPBoth:
            if err := bodyContent.Render(c.RequestCtx(), c.Response().BodyWriter()); err != nil {
                return err
            }

            return template.Metadata(page.Metadata).Render(c.RequestCtx(), c.Response().BodyWriter())
        default:
            return fiber.NewError(fiber.StatusBadRequest, "Invalid SPA Target")
        }
    }

    var fullPage templ.Component = template.Base(template.Metadata(page.Metadata), bodyContent, components.Footer())

    return fullPage.Render(c.Req().RequestCtx(), c.Response().BodyWriter())
}

