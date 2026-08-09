// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package analytics

import (
    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func GetTopPages(c fiber.Ctx) error {
    // Query parameters
    var pageStr = c.Query("page", "0")
    var limitStr = c.Query("limit", "10")
    var sortOrder = c.Query("sort", "1") // asc or desc (default desc for top pages)

    var page, limit, sort, parseErr = parsePaginationParam(pageStr, limitStr, sortOrder)
    if parseErr != nil {
        return c.Status(fiber.StatusBadRequest).JSON(*parseErr)
    }

    // Get total count of distinct pages
    var countQuery = `SELECT COUNT(DISTINCT page_path) FROM analytics_events`

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

    // Get total visits for percentage calculation
    var totalVisits int
    var totalVisitsQuery = `SELECT COUNT(*) FROM analytics_events`
    if err := db.DB.QueryRow(totalVisitsQuery).Scan(&totalVisits); err != nil {
        logger.Error(c.Path(), err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }

    // Query with pagination
    var sortDir = "DESC"
    if sort == types.SOAsc {
        sortDir = "ASC"
    }

    var query = `
        SELECT 
            page_path,
            COUNT(*) as visit_count
        FROM analytics_events
        GROUP BY page_path
        ORDER BY visit_count ` + sortDir + `
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

    var data []types.TopPageResponse
    for rows.Next() {
        var resp types.TopPageResponse
        if err := rows.Scan(
            &resp.PagePath,
            &resp.VisitCount,
        ); err != nil {
            logger.Error(c.Path(), err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: "Internal Server Error",
            })
        }
        if totalVisits > 0 {
            resp.Percentage = (float64(resp.VisitCount) / float64(totalVisits)) * 100
        }
        data = append(data, resp)
    }

    var hasMore bool = (offset + limit) < totalRows

    return c.Status(fiber.StatusOK).JSON(types.PaginatedTopPagesResponse{
        Data:      data,
        Page:      page,
        Limit:     limit,
        HasMore:   hasMore,
        TotalRows: totalRows,
    })
}

