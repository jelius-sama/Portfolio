// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package analytics

import (
    "database/sql"
    "errors"
    "strconv"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

// GetPageVisitCount returns the total number of visits for a specific page
func GetPageVisitCount(c fiber.Ctx, buf ...*strings.Builder) error {
    var pagePath string

    if len(buf) != 0 {
        pagePath = buf[0].String()
    } else {
        pagePath = c.Query("page")
        if len(pagePath) == 0 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "page query parameter is required",
            })
        }
    }

    if _, exists := types.Pages[pagePath]; !exists {
        // coule be a dynamic route if not a direct match
        for templatePattern, _ := range types.Pages {
            // Evaluate the raw URL against the map key pattern
            if fiber.RoutePatternMatch(pagePath, templatePattern) {
                exists = true
                break
            }
        }

        if !exists && len(buf) == 0 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "invalid page parameter provided",
            })
        }

        if !exists && len(buf) != 0 {
            return errors.New("invalid page path")
        }
    }

    // Query to count visits for the specific page
    var query = `
    SELECT COUNT(*) as visit_count
    FROM analytics_events
    WHERE page_path = ?
    `

    var visitCount int
    if err := db.DB.QueryRow(query, pagePath).Scan(&visitCount); err != nil && err != sql.ErrNoRows {
        if len(buf) == 0 {
            logger.Error(c.Path(), err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: "Internal Server Error",
            })
        } else {
            return err
        }
    }

    if len(buf) == 0 {
        return c.SendString(strconv.Itoa(visitCount))
    } else {
        buf[0].Reset()
        buf[0].WriteString(strconv.Itoa(visitCount))
    }

    return nil
}

