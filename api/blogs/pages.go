package blogs

import (
    "strconv"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
)

func GetBlogs(c fiber.Ctx) error {
    var pageStr = c.Query("page", "1")
    var sort types.BlogsSortOrder = types.BSONew

    if s := c.Query("sort", ""); len(s) != 0 {
        switch s {
        case types.BSONew.String():
            sort = types.BSONew
        case types.BSOPopular.String():
            sort = types.BSOPopular
        case types.BSOOld.String():
            sort = types.BSOOld
        default:
            return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
                Code:    fiber.StatusBadGateway,
                Message: "Invalid sorting order requested!",
            })
        }
    }

    var page int
    if p, err := strconv.Atoi(pageStr); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadGateway,
            Message: "Invalid page number!",
        })
    } else {
        page = p
    }

    var buf strings.Builder

    if page == 1 {
        if err := pages.BlogsSection(pages.BlogsSectionArgs{
            Post:        types.SampleBlogPosts[:5],
            HasMore:     true,
            TotalPosts:  len(types.SampleBlogPosts),
            LoadedPages: page,
            TotalPages:  (len(types.SampleBlogPosts) + types.PostPerPage - 1) / types.PostPerPage,
            Sort:        sort,
        }).Render(c.RequestCtx(), &buf); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: err.Error(),
            })
        }

        c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
        return c.SendString(buf.String())
    }

    var skip = (page - 1) * types.PostPerPage
    var take = types.PostPerPage

    start := min(skip, len(types.SampleBlogPosts))
    end := min(start+take, len(types.SampleBlogPosts))

    if err := pages.BlogPostsOOB(types.SampleBlogPosts[start:end]).Render(c.RequestCtx(), &buf); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: err.Error(),
        })
    }

    if err := pages.BlogInfoOOBUpdate(pages.BlogInfoArgs{
        Sort:        sort,
        LoadedPosts: ((page - 1) * types.PostPerPage) + len(types.SampleBlogPosts[start:end]),
        TotalPosts:  len(types.SampleBlogPosts),
        LoadedPages: page,
        TotalPages:  (len(types.SampleBlogPosts) + types.PostPerPage - 1) / types.PostPerPage,
    }).Render(c.RequestCtx(), &buf); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: err.Error(),
        })
    }

    if page >= (len(types.SampleBlogPosts)+types.PostPerPage-1)/types.PostPerPage {
        if err := pages.BlogEndOfPosts().Render(c.RequestCtx(), &buf); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: err.Error(),
            })
        }
    } else {
        if err := pages.BlogLoadMoreTrigger(page+1, types.BSONew).Render(c.RequestCtx(), &buf); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: err.Error(),
            })
        }
    }

    c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
    return c.SendString(buf.String())
}

