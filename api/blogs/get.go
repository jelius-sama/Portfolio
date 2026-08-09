// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package blogs

import (
    "bytes"
    "database/sql"
    "encoding/gob"
    "errors"
    "strconv"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/api/analytics"
    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

// NOTE:
// var tempTestDirectCall = func() {
//     buf := bytes.NewBufferString("abc1230") // Blog ID
//     err := blogs.GetAllBlogs(nil, buf)
//     if err == nil {
//         var decodedResponse types.BlogResponse
//         decErr := gob.NewDecoder(buf).Decode(&decodedResponse)
//         if decErr != nil {
//             logger.Error("Failed to decode Gob data:", decErr.Error())
//             return
//         }
//         decodedResponse.Title // access fields
//         // convert into human-readible
//         humanReadable, _ := json.MarshalIndent(decodedResponse, "", "  ")
//         logger.Okay("Buffer output:\n", string(humanReadable))
//     } else {
//         logger.Error("Error output:", err.Error())
//     }
// }

// GetBlog retrieves a blog post by ID with its prequel and sequel chain
func GetBlog(c fiber.Ctx, buf ...*bytes.Buffer) error {
    var id string
    if len(buf) != 0 {
        id = buf[0].String()
    } else {
        id = c.Params("id")
    }

    if len(id) == 0 && len(buf) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "Blog ID is required",
        })
    }

    if len(id) == 0 && len(buf) != 0 {
        return errors.New("Blog ID is not provided")
    }

    var blog, err = getBlogResponse(id)
    if err != nil {
        if len(buf) != 0 {
            return err
        }

        if err == sql.ErrNoRows {
            return c.Status(fiber.StatusNotFound).JSON(types.ErrorResp{
                Code:    fiber.StatusNotFound,
                Message: "Blog not found",
            })
        }
        logger.Error(c.Path(), err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }

    if len(buf) != 0 {
        buf[0].Reset()
        return gob.NewEncoder(buf[0]).Encode(blog)
    }

    return c.Status(fiber.StatusOK).JSON(blog)
}

// Flattened db scanning struct to extract depth metadata
type scannedBlog struct {
    Response  types.BlogResponse
    PrequelID sql.NullString
    SequelID  sql.NullString
    Depth     int
}

func getBlogResponse(targetID string) (*types.BlogResponse, error) {
    var query = `
        WITH RECURSIVE 
        prequels AS (
            SELECT id, title, excerpt, published_at, updated_at, deleted_at, prequel_id, sequel_id, 0 AS depth
            FROM blogs WHERE id = ? AND deleted_at IS NULL
            UNION ALL
            SELECT b.id, b.title, b.excerpt, b.published_at, b.updated_at, b.deleted_at, b.prequel_id, b.sequel_id, p.depth - 1
            FROM blogs b JOIN prequels p ON b.id = p.prequel_id WHERE b.deleted_at IS NULL
        ),
        sequels AS (
            SELECT id, title, excerpt, published_at, updated_at, deleted_at, prequel_id, sequel_id, 0 AS depth
            FROM blogs WHERE id = ? AND deleted_at IS NULL
            UNION ALL
            SELECT b.id, b.title, b.excerpt, b.published_at, b.updated_at, b.deleted_at, b.prequel_id, b.sequel_id, s.depth + 1
            FROM blogs b JOIN sequels s ON b.id = s.sequel_id WHERE b.deleted_at IS NULL
        ),
        combined_chain AS (
            SELECT * FROM prequels UNION SELECT * FROM sequels
        )
        SELECT id, title, excerpt, published_at, updated_at, deleted_at, prequel_id, sequel_id, depth 
        FROM combined_chain ORDER BY depth ASC;
    `

    // Pass targetID twice: once for the prequels CTE, once for the sequels CTE
    var rows, err = db.DB.Query(query, targetID, targetID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    blogMap := make(map[string]*scannedBlog)
    var rowCount int

    for rows.Next() {
        var sb scannedBlog
        var deletedAt sql.NullTime

        err := rows.Scan(
            &sb.Response.ID, &sb.Response.Title, &sb.Response.Excerpt,
            &sb.Response.PublishedAt, &sb.Response.UpdatedAt, &deletedAt,
            &sb.PrequelID, &sb.SequelID, &sb.Depth,
        )
        if err != nil {
            return nil, err
        }
        if deletedAt.Valid {
            sb.Response.DeletedAt = &deletedAt.Time
        }

        blogMap[sb.Response.ID] = &sb
        rowCount++
    }
    rows.Close() // Close early to free the database file lock before fetching analytics

    if rowCount == 0 {
        return nil, sql.ErrNoRows
    }

    // First, populate view counts sequentially
    for _, item := range blogMap {
        var viewCountBuf strings.Builder
        viewCountBuf.WriteString("/blog/")
        viewCountBuf.WriteString(item.Response.ID)
        if err := analytics.GetPageVisitCount(nil, &viewCountBuf); err == nil {
            if views, err := strconv.Atoi(viewCountBuf.String()); err == nil {
                item.Response.Views = uint(views)
            }
        }
    }

    // Stitch pointers based on depth to prevent JSON marshal cycles
    for _, item := range blogMap {
        // Root item (Depth == 0) gets BOTH prequel and sequel tracks populated
        if item.Depth == 0 {
            if item.PrequelID.Valid && item.PrequelID.String != "" {
                if pRow, exists := blogMap[item.PrequelID.String]; exists {
                    item.Response.Prequel = &pRow.Response
                }
            }
            if item.SequelID.Valid && item.SequelID.String != "" {
                if sRow, exists := blogMap[item.SequelID.String]; exists {
                    item.Response.Sequel = &sRow.Response
                }
            }
            continue
        }

        // Forward timeline chain (Depth > 0): Only link sequels to prevent backward loops
        if item.Depth > 0 {
            if item.SequelID.Valid && item.SequelID.String != "" {
                if sRow, exists := blogMap[item.SequelID.String]; exists {
                    item.Response.Sequel = &sRow.Response
                }
            }
        }

        // Backward timeline chain (Depth < 0): Only link prequels to prevent forward loops
        if item.Depth < 0 {
            if item.PrequelID.Valid && item.PrequelID.String != "" {
                if pRow, exists := blogMap[item.PrequelID.String]; exists {
                    item.Response.Prequel = &pRow.Response
                }
            }
        }
    }

    // Extract the requested target post
    var targetBlog, exists = blogMap[targetID]
    if !exists {
        return nil, sql.ErrNoRows
    }

    return &targetBlog.Response, nil
}

