package renderer

import (
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/gofiber/fiber/v3"
)

func (v *ViewManager) RenderBlogs(c fiber.Ctx) error {
    var metadata types.Metadata = types.Metadata{
        Title:       "Blogs",
        Description: "My blog posts",
    }
    return Renderer(c, metadata, pages.Blogs())
}

