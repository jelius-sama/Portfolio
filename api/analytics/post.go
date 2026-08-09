// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package analytics

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "sync"

    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

type TZResolver struct {
    mu    sync.RWMutex
    tzMap map[string]string
}

var resolver *TZResolver

func init() {
    if r, err := NewTZResolver(); err != nil {
        logger.Panic(err)
    } else {
        resolver = r
    }
}

// NewTZResolver reads the Linux zone.tab file and loads it into memory
func NewTZResolver() (*TZResolver, error) {
    // zone.tab is the legacy standard, zone1970.tab is the newer standard.
    // We check zone.tab first as it's universally present on Linux systems.
    filePath := "/usr/share/zoneinfo/zone.tab"

    file, err := os.Open(filePath)
    if err != nil {
        // Fallback to zone1970.tab if zone.tab doesn't exist
        filePath = "/usr/share/zoneinfo/zone1970.tab"
        file, err = os.Open(filePath)
        if err != nil {
            return nil, fmt.Errorf("failed to open timezone data files: %w", err)
        }
    }
    defer file.Close()

    tzMap := make(map[string]string)
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()

        // Skip comments and empty lines
        if len(line) == 0 || line[0] == '#' {
            continue
        }

        // zone.tab columns are separated by tabs: CountryCode \t Coordinates \t TZ_Name
        fields := strings.Split(line, "\t")
        if len(fields) < 3 {
            continue
        }

        countryCode := fields[0]
        tzName := fields[2]

        // Map the timezone string to the country code
        tzMap[tzName] = countryCode
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error reading timezone file: %w", err)
    }

    return &TZResolver{tzMap: tzMap}, nil
}

// GetCountryCode looks up the timezone string in the loaded map
func (r *TZResolver) GetCountryCode(tz string) string {
    r.mu.RLock()
    defer r.mu.RUnlock()

    if code, exists := r.tzMap[tz]; exists {
        return code
    }
    return "UNKNOWN"
}

func TrackAnalytics(c fiber.Ctx) error {
    var req types.TrackAnalyticsRequest

    // Parse request body
    if err := c.Bind().JSON(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "Invalid request body",
        })
    }

    // Validate required fields
    if len(req.UserTimeZone) == 0 || len(req.PagePath) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
            Code:    fiber.StatusBadRequest,
            Message: "user_time_zone and page_path are required",
        })
    }

    // Insert into database
    var query = `
        INSERT INTO analytics_events (country_code, page_path, timestamp)
        VALUES (?, ?, datetime('now', 'utc'))
    `

    if _, err := db.DB.Exec(query, resolver.GetCountryCode(req.UserTimeZone), req.PagePath); err != nil {
        logger.Error(c.Path(), err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResp{
            Code:    fiber.StatusInternalServerError,
            Message: "Internal Server Error",
        })
    }

    return c.SendStatus(fiber.StatusCreated)
}

