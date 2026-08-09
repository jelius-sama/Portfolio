// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package renderer

import (
    "context"
    "database/sql"
    "fmt"

    "git.jelius.dev/jelius-sama/Portfolio/db"
    "git.jelius.dev/jelius-sama/Portfolio/template/components"
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"

    "github.com/gofiber/fiber/v3"
    "github.com/jelius-sama/logger"
)

func getHomePageData(ctx context.Context) (pages.HomeData, error) {
    var homeData pages.HomeData
    var homePageID int

    var pageQuery = `
        SELECT
            id, cv_path, image_path, image_alt_text, first_name, last_name,
            education_level, study_domain, who_am_i, email, qr_image_path
        FROM home_page
        ORDER BY updated_at DESC
        LIMIT 1
    `

    if err := db.DB.QueryRowContext(ctx, pageQuery).Scan(
        &homePageID,
        &homeData.HeroSectionData.CVPath,
        &homeData.HeroSectionData.ImagePath,
        &homeData.HeroSectionData.ImageAltText,
        &homeData.HeroSectionData.FirstName,
        &homeData.HeroSectionData.LastName,
        &homeData.HeroSectionData.EducationLevel,
        &homeData.HeroSectionData.StudyDomain,
        &homeData.HeroSectionData.WhoAmI,
        &homeData.ContactSectionData.Email,
        &homeData.ContactSectionData.QRImagePath,
    ); err != nil {
        if err == sql.ErrNoRows {
            return pages.HomeData{}, fmt.Errorf("home page not found")
        }
        return pages.HomeData{}, err
    }

    // Specialized domains
    if domainRows, err := db.DB.QueryContext(ctx,
        `SELECT domain FROM home_specialized_domain WHERE home_page_id = ? ORDER BY position ASC`,
        homePageID,
    ); err != nil {
        return pages.HomeData{}, err
    } else {
        defer domainRows.Close()

        for domainRows.Next() {
            var domain string
            if err := domainRows.Scan(&domain); err != nil {
                return pages.HomeData{}, err
            }
            homeData.HeroSectionData.SpecializedDomain = append(homeData.HeroSectionData.SpecializedDomain, domain)
        }
        if err = domainRows.Err(); err != nil {
            return pages.HomeData{}, err
        }
    }

    // About messages
    if messageRows, err := db.DB.QueryContext(ctx,
        `SELECT message FROM home_about_message WHERE home_page_id = ? ORDER BY position ASC`,
        homePageID,
    ); err != nil {
        return pages.HomeData{}, err
    } else {
        defer messageRows.Close()

        for messageRows.Next() {
            var message string
            if err := messageRows.Scan(&message); err != nil {
                return pages.HomeData{}, err
            }
            homeData.AboutSectionMessage = append(homeData.AboutSectionMessage, message)
        }
        if err = messageRows.Err(); err != nil {
            return pages.HomeData{}, err
        }
    }

    // Skills
    if skillRows, err := db.DB.QueryContext(ctx,
        `SELECT category, skill FROM home_skill WHERE home_page_id = ? ORDER BY category ASC, position ASC`,
        homePageID,
    ); err != nil {
        return pages.HomeData{}, err
    } else {
        defer skillRows.Close()

        for skillRows.Next() {
            var category, skill string
            if err := skillRows.Scan(&category, &skill); err != nil {
                return pages.HomeData{}, err
            }
            switch category {
            case "language":
                homeData.SkillSectionData.Languages = append(homeData.SkillSectionData.Languages, skill)
            case "framework":
                homeData.SkillSectionData.Frameworks = append(homeData.SkillSectionData.Frameworks, skill)
            case "other":
                homeData.SkillSectionData.Others = append(homeData.SkillSectionData.Others, skill)
            }
        }
        if err = skillRows.Err(); err != nil {
            return pages.HomeData{}, err
        }
    }

    // Projects
    if projectRows, err := db.DB.QueryContext(ctx,
        `SELECT id, title, description, thumbnail FROM home_project WHERE home_page_id = ? ORDER BY position ASC`,
        homePageID,
    ); err != nil {
        return pages.HomeData{}, err
    } else {
        defer projectRows.Close()

        for projectRows.Next() {
            var projectID int
            var project components.ProjectData

            if err := projectRows.Scan(&projectID, &project.Title, &project.Description, &project.Thumbnail); err != nil {
                return pages.HomeData{}, err
            }

            if tagRows, err := db.DB.QueryContext(ctx,
                `SELECT tag FROM home_project_tag WHERE project_id = ? ORDER BY position ASC`,
                projectID,
            ); err != nil {
                return pages.HomeData{}, err
            } else {
                defer tagRows.Close()

                for tagRows.Next() {
                    var tag string
                    if err := tagRows.Scan(&tag); err != nil {
                        return pages.HomeData{}, err
                    }
                    project.Tags = append(project.Tags, tag)
                }
                if err = tagRows.Err(); err != nil {
                    return pages.HomeData{}, err
                }
            }

            if linkRows, err := db.DB.QueryContext(ctx,
                `SELECT label, href, kind FROM home_project_link WHERE project_id = ? ORDER BY position ASC`,
                projectID,
            ); err != nil {
                return pages.HomeData{}, err
            } else {
                defer linkRows.Close()

                for linkRows.Next() {
                    var link components.ProjectLink
                    if err := linkRows.Scan(&link.Label, &link.Href, &link.Kind); err != nil {
                        return pages.HomeData{}, err
                    }
                    project.Links = append(project.Links, link)
                }
                if err = linkRows.Err(); err != nil {
                    return pages.HomeData{}, err
                }
            }

            homeData.ProjectSectionData = append(homeData.ProjectSectionData, project)
        }
        if err = projectRows.Err(); err != nil {
            return pages.HomeData{}, err
        }
    }

    // Experience
    if experienceRows, err := db.DB.QueryContext(ctx,
        `SELECT name, description, date_range, is_active FROM home_experience WHERE home_page_id = ? ORDER BY position ASC`,
        homePageID,
    ); err != nil {
        return pages.HomeData{}, err
    } else {
        defer experienceRows.Close()

        for experienceRows.Next() {
            var entry components.ExperienceEntry
            if err := experienceRows.Scan(&entry.Name, &entry.Description, &entry.DateRange, &entry.IsActive); err != nil {
                return pages.HomeData{}, err
            }
            homeData.ExperienceSectionData = append(homeData.ExperienceSectionData, entry)
        }
        if err = experienceRows.Err(); err != nil {
            return pages.HomeData{}, err
        }
    }

    // Contact quick links
    if quickLinkRows, err := db.DB.QueryContext(ctx,
        `SELECT platform, url, icon FROM home_contact_quick_link WHERE home_page_id = ? ORDER BY position ASC`,
        homePageID,
    ); err != nil {
        return pages.HomeData{}, err
    } else {
        defer quickLinkRows.Close()

        for quickLinkRows.Next() {
            var ql components.QuickContactLinks
            if err := quickLinkRows.Scan(&ql.Platform, &ql.URL, &ql.Icon); err != nil {
                return pages.HomeData{}, err
            }
            homeData.ContactSectionData.QuickLinks = append(homeData.ContactSectionData.QuickLinks, ql)
        }
        if err = quickLinkRows.Err(); err != nil {
            return pages.HomeData{}, err
        }
    }

    return homeData, nil
}

func (v *ViewManager) RenderHome(c fiber.Ctx) error {
    if metadata, err := GetMetadata(c); err != nil {
        logger.Error(c.Path(), err.Error())
        return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
    } else {
        if data, err := getHomePageData(c.RequestCtx()); err == nil {
            return Renderer(c, metadata, pages.Home(data))
        } else {
            logger.Error(c.Path(), err.Error())
            return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
        }
    }
}

