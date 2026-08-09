// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package api

import (
    "git.jelius.dev/jelius-sama/Portfolio/db"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func Healthz(c fiber.Ctx) error {
    if err := db.DB.Ping(); err != nil {
        logger.Fatal("Failed to ping database:", err.Error())
        return c.Status(fiber.StatusServiceUnavailable).SendString("Database down or degraded!")
    }

    return c.SendString("Healthy")
}

