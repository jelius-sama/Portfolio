package renderer

import (
    "bytes"
    "database/sql"
    "encoding/gob"
    "fmt"
    "io"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/api/blogs"
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func (v *ViewManager) RenderBlog(c fiber.Ctx) error {
    var buf = bytes.NewBufferString(c.Params("id"))
    if err := blogs.GetBlog(nil, buf); err != nil {
        if err == sql.ErrNoRows {
            c.Locals("pseudo_path", "#not_found")
            if metadata, metadataErr := GetMetadata(c); metadataErr != nil {
                logger.Error(c.Path(), err.Error())
                return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
            } else {
                return Renderer(c, metadata, pages.BlogPost(c, nil, nil))
            }
        }
        logger.Error("Failed to fetch blog post data:", err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    var decodedResponse types.BlogResponse
    if decErr := gob.NewDecoder(buf).Decode(&decodedResponse); decErr != nil {
        logger.Error("Failed to decode Gob data:", decErr.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    var stream io.ReadCloser = io.NopCloser(strings.NewReader(c.Params("id")))
    if err := blogs.GetBlogMarkdown(nil, &stream); err != nil {
        logger.Error("Failed to fetch markdown content:", err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }
    defer stream.Close()

    var markdownContent string
    if content, err := io.ReadAll(stream); err != nil {
        logger.Error("Failed to read markdown content:", err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    } else {
        markdownContent = string(content)
        c.Locals("context", "blog")
        c.Locals("pseudo_path", "*")

        if metadata, metadataErr := GetMetadata(c); metadataErr != nil {
            logger.Error(c.Path(), metadataErr.Error())
            return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
        } else {
            c.Locals("title", fmt.Sprintf("%s | Jelius", decodedResponse.Title))
            c.Locals("description", decodedResponse.Excerpt)
            GetDynamicRouteMetadata(c, metadata)
            return Renderer(c, metadata, pages.BlogPost(c, &decodedResponse, &markdownContent))
        }
    }
}

