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

    app.Get("/api/healthz", api.Healthz)

    if types.EVEnv.Get().Value == types.EMDev.String() {
        sub, _ := fs.Sub(embed.AssetFS, "assets")
        app.Get("/assets/*", static.New("", static.Config{
            FS:            sub,
            CacheDuration: 0,
            Compress:      true,
            NotFoundHandler: func(c fiber.Ctx) error {
                return middleware.ErrHandler(c, fiber.ErrNotFound)
            },
        }))
    }

    ui := renderer.New()
    app.Get("/", ui.RenderHome)
    app.Get("/links", ui.RenderLinks)
}

