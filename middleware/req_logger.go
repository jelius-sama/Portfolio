// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package middleware

import (
    "strings"

    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func RequestLogger() fiber.Handler {
    return func(c fiber.Ctx) error {
        if path := c.Path(); strings.HasPrefix(path, "/api/") {
            logger.Info(c.Method(), path)
        }

        return c.Next()
    }
}

