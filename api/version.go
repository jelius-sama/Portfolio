package api

import (
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
)

func Version(c fiber.Ctx) error {
    return c.SendString(types.EVVersion.Get().Value)
}

