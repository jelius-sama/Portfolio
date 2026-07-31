package middleware

import (
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/renderer"
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func ErrHandler(c fiber.Ctx, err error) error {
    var code int
    if err == nil {
        logger.Panic("Expected error got nil, panicked to prevent nil pointer deference, this should be caught by panic handler middleware.")
    }

    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
    }

    if path := c.Path(); strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/assets/") {
        return c.Status(code).JSON(fiber.Map{
            "status":  code,
            "message": err.Error(),
        })
    }

    c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
    c.Status(fiber.StatusNotFound)

    return renderer.Renderer(c, types.Metadata{}, pages.Error(code, err))
}

