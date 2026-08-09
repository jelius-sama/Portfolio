// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package renderer

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "git.jelius.dev/jelius-sama/Portfolio/api/analytics"
    "git.jelius.dev/jelius-sama/Portfolio/api/blogs"
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
    "golang.org/x/sys/unix"
)

type basicFileStat struct {
    UpdatedAt   time.Time
    PublishedAt time.Time
}

func getFileCreationTime(filePath string, stats *basicFileStat) error {
    if s, err := os.Stat(filePath); err != nil {
        return err
    } else {
        stats.UpdatedAt = s.ModTime()
    }

    var statx unix.Statx_t

    if err := unix.Statx(
        unix.AT_FDCWD,
        filePath,
        unix.AT_SYMLINK_NOFOLLOW,
        unix.STATX_BTIME,
        &statx,
    ); err != nil {
        return fmt.Errorf("statx syscall failed: %w", err)
    }

    if statx.Mask&unix.STATX_BTIME == 0 {
        return fmt.Errorf("birth time (btime) is not supported by this file system")
    }

    stats.PublishedAt = time.Unix(statx.Btime.Sec, int64(statx.Btime.Nsec))
    return nil
}

func (v *ViewManager) RenderAchievements(c fiber.Ctx) error {
    if metadata, err := GetMetadata(c); err != nil {
        logger.Error(c.Path(), err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    } else {
        var stream io.ReadCloser = io.NopCloser(strings.NewReader("achievements"))
        if err := blogs.GetBlogMarkdown(nil, &stream); err != nil {
            logger.Error("Failed to fetch markdown content:", err.Error())
            return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
        }
        defer stream.Close()

        var viewCountBuf strings.Builder
        viewCountBuf.WriteString(c.Path())
        if err := analytics.GetPageVisitCount(nil, &viewCountBuf); err != nil {
            logger.Error(c.Path(), err)
            viewCountBuf.Reset()
            viewCountBuf.WriteString("0")
        }

        var views, _ = strconv.Atoi(viewCountBuf.String())
        var filePath = filepath.Join(types.EVDataDir.Get().Value, "blogs", "achievements.md")
        var stats basicFileStat

        if err := getFileCreationTime(filePath, &stats); err != nil {
            logger.Error(c.Path(), err)
            stats.UpdatedAt = time.Now()
            stats.PublishedAt = time.Now()
        }

        var achievementMetadata types.BlogResponse = types.BlogResponse{
            ID:          "achievements",
            PublishedAt: stats.PublishedAt,
            UpdatedAt:   stats.UpdatedAt,
            Prequel:     nil,
            Sequel:      nil,
            Title:       "My Achievements",
            Excerpt:     "A chronological record of my academic, professional, and personal achievements.\nThis page will be continuously updated as I progress through my journey.",
            Views:       uint(views),
        }

        var markdownContent string
        if content, err := io.ReadAll(stream); err != nil {
            logger.Error("Failed to read markdown content:", err.Error())
            return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
        } else {
            markdownContent = string(content)
            c.Locals("context", "achievement")
            return Renderer(c, metadata, pages.BlogPost(c, &achievementMetadata, &markdownContent))
        }

    }
}

