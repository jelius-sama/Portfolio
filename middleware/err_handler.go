package middleware

import (
    "fmt"
    "html"
    "strings"

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

    // TODO: Integrate errors with GoTH stack
    // GoTH == Go + Templ + HTMX
    var htmlResponse = fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Error %d</title>
            <style>
                body { font-family: sans-serif; text-align: center; padding: 50px; background: #f4f6f9; color: #333; }
                h1 { font-size: 50px; color: #e74c3c; margin-bottom: 10px; }
                p { font-size: 18px; color: #666; }
                a { color: #3498db; text-decoration: none; font-weight: bold; }
            </style>
        </head>
        <body>
            <h1>Oops! Error %d</h1>
            <p>%s</p>
            <hr style="max-width: 400px; border: 0; border-top: 1px solid #ccc; margin: 20px auto;">
            <p><a href="/">Go Back Home</a></p>
        </body>
        </html>
    `, code, code, html.EscapeString(err.Error()))

    c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
    return c.Status(code).SendString(htmlResponse)
}

