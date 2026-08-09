package types

import (
    "github.com/gofiber/fiber/v3"
)

type MiddlewareHandler uint8

const (
    MHNoCache MiddlewareHandler = iota
    MHStaticAsset
    MHHTMXCache
    MHStaticPages
)

type MiddlewareHandlerMap map[MiddlewareHandler]fiber.Handler

type Page struct {
    Handler  any
    Handlers []any
}

// Would contain all the register page routes
// This will be initialized by the `Router()` function in `main` package.
var Pages map[string]Page

