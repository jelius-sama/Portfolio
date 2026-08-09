// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package analytics

import (
    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func GetAllAnalyticsEvents(c fiber.Ctx) error {
    // Query parameters
    var pageStr = c.Query("page", "0")
    var limitStr = c.Query("limit", "10")
    var sortOrder = c.Query("sort", "1") // asc or desc

    var page, limit, sort, parseErr = parsePaginationParam(pageStr, limitStr, sortOrder)
    if parseErr != nil {
        return c.Status(fiber.StatusBadRequest).JSON(*parseErr)
    }

    // Get total count
    var countQuery = `SELECT COUNT(*) FROM analytics_events`

    var totalRows int
    if err := db.DB.QueryRow(countQuery).Scan(&totalRows); err != nil {
        logger.Error(c.Path(), err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }

    // Calculate offset
    var offset = page * limit

    // Query with pagination
    var sortDir = "ASC"
    if sort == types.SODesc {
        sortDir = "DESC"
    }

    var query = `
        SELECT 
            event_id,
            country_code,
            page_path,
            strftime('%Y-%m-%d %H:%M:%S', timestamp) as timestamp
        FROM analytics_events
        ORDER BY timestamp ` + sortDir + `
        LIMIT ? OFFSET ?
    `

    var rows, queryErr = db.DB.Query(query, limit, offset)
    if queryErr != nil {
        logger.Error(c.Path(), queryErr.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }
    defer rows.Close()

    var data []types.AnalyticsEventResponse
    for rows.Next() {
        var resp types.AnalyticsEventResponse
        if err := rows.Scan(
            &resp.EventID,
            &resp.CountryCode,
            &resp.PagePath,
            &resp.TimestampUTC,
        ); err != nil {
            logger.Error(c.Path(), err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusBadRequest,
                Message: "Internal Server Error",
            })
        }
        data = append(data, resp)
    }

    var hasMore bool = (offset + limit) < totalRows

    return c.Status(fiber.StatusOK).JSON(types.PaginatedEventsResponse{
        Data:      data,
        Page:      page,
        Limit:     limit,
        HasMore:   hasMore,
        TotalRows: totalRows,
    })
}

