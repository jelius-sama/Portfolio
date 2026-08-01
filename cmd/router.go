package main

import (
    "io/fs"
    "net/url"
    "os"
    "path/filepath"
    "time"

    "git.jelius.dev/jelius-sama/Portfolio/api"
    "git.jelius.dev/jelius-sama/Portfolio/middleware"
    "git.jelius.dev/jelius-sama/Portfolio/renderer"
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
        CustomHeader: "no-store, no-cache, must-revalidate, proxy-revalidate",
    })

    types.Pages = map[string]types.Page{
        "/":      types.Page{Handler: routerCtx.MiddlewareHandlers[types.MHNoCache], Handlers: []any{routerCtx.UI.RenderHome}},
        "/links": types.Page{Handler: routerCtx.MiddlewareHandlers[types.MHNoCache], Handlers: []any{routerCtx.UI.RenderLinks}},
        // TODO:
        "/acheivements": types.Page{Handler: routerCtx.MiddlewareHandlers[types.MHNoCache], Handlers: []any{routerCtx.UI.RenderLinks}},
        "/blogs":        types.Page{Handler: routerCtx.MiddlewareHandlers[types.MHNoCache], Handlers: []any{routerCtx.UI.RenderLinks}},
    }
}

func Router(app *fiber.App) {
    app.Use(middleware.RecoveryMiddleware())
    app.Use(middleware.RequestLogger())

    var apiHandle = app.Group("/api", routerCtx.MiddlewareHandlers[types.MHNoCache])

    apiHandle.Get("/healthz", api.Healthz)

    if types.EVEnv.Get().Value == types.EMDev.String() {
        if _, err := os.Stat("assets"); os.IsNotExist(err) {
            logger.Panic("assets dir missing; golang's lack of #ifdef forces runtime hacks just to stop assets leaking into prod")
        }
        app.Get("/assets/*", static.New("assets", static.Config{
            CacheDuration: 0,
            Compress:      true,
            NotFoundHandler: func(c fiber.Ctx) error {
                return middleware.ErrHandler(c, fiber.ErrNotFound)
            },
        }))
    } else {
        var parsedURL, err = url.Parse(types.EVHostname.Get().Value)
        if err != nil {
            logger.Panic(err)
        }

        var assetDir = filepath.Join(os.Getenv("XDG_DATA_HOME"), parsedURL.Hostname())

        if _, err := os.Stat(assetDir); err != nil {
            logger.Error("Assets directory could not be found, disabling static asset route!")
        } else {
            app.Get("/assets/*", static.New(assetDir, static.Config{
                Browse:        false,
                MaxAge:        3600,
                CacheDuration: 10 * time.Second,
                Compress:      true,
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

