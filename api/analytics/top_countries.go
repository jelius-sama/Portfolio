package analytics

import (
    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func GetTopCountries(c fiber.Ctx) error {
    // Query parameters
    var pageStr = c.Query("page", "0")
    var limitStr = c.Query("limit", "10")
    var sortOrder = c.Query("sort", "1") // asc or desc (default desc for top countries)

    var page, limit, sort, parseErr = parsePaginationParam(pageStr, limitStr, sortOrder)
    if parseErr != nil {
        return c.Status(fiber.StatusBadRequest).JSON(*parseErr)
    }

    // Get total count of distinct countries
    var countQuery = `SELECT COUNT(DISTINCT country_code) FROM analytics_events`

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

    query := `
        SELECT 
            country_code,
            COUNT(*) as visit_count
        FROM analytics_events
        GROUP BY country_code
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

    var data []types.TopCountryResponse
    for rows.Next() {
        var resp types.TopCountryResponse
        if err := rows.Scan(
            &resp.CountryCode,
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

    return c.Status(fiber.StatusOK).JSON(types.PaginatedTopCountriesResponse{
        Data:      data,
        Page:      page,
        Limit:     limit,
        HasMore:   hasMore,
        TotalRows: totalRows,
    })
}

