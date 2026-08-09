// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package renderer

import (
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/template"
    "git.jelius.dev/jelius-sama/Portfolio/template/components"
    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/a-h/templ"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

type ViewManager struct{}

func New() *ViewManager {
    return &ViewManager{}
}

func Renderer(c fiber.Ctx, metadata *types.Metadata, bodyContent templ.Component) error {
    var page types.PageContext = types.PageContext{}
    page.DetermineRenderMode(c)
    c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

    if (page.IsPartial && page.TargetPart != types.TPBody) || !page.IsPartial {
        page.Metadata = metadata
    }

    var buf strings.Builder

    if page.IsPartial {
        switch page.TargetPart {
        case types.TPHead:
            if err := template.Metadata(page.Metadata).Render(c.RequestCtx(), &buf); err != nil {
                logger.Error(err)
                return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
            }

            c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
            return c.SendString(buf.String())
        case types.TPBody:
            if err := bodyContent.Render(c.RequestCtx(), &buf); err != nil {
                logger.Error(err)
                return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
            }
        case types.TPBoth:
            if err := bodyContent.Render(c.RequestCtx(), &buf); err != nil {
                return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
            }
            template.Metadata(page.Metadata).Render(c.RequestCtx(), &buf)

            c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
            return c.SendString(buf.String())
        default:
            return fiber.NewError(fiber.StatusBadRequest, "Invalid SPA Target")
        }
    }

    var fullPage templ.Component = template.Base(c, page.Metadata, bodyContent, components.Footer())

    if err := fullPage.Render(c.Req().RequestCtx(), &buf); err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
    return c.SendString(buf.String())
}

