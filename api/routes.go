package api

import (
	"KazuFolio/api/handler"
	"net/http"
)

var ApiRoutes = map[string]http.HandlerFunc{
	"GET /healthz":        Cacheable(300, handler.Healthz),
	"GET /latest_commit":  handler.LatestCommit,
	"GET /neofetch":       Cacheable(3600, handler.NeofetchInfo),
	"POST /sudo":          handler.VerifySudo,
	"POST /update_server": handler.UpdateServer,
	"GET /version":        handler.GetVersion,
	"POST /analytics":     handler.SaveAnalytics,
	"GET /blogs":          Cacheable(600, handler.GetBlogs),
	"GET /blog":           Cacheable(600, handler.GetBlog),
	"POST /blog":          handler.PostBlog,
	"GET /blog_file":      Cacheable(600, handler.GetBlogFile),
	"GET /analytics":      handler.GetAnalytics,
	"POST /authenticate":  handler.Authenticate,
	"GET /verify_auth":    handler.VerifyAuthStatus,
}
