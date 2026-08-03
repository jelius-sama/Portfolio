package analytics

import (
    "strconv"

    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func GetAvgVisitsPerHour(c fiber.Ctx) error {
    // Query parameters
    var hoursStr = c.Query("hours", "1")
    var pageStr = c.Query("page", "0")
    var limitStr = c.Query("limit", "10")
    var sortOrder = c.Query("sort", "0")

    var page, limit, sort, parseErr = parsePaginationParam(pageStr, limitStr, sortOrder)
    if parseErr != nil {
        return c.Status(fiber.StatusBadRequest).JSON(*parseErr)
    }

    var hours, hrsParseErr = strconv.Atoi(hoursStr)
    if hrsParseErr != nil || hours < 1 {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "hours must be a positive integer",
        })
    }

    // Get total count
    var countQuery = `
        SELECT COUNT(DISTINCT DATE_FORMAT(timestamp, '%Y-%m-%d %H:00:00'), country_code, page_path)
        FROM analytics_events
        WHERE timestamp >= DATE_SUB(NOW(), INTERVAL ? HOUR)
    `

    var totalRows int
    if err := db.DB.QueryRow(countQuery, hours).Scan(&totalRows); err != nil {
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
            DATE_FORMAT(timestamp, '%Y-%m-%d %H:00:00') as time_window,
            country_code,
            page_path,
            COUNT(*) as visit_count,
            DATE(timestamp) as date
        FROM analytics_events
        WHERE timestamp >= DATE_SUB(NOW(), INTERVAL ? HOUR)
        GROUP BY time_window, country_code, page_path
        ORDER BY time_window ` + sortDir + `, country_code
        LIMIT ? OFFSET ?
    `

    var rows, queryErr = db.DB.Query(query, hours, limit, offset)
    if queryErr != nil {
        logger.Error(c.Path(), queryErr.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }
    defer rows.Close()

    var data []types.AvgVisitsResponse
    for rows.Next() {
        var resp types.AvgVisitsResponse
        var visitCount int
        if err := rows.Scan(
            &resp.TimeWindow,
            &resp.CountryCode,
            &resp.PagePath,
            &visitCount,
            &resp.Date,
        ); err != nil {
            logger.Error(c.Path(), err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: "Internal Server Error",
            })
        }
        resp.AvgVisits = float64(visitCount) / float64(hours)
        data = append(data, resp)
    }

    var hasMore bool = (offset + limit) < totalRows

    return c.Status(fiber.StatusOK).JSON(types.PaginatedAvgVisitsResponse{
        Data:      data,
        Page:      page,
        Limit:     limit,
        HasMore:   hasMore,
        TotalRows: totalRows,
    })
}

