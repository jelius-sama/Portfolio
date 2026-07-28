package main

import (
    "git.jelius.dev/jelius-sama/Portfolio/api"
    "git.jelius.dev/jelius-sama/Portfolio/middleware"
    "github.com/gofiber/fiber/v3"
)

func Router(app *fiber.App) {
    app.Use(middleware.RequestLogger())
    app.Get("/api/healthz", api.Healthz)
}

