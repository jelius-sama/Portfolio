package middleware

import (
    "fmt"
    "time"

    "github.com/gofiber/fiber/v3"
)

// CacheConfig defines the configuration for the middleware
type CacheConfig struct {
    // MaxAge defines the duration the resource should be cached
    MaxAge time.Duration
    // Like max-age but only for shared caches (CDNs). Overrides max-age for them.
    SMaxAge *time.Duration
    // Public determines if the resource can be cached by public caches (CDNs/Proxies)
    Public bool
    // MustRevalidate forces caches to submit a request to the origin server for validation
    MustRevalidate bool
    // CustomHeader allows overriding the entire string manually if needed
    CustomHeader string
}

func NewCacheControl(config CacheConfig) fiber.Handler {
    var headerValue string

    if len(config.CustomHeader) != 0 {
        headerValue = config.CustomHeader
    } else {
        var cacheType = "private"
        if config.Public {
            cacheType = "public"
        }

        var maxAgeSeconds = int(config.MaxAge.Seconds())
        headerValue = fmt.Sprintf("%s, max-age=%d", cacheType, maxAgeSeconds)

        if sMaxAge := config.SMaxAge; sMaxAge != nil {
            var maxSAgeSeconds = int(sMaxAge.Seconds())
            headerValue += fmt.Sprintf(", s-maxage=%d", maxSAgeSeconds)
        }

        if config.MustRevalidate {
            headerValue += ", must-revalidate"
        }
    }

    return func(c fiber.Ctx) error {
        // Set the header on the response context
        c.Set("Cache-Control", headerValue)

        return c.Next()
    }
}

