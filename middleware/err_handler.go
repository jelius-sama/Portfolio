package middleware

import (
    "fmt"
    "strconv"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/renderer"
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func ErrHandler(c fiber.Ctx, err error) error {
    c.Response().Header.Del("Cache-Control")
    var code int
    if err == nil {
        logger.Panic("Expected error got nil, panicked to prevent nil pointer deference, this should be caught by panic handler middleware.")
    }

    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
    }

    if path := c.Path(); strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/assets/") {
        if err := c.Status(code).JSON(fiber.Map{
            "status":  code,
            "message": err.Error(),
        }); err != nil {
            logger.Panic(err)
        }
        return nil
    }

    c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
    c.Status(code)

    var metadata types.Metadata = types.Metadata{
        Title:       strconv.Itoa(code),
        Description: err.Error(),
    }
    if renderErr := renderer.Renderer(c, metadata, pages.Error(code, err)); renderErr != nil {
        logger.Panic(fmt.Sprintf(
            "double fault in error middleware:\n"+
                " -> Original Error: %v\n"+
                " -> Render Error  : %v",
            err, renderErr,
        ))
    }

    return nil
}

