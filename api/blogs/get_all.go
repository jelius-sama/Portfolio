package blogs

import (
    "bytes"
    "encoding/gob"
    "fmt"
    "strconv"
    "strings"

    "git.jelius.dev/jelius-sama/Portfolio/api/analytics"
    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

// NOTE:
//	var tempTestDirectCall = func() {
//	    buf := bytes.NewBufferString("page=1sort=0")
//	    err := blogs.GetAllBlogs(nil, buf)
//	    if err == nil {
//	        var decodedResponse types.PaginatedBlogsResponse
//	        decErr := gob.NewDecoder(buf).Decode(&decodedResponse)
//	        if decErr != nil {
//	            logger.Error("Failed to decode Gob data:", decErr.Error())
//	            return
//	        }
//          decodedResponse.Title // access fields
//          // convert into human-readible
//	        humanReadable, _ := json.MarshalIndent(decodedResponse, "", "  ")
//	        logger.Okay("Buffer output:\n", string(humanReadable))
//	    } else {
//	        logger.Error("Error output:", err.Error())
//	    }
//	}

// GetAllBlogs retrieves all non-deleted blog posts with pagination and sorting
func GetAllBlogs(c fiber.Ctx, buf ...*bytes.Buffer) error {
    var page int = 0
    var sort types.BlogsSortOrder = types.BSONew

    var validateSort = func(t any, s *types.BlogsSortOrder) error {
        switch t := t.(type) {
        case string:
            if order, err := strconv.Atoi(t); err != nil {
                return fmt.Errorf("invalid sorting parameter")
            } else if bso := types.BlogsSortOrder(order); bso != types.BSONew && bso != types.BSOOld && bso != types.BSOPopular {
                return fmt.Errorf("invalid sorting parameter")
            } else {
                *s = bso
            }
            return nil
        case types.BlogsSortOrder:
            if t != types.BSONew && t != types.BSOOld && t != types.BSOPopular {
                return fmt.Errorf("invalid sorting parameter")
            } else {
                *s = t
            }
            return nil
        default:
            return fmt.Errorf("unsupported type in validator function")
        }
    }

    var readDigits = func(s string) string {
        var end int
        for end < len(s) && s[end] >= '0' && s[end] <= '9' {
            end++
        }
        return s[:end]
    }

    if has := len(buf) > 0; has && buf[0] == nil {
        return fmt.Errorf("length of variadic greater than 0 but nil buffer\n")
    } else if has {
        var str = buf[0].String()
        if _, after, found := strings.Cut(str, "page="); found {
            if val, err := strconv.Atoi(readDigits(after)); err == nil {
                page = val
            }
        }

        if _, after, found := strings.Cut(str, "sort="); found {
            if val, err := strconv.Atoi(readDigits(after)); err == nil {
                sort = types.BlogsSortOrder(val)
            }
        }

        if err := validateSort(sort, &sort); err != nil {
            return fmt.Errorf("page must be a non-negative integer\n")
        }
    } else {
        if p, pageErr := strconv.Atoi(c.Query("page", "0")); pageErr != nil || p < 0 {
            return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
                Code:    fiber.StatusBadRequest,
                Message: "page must be a non-negative integer",
            })
        } else {
            page = p
        }

        if err := validateSort(c.Query("sort", "0"), &sort); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
                Code:    fiber.StatusBadRequest,
                Message: "sort must be a valid sorting parameter",
            })
        }
    }

    // Get total count of non-deleted blogs
    var countQuery = `SELECT COUNT(*) FROM blogs WHERE deleted_at IS NULL`
    var totalRows int
    if err := db.DB.QueryRow(countQuery).Scan(&totalRows); err != nil {
        if len(buf) != 0 {
            return err
        }

        logger.Error(c.Path(), err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }

    // Calculate offset
    var offset = page * types.PostPerPage

    // Build query with sorting
    var orderBy string
    switch sort {
    case types.BSONew:
        orderBy = "published_at DESC"
    case types.BSOOld:
        orderBy = "published_at ASC"
    case types.BSOPopular:
        orderBy = "views DESC"
    }

    var query = `
        SELECT 
            id, title, excerpt, published_at, updated_at, deleted_at, prequel_id, sequel_id
        FROM blogs
        WHERE deleted_at IS NULL
        ORDER BY ` + orderBy + `
        LIMIT ? OFFSET ?
    `

    var rows, queryErr = db.DB.Query(query, types.PostPerPage, offset)
    if queryErr != nil {
        if len(buf) != 0 {
            return queryErr
        }

        logger.Error(c.Path(), queryErr.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }
    defer rows.Close()

    var data []types.BlogPost
    for rows.Next() {
        var post types.BlogPost
        if err := rows.Scan(
            &post.ID,
            &post.Title,
            &post.Excerpt,
            &post.PublishedAt,
            &post.UpdatedAt,
            &post.DeletedAt,
            &post.PrequelID,
            &post.SequelID,
        ); err != nil {
            if len(buf) != 0 {
                return err
            }

            logger.Error(c.Path(), err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: "Internal Server Error",
            })
        }

        // Get view count for this blog post
        var viewCountBuf strings.Builder
        viewCountBuf.WriteString("/blog/")
        viewCountBuf.WriteString(post.ID)
        if err := analytics.GetPageVisitCount(nil, &viewCountBuf); err != nil {
            if len(buf) != 0 {
                return err
            }

            logger.Error(c.Path(), err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
                Code:    fiber.StatusInternalServerError,
                Message: "Internal Server Error",
            })
        }

        views, _ := strconv.Atoi(viewCountBuf.String())
        post.Views = uint(views)

        data = append(data, post)
    }

    var hasMore bool = (offset + types.PostPerPage) < totalRows
    var resp = types.PaginatedBlogsResponse{
        Data:      data,
        Page:      page,
        Limit:     types.PostPerPage,
        HasMore:   hasMore,
        TotalRows: totalRows,
        Sort:      sort.String(),
    }

    if len(buf) != 0 {
        buf[0].Reset()
        return gob.NewEncoder(buf[0]).Encode(resp)
    }

    return c.Status(fiber.StatusOK).JSON(resp)
}

