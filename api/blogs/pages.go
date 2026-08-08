package blogs

import (
    "bytes"
    "encoding/gob"
    "fmt"
    "strconv"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func GetBlogsPage(c fiber.Ctx) error {
    var pageStr = c.Query("page", "1")
    var sort types.BlogsSortOrder = types.BSONew

    if s := c.Query("sort"); len(s) != 0 {
        switch s {
        case types.BSONew.String():
            sort = types.BSONew
        case types.BSOPopular.String():
            sort = types.BSOPopular
        case types.BSOOld.String():
            sort = types.BSOOld
        default:
            return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
                Code:    fiber.StatusBadRequest,
                Message: "Invalid sorting order requested!",
            })
        }
    }

    var page int
    if p, err := strconv.Atoi(pageStr); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "Invalid page number!",
        })
    } else {
        page = p
    }

    var dataBuf bytes.Buffer
    if _, err := fmt.Fprintf(&dataBuf, "page=%dsort=%d", page, sort); err != nil {
        logger.Error(c.Path(), err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    if err := GetAllBlogs(nil, &dataBuf); err != nil {
        logger.Error(c.Path(), err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    var decodedResponse types.PaginatedBlogsResponse
    if err := gob.NewDecoder(&dataBuf).Decode(&decodedResponse); err != nil {
        logger.Error(c.Path(), err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    }

    var buf strings.Builder

    if page == 1 {
        if err := pages.BlogsSection(pages.BlogsSectionArgs{
            Post:        decodedResponse.Data,
            HasMore:     decodedResponse.HasMore,
            TotalPosts:  decodedResponse.TotalRows,
            LoadedPages: decodedResponse.Page,
            TotalPages:  (decodedResponse.TotalRows + decodedResponse.Limit - 1) / decodedResponse.Limit,
            Sort:        decodedResponse.Sort,
        }).Render(c.RequestCtx(), &buf); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: err.Error(),
            })
        }

        c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
        return c.SendString(buf.String())
    }

    if err := pages.BlogPostsOOB(decodedResponse.Data).Render(c.RequestCtx(), &buf); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: err.Error(),
        })
    }

    if err := pages.BlogInfoOOBUpdate(pages.BlogInfoArgs{
        Sort:        sort,
        LoadedPosts: ((page - 1) * decodedResponse.Limit) + len(decodedResponse.Data),
        // NOTE: For mathematicians out there, golang stores the length of an array internally, it is more efficient to use that length than to do fancy math which only waste more CPU cycles
        // LoadedPosts: ((page - 1) * decodedResponse.Limit) + min(decodedResponse.Limit, decodedResponse.TotalRows - ((page - 1) * decodedResponse.Limit)),
        TotalPosts:  decodedResponse.TotalRows,
        LoadedPages: decodedResponse.Page,
        TotalPages:  (decodedResponse.TotalRows + decodedResponse.Limit - 1) / decodedResponse.Limit,
    }).Render(c.RequestCtx(), &buf); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: err.Error(),
        })
    }

    if page >= (decodedResponse.TotalRows+decodedResponse.Limit-1)/decodedResponse.Limit {
        if err := pages.BlogEndOfPosts().Render(c.RequestCtx(), &buf); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: err.Error(),
            })
        }
    } else {
        if err := pages.BlogLoadMoreTrigger(decodedResponse.Page, sort).Render(c.RequestCtx(), &buf); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: err.Error(),
            })
        }
    }

    c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
    return c.SendString(buf.String())
}

