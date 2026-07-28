package renderer

import (
    "github.com/gofiber/fiber/v3"
)

func (v *ViewManager) RenderHome(c fiber.Ctx) error {
    var page PageContext = PageContext{Title: "Home"}
    page.DetermineRenderMode(c)

    var data fiber.Map = fiber.Map{
        "Title": page.Title,
    }

    if page.IsPartial {
        switch page.TargetPart {
        case TPHead:
            return c.Render("template/head", data)
        case TPBody:
            return c.Render("template/home", data)
        default:
            return fiber.NewError(fiber.StatusBadRequest, "Invalid SPA Target")
        }
    }

    return c.Render("template/home", data, "template/base")
}

