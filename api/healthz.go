package api

import "github.com/gofiber/fiber/v3"

func Healthz(c fiber.Ctx) error {
    return c.SendString("Healthy")
}

