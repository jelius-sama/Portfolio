package main

import (
    "os"
    "path/filepath"
    "time"

    "git.jelius.dev/jelius-sama/Portfolio/api"
    "git.jelius.dev/jelius-sama/Portfolio/api/analytics"
    "git.jelius.dev/jelius-sama/Portfolio/api/blogs"
    "git.jelius.dev/jelius-sama/Portfolio/middleware"
    "git.jelius.dev/jelius-sama/Portfolio/renderer"
    "git.jelius.dev/jelius-sama/Portfolio/template/pages"
    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/static"
    "github.com/jelius-sama/logger"
)

type RouterCtx struct {
    UI                 *renderer.ViewManager
    MiddlewareHandlers types.MiddlewareHandlerMap
}

var routerCtx RouterCtx = RouterCtx{}

func init() {
    routerCtx.UI = renderer.New()
    routerCtx.MiddlewareHandlers = make(types.MiddlewareHandlerMap)
    routerCtx.MiddlewareHandlers[types.MHNoCache] = middleware.NewCacheControl(middleware.CacheConfig{
        CustomHeader: "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate",
    })
    routerCtx.MiddlewareHandlers[types.MHStaticPages] = middleware.NewCacheControl(middleware.CacheConfig{
        Public:         true,
        MaxAge:         time.Hour,
        SMaxAge:        new(48 * time.Hour),
        MustRevalidate: true,
    })
    routerCtx.MiddlewareHandlers[types.MHStaticAsset] = middleware.NewCacheControl(middleware.CacheConfig{
        Public:         true,
        MaxAge:         31536000 * time.Second, // 1 year
        MustRevalidate: true,
    })

    var pageCache = routerCtx.MiddlewareHandlers[types.MHStaticPages]
    if types.EVEnv.Get().Value == types.EMDev.String() {
        pageCache = routerCtx.MiddlewareHandlers[types.MHNoCache]
    }

    var inDev = func(c fiber.Ctx) error {
        return renderer.Renderer(c, new(types.Metadata{
            Title:       "501",
            Description: "Feature in development",
        }), pages.Error(501, &fiber.Error{Code: 501, Message: "Active work on this feature is ongoing, in the mean time you can visit other parts of this application."}))
    }

    types.Pages = map[string]types.Page{
        "/":            types.Page{Handler: pageCache, Handlers: []any{routerCtx.UI.RenderHome}},
        "/links":       types.Page{Handler: pageCache, Handlers: []any{routerCtx.UI.RenderLinks}},
        "/blogs":       types.Page{Handler: pageCache, Handlers: []any{routerCtx.UI.RenderBlogs}},
        "/robots.txt":  types.Page{Handler: routerCtx.MiddlewareHandlers[types.MHNoCache], Handlers: []any{api.GenerateRobots}},
        "/sitemap.xml": types.Page{Handler: routerCtx.MiddlewareHandlers[types.MHNoCache], Handlers: []any{api.GenerateSitemap}},
        // TODO: Implement these pages
        // TODO: Implement dynamic metadata for blog pages
        "/blog/:id":     types.Page{Handler: pageCache, Handlers: []any{inDev}},
        "/acheivements": types.Page{Handler: pageCache, Handlers: []any{inDev}},
    }
}

func Router(app *fiber.App) {
    app.Use(middleware.RecoveryMiddleware())
    app.Use(middleware.RequestLogger())

    var apiHandle = app.Group("/api", routerCtx.MiddlewareHandlers[types.MHNoCache])

    apiHandle.Get("/healthz", api.Healthz)
    apiHandle.Get("/version", api.Version)

    // Analytics endpoints
    apiHandle.Get("/analytics/get/all", analytics.GetAllAnalyticsEvents)
    apiHandle.Get("/analytics/get/visit-count", func(c fiber.Ctx) error { return analytics.GetPageVisitCount(c) })
    apiHandle.Get("/analytics/get/avg-visits", analytics.GetAvgVisitsPerHour)
    apiHandle.Get("/analytics/get/top-countries", analytics.GetTopCountries)
    apiHandle.Get("/analytics/get/top-pages", analytics.GetTopPages)
    apiHandle.Post("/analytics/track", analytics.TrackAnalytics)

    apiHandle.Get("/blogs", blogs.GetBlogsPage)

    apiHandle.Get("/blog/all", func(c fiber.Ctx) error { return blogs.GetAllBlogs(c) })
    apiHandle.Get("/blog/md/:id", func(c fiber.Ctx) error { return blogs.GetBlogMarkdown(c) })
    apiHandle.Get("/blog/:id", func(c fiber.Ctx) error { return blogs.GetBlog(c) })
    apiHandle.Post("/blog", blogs.CreateBlog)

    if types.EVEnv.Get().Value == types.EMDev.String() {
        if _, err := os.Stat("assets"); os.IsNotExist(err) {
            logger.Panic("assets dir missing; golang's lack of #ifdef forces runtime hacks just to stop assets leaking into prod")
        }
        app.Get("/assets/*", routerCtx.MiddlewareHandlers[types.MHNoCache], static.New("assets", static.Config{
            Compress:      true,
            CacheDuration: 0 * time.Second,
            NotFoundHandler: func(c fiber.Ctx) error {
                return middleware.ErrHandler(c, fiber.ErrNotFound)
            },
        }))
    } else {
        var assetDir = filepath.Join(types.EVDataDir.Get().Value, "assets")

        if _, err := os.Stat(assetDir); err != nil {
            logger.Error("Assets directory could not be found, disabling static asset route!")
        } else {
            logger.Okay("Serving assets from:", assetDir)
            app.Get("/assets/*", routerCtx.MiddlewareHandlers[types.MHStaticAsset], static.New(assetDir, static.Config{
                Browse:   false,
                Compress: true,
                NotFoundHandler: func(c fiber.Ctx) error {
                    return middleware.ErrHandler(c, fiber.ErrNotFound)
                },
            }))
        }
    }

    for k, v := range types.Pages {
        app.Get(k, v.Handler, v.Handlers...)
    }
}

