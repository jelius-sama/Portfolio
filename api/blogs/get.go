package blogs

import (
    "bytes"
    "database/sql"
    "encoding/gob"
    "strconv"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/api/analytics"
    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

// GetBlog retrieves a blog post by ID with its prequel and sequel chain
func GetBlog(c fiber.Ctx, buf ...*bytes.Buffer) error {
    var id string
    if len(buf) != 0 {
        id = buf[0].String()
    } else {
        id = c.Params("id")
    }

    if len(id) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "Blog ID is required",
        })
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
        if err := gob.NewEncoder(buf[0]).Encode(blog); err != nil {
            return err
        }
        return nil
    }

    return c.Status(fiber.StatusOK).JSON(blog)
}

// getBlogResponse recursively fetches a blog and its prequel/sequel chain
func getBlogResponse(id string) (*types.BlogResponse, error) {
    var query = `
        SELECT id, title, excerpt, published_at, updated_at, deleted_at, prequel_id, sequel_id
        FROM blogs
        WHERE id = ? AND deleted_at IS NULL
    `
    var blog types.BlogResponse
    var deletedAt sql.NullTime
    var prequelID sql.NullString
    var sequelID sql.NullString

    if err := db.DB.QueryRow(query, id).Scan(
        &blog.ID,
        &blog.Title,
        &blog.Excerpt,
        &blog.PublishedAt,
        &blog.UpdatedAt,
        &deletedAt,
        &prequelID,
        &sequelID,
    ); err != nil {
        return nil, err
    }

    if deletedAt.Valid {
        blog.DeletedAt = &deletedAt.Time
    }

    // Get view count for this blog post
    var viewCountBuf strings.Builder
    viewCountBuf.WriteString("/blog/")
    viewCountBuf.WriteString(blog.ID)

    if err := analytics.GetPageVisitCount(nil, &viewCountBuf); err == nil {
        if views, err := strconv.Atoi(viewCountBuf.String()); err == nil {
            blog.Views = uint(views)
        }
    }

    // Recursively fetch prequel if it exists
    if prequelID.Valid && len(prequelID.String) > 0 {
        var prequel, err = getBlogResponse(prequelID.String)
        if err == nil {
            blog.Prequel = prequel
        }
    }

    // Recursively fetch sequel if it exists
    if sequelID.Valid && len(sequelID.String) > 0 {
        var sequel, err = getBlogResponse(sequelID.String)
        if err == nil {
            blog.Sequel = sequel
        }
    }

    return &blog, nil
}

