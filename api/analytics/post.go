package analytics

import (
    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func TrackAnalytics(c fiber.Ctx) error {
    var req types.TrackAnalyticsRequest

    // Parse request body
    if err := c.Bind().JSON(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "Invalid request body",
        })
    }

    // Validate required fields
    if len(req.CountryCode) == 0 || len(req.PagePath) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "country_code and page_path are required",
        })
    }

    // Insert into database
    var query = `
        INSERT INTO analytics_events (country_code, page_path, timestamp)
        VALUES (?, ?, datetime('now', 'utc'))
    `

    if _, err := db.DB.Exec(query, req.CountryCode, req.PagePath); err != nil {
        logger.Error(c.Path(), err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }

    return c.SendStatus(fiber.StatusCreated)
}

