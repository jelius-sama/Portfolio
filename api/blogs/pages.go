package blogs

import (
    "strconv"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "golang.org/x/text/cases"
    "golang.org/x/text/language"
)

func GetBlogs(c fiber.Ctx) error {
    var pageStr = c.Query("page", "1")
    var sort = c.Query("sort", "newest")
    sort = strings.ReplaceAll(sort, "-", " ")
    var caser = cases.Title(language.English)
    sort = caser.String(sort)

    var page int
    if p, err := strconv.Atoi(pageStr); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadGateway,
            Message: "Invalid page number!",
        })
    } else {
        page = p
    }

    if sort != "Newest" && sort != "Oldest" && sort != "Most Viewed" {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadGateway,
            Message: "Invalid sorting order requested!",
        })
    }

    var buf strings.Builder

    if err := pages.BlogPostsOOB(pages.SampleBlogPosts[5:]).Render(c.RequestCtx(), &buf); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: err.Error(),
        })
    }

    if err := pages.BlogInfoOOBUpdate(sort, 10, 10, page, 2).Render(c.RequestCtx(), &buf); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: err.Error(),
        })
    }

    if err := pages.BlogEndOfPosts().Render(c.RequestCtx(), &buf); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: err.Error(),
        })
    }

    c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
    return c.SendString(buf.String())
}

