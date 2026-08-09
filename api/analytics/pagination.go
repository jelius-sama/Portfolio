// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package analytics

import (
    "strconv"

    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
)

func parsePaginationParam(pageStr, limitStr, sortOrder string) (int, int, types.SortOrder, *types.ErrorResp) {
    var page, pageErr = strconv.Atoi(pageStr)
    if pageErr != nil || page < 0 {
        return 0, 0, 0, &types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "page must be a non-negative integer",
        }
    }

    var limit, limitErr = strconv.Atoi(limitStr)
    if limitErr != nil || limit < 1 {
        return 0, 0, 0, &types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "limit must be a positive integer",
        }
    }

    var sort types.SortOrder
    if order, err := strconv.Atoi(sortOrder); err != nil {
        return 0, 0, 0, &types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "invalid sort order",
        }
    } else if types.SortOrder(order) != types.SOAsc && types.SortOrder(order) != types.SODesc {
        return 0, 0, 0, &types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "invalid sort order",
        }
    } else {
        sort = types.SortOrder(order)
    }

    return page, limit, sort, nil
}

