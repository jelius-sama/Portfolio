package middleware

import (
    "strings"

    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func RecoveryMiddleware() fiber.Handler {
    return func(c fiber.Ctx) error {
        defer func() {
            if err := recover(); err != nil {
                logger.Error(err)
                logger.Error("Encountered a panic, returning 500 to client and recovering the server ASAP!")
                if strings.HasPrefix(c.Path(), "/api/") || strings.HasPrefix(c.Path(), "/assets/") {
                    c.Set("Content-Type", "application/json")
                    c.Status(500)
                    c.Response().SetBody([]byte(`{"status":500,"message":"Internal Server Error"}`))
                    return
                }

                c.Set("Content-Type", "text/html; charset=utf-8")
                c.Status(500)
                c.Response().AppendBodyString(`
                    <main id="internal-server-error">
                        <h2>500</h2>
                        <p>Internal Server Error</p>
                    </main>
                `)
                return
            }
        }()

        return c.Next()
    }
}

