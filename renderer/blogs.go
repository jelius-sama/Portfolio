package renderer

import (
    "bytes"
    "encoding/gob"

    "git.jelius.dev/jelius-sama/Portfolio/api/blogs"
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func (v *ViewManager) RenderBlogs(c fiber.Ctx) error {
    var metadata types.Metadata = types.Metadata{
        Title:       "Blogs",
        Description: "My blog posts",
    }

    var buf = bytes.NewBufferString("page=1sort=0")
    if err := blogs.GetAllBlogs(nil, buf); err != nil {
        logger.Error(c.Path(), err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    var decodedResponse types.PaginatedBlogsResponse
    if err := gob.NewDecoder(buf).Decode(&decodedResponse); err != nil {
        logger.Error(c.Path(), err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    return Renderer(c, metadata, pages.Blogs(pages.BlogsSectionArgs{
        Post:        decodedResponse.Data,
        HasMore:     decodedResponse.HasMore,
        TotalPosts:  decodedResponse.TotalRows,
        LoadedPages: decodedResponse.Page,
        TotalPages:  max(1, (decodedResponse.TotalRows+decodedResponse.Limit-1)/decodedResponse.Limit),
        Sort:        decodedResponse.Sort,
    }))
}

