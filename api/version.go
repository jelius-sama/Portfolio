// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package api

import (
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
)

func Version(c fiber.Ctx) error {
    return c.SendString(types.EVVersion.Get().Value)
}

