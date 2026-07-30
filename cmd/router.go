package main

import (
    "io/fs"

    embed "git.jelius.dev/jelius-sama/Portfolio"
    "git.jelius.dev/jelius-sama/Portfolio/api"
    "git.jelius.dev/jelius-sama/Portfolio/middleware"
    "git.jelius.dev/jelius-sama/Portfolio/renderer"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/static"
)

func Router(app *fiber.App) {
    app.Use(middleware.RecoveryMiddleware())
    app.Use(middleware.RequestLogger())

    var noCache = middleware.NewCacheControl(middleware.CacheConfig{
        CustomHeader: "no-store, no-cache, must-revalidate, proxy-revalidate",
    })
    var apiHandle = app.Group("/api", noCache)

    apiHandle.Get("/healthz", api.Healthz)

    if types.EVEnv.Get().Value == types.EMDev.String() {
        var sub, _ = fs.Sub(embed.AssetFS, "assets")
        app.Get("/assets/*", static.New("", static.Config{
            FS:            sub,
            CacheDuration: 0,
            Compress:      true,
            NotFoundHandler: func(c fiber.Ctx) error {
                return middleware.ErrHandler(c, fiber.ErrNotFound)
            },
        }))
    }

    var ui = renderer.New()
    // INFO: Due to how htmx works we might not be able to cache this very well.
    // I don't know if CDN allows a content to be differentiated with headers alone.
    // Since HTMX natively uses headers to communicate, disregarding that fact will
    // make the custom SPA solution fail, for now just stop caching it.
    app.Get("/", noCache, ui.RenderHome)
    app.Get("/links", noCache, ui.RenderLinks)
}

