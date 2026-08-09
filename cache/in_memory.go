// Package cache provides a short-lived, in-process response cache for Fiber v3
// route handlers. It exists to solve two things at once:
//
//  1. Avoid re-rendering the same page for a burst of concurrent identical
//     requests (singleflight-style coalescing — only the first request in a
//     window actually renders; everyone else waits on that result).
//  2. Avoid the exact bug the CDN caching layer had: a cache key that
//     ignores request headers HTMX uses to distinguish a full-page request
//     from a partial/SPA-swap request. Any header your handler branches on
//     MUST be included in the cache key, or two requests that should return
//     different content will collide into one cache entry.
//
// Usage:
//
//	blogCache := cache.New(5*time.Minute, "HX-Request", "HX-Target", "HX-Current-URL", "HX-Boosted")
//	routerCtx.MiddlewareHandlers[types.MHShortCache] = blogCache.Middleware()
//
//	types.Pages["/blogs"] = types.Page{
//	    Handler:  routerCtx.MiddlewareHandlers[types.MHShortCache],
//	    Handlers: []any{routerCtx.UI.RenderBlogs},
//	}
//
// Each route that needs caching should generally get its own Store (or at
// least its own header set) rather than sharing one global Store across
// unrelated routes, since the "right" header set can differ per route.
package cache

import (
    "crypto/sha256"
    "encoding/hex"
    "sync"
    "time"

    "github.com/gofiber/fiber/v3"
)

// entry is one cached response.
type entry struct {
    body        []byte
    contentType string
    status      int
    expiresAt   time.Time
}

// call tracks a single in-flight render for a given cache key, so that
// concurrent requests for the same (path + relevant headers) combination
// only trigger one actual render — everyone else waits on the result
// instead of rendering the component redundantly.
type call struct {
    done chan struct{}
    val  *entry
    err  error
}

// Store is a single route-scoped (or app-scoped) short-lived cache.
// Safe for concurrent use.
type Store struct {
    mu         sync.RWMutex
    entries    map[string]*entry
    inflight   sync.Map // key(string) -> *call
    ttl        time.Duration
    headerKeys []string // request headers that must be part of the cache key
}

// New creates a Store with the given TTL. headerKeys should list every
// request header your handlers (or DetermineRenderMode) branch on —
// HX-Request, HX-Target, etc. Anything NOT listed here is invisible to the
// cache key, so two requests that differ only in an unlisted header will
// incorrectly share a cache entry. This is the exact bug that bit the CDN
// setup — don't repeat it here.
func New(ttl time.Duration, headerKeys ...string) *Store {
    s := &Store{
        entries:    make(map[string]*entry),
        ttl:        ttl,
        headerKeys: headerKeys,
    }
    go s.janitor()
    return s
}

// janitor periodically evicts expired entries so the map doesn't grow
// forever.
func (s *Store) janitor() {
    t := time.NewTicker(30 * time.Minute)
    defer t.Stop()
    for range t.C {
        now := time.Now()
        s.mu.Lock()
        for k, e := range s.entries {
            if now.After(e.expiresAt) {
                delete(s.entries, k)
            }
        }
        s.mu.Unlock()
    }
}

// key builds the cache key from method + full URL (path + query) + the
// configured header set. Hashed only to keep map keys a fixed, small size —
// there's no security requirement here, sha256 is just a convenient way to
// fold an arbitrary number of header values into one string.
func (s *Store) key(c fiber.Ctx) string {
    h := sha256.New()
    h.Write([]byte(c.Method()))
    h.Write([]byte{'|'})
    h.Write([]byte(c.OriginalURL()))
    for _, hk := range s.headerKeys {
        h.Write([]byte{'|'})
        h.Write([]byte(hk))
        h.Write([]byte{'='})
        h.Write([]byte(c.Get(hk)))
    }
    return hex.EncodeToString(h.Sum(nil))
}

func (s *Store) get(key string) *entry {
    s.mu.RLock()
    defer s.mu.RUnlock()
    e, ok := s.entries[key]
    if !ok || time.Now().After(e.expiresAt) {
        return nil
    }
    return e
}

func (s *Store) set(key string, e *entry) {
    s.mu.Lock()
    s.entries[key] = e
    s.mu.Unlock()
}

// writeEntry sends a cached (or freshly rendered) entry to the client and
// makes sure downstream proxies/browsers don't ALSO cache it — this cache
// is meant to live only inside your server process, not leak into
// Cloudflare or the browser and reintroduce the original bug at a
// different layer.
func writeEntry(c fiber.Ctx, e *entry, cacheStatus string) error {
    c.Set(fiber.HeaderContentType, e.contentType)
    c.Set(fiber.HeaderCacheControl, "no-store")
    c.Set("X-Cache", cacheStatus)
    return c.Status(e.status).Send(e.body)
}

// Middleware returns a fiber.Handler you can drop in as a route's leading
// handler (the `Handler` field in your types.Page struct, in place of
// MHNoCache). Non-GET requests always pass straight through uncached.
func (s *Store) Middleware() fiber.Handler {
    return func(c fiber.Ctx) error {
        if c.Method() != fiber.MethodGet {
            return c.Next()
        }

        key := s.key(c)

        if e := s.get(key); e != nil {
            return writeEntry(c, e, "HIT")
        }

        // Singleflight: only the first request for a given key actually
        // renders. Every concurrent request for the same key blocks on
        // <-cl.done and reuses that result instead of rendering again.
        leaderCall := &call{done: make(chan struct{})}
        actual, loaded := s.inflight.LoadOrStore(key, leaderCall)
        cl := actual.(*call)

        if loaded {
            <-cl.done
            if cl.err != nil {
                return cl.err
            }
            return writeEntry(c, cl.val, "HIT")
        }

        defer func() {
            s.inflight.Delete(key)
            close(cl.done)
        }()

        if err := c.Next(); err != nil {
            cl.err = err
            return err
        }

        resp := c.Response()
        e := &entry{
            // resp.Body() returns fasthttp's internal buffer, which gets
            // reused after this request completes — it MUST be copied
            // before caching, or the cached bytes will get corrupted by a
            // later, unrelated request.
            body:        append([]byte(nil), resp.Body()...),
            contentType: string(resp.Header.ContentType()),
            status:      resp.StatusCode(),
            expiresAt:   time.Now().Add(s.ttl),
        }
        cl.val = e

        // Don't cache error responses — a transient 500 shouldn't get
        // pinned for 5 minutes. The leader (and anyone already waiting on
        // it) still gets this exact response; it's just not stored for
        // future requests.
        if e.status < fiber.StatusBadRequest {
            s.set(key, e)
            c.Set("X-Cache", "MISS")
        } else {
            c.Set("X-Cache", "SKIP")
        }

        c.Set(fiber.HeaderCacheControl, "no-store")
        return nil
    }
}

