package renderer

import (
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/gofiber/fiber/v3"
)

func (v *ViewManager) RenderHome(c fiber.Ctx) error {
    var metadata types.Metadata = types.Metadata{
        Title:       "Home",
        Description: "Portfolio",
    }
    return Renderer(c, metadata, pages.Home())
}

