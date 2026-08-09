// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package blogs

import (
    "crypto/rand"
    "crypto/sha1"
    "encoding/hex"
    "os"
    "path/filepath"

    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

// CreateBlog inserts a new blog post into the database and saves the markdown file
func CreateBlog(c fiber.Ctx) error {
    var req types.CreateBlogPost
    if err := c.Bind().Body(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "Invalid request body",
        })
    }

    // Validate required fields
    if len(req.Title) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "Title is required",
        })
    }

    // Get markdown file from form
    var file, err = c.FormFile("markdown")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "Markdown file is required",
        })
    }

    var randomBytes = make([]byte, 16)
    if _, err := rand.Read(randomBytes); err != nil {
        logger.Error(err)
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }

    var hash = sha1.Sum(randomBytes)
    var id = hex.EncodeToString(hash[:])[:7]

    var query = `
        INSERT INTO blogs (id, title, excerpt, prequel_id, sequel_id, published_at, updated_at)
        VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
    `

    if _, err := db.DB.Exec(query, id, req.Title, req.Excerpt, req.PrequelID, req.SequelID); err != nil {
        logger.Error(c.Path(), err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Failed to create blog",
        })
    }

    // Create blogs directory if it doesn't exist
    var blogsDir = filepath.Join(types.EVDataDir.Get().Value, "blogs")
    if err := os.MkdirAll(blogsDir, 0o755); err != nil {
        logger.Error(c.Path(), err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Failed to create blog directory",
        })
    }

    // Save markdown file
    var filePath = filepath.Join(blogsDir, id+".md")
    if err := c.SaveFile(file, filePath); err != nil {
        logger.Error(c.Path(), err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Failed to save markdown file",
        })
    }

    return c.SendStatus(fiber.StatusCreated)
}

