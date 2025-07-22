package api

import (
	"KazuFolio/api/handler"
	"net/http"
)

var ApiRoutes = map[string]http.HandlerFunc{
	"GET /healthz":        handler.Healthz,
	"GET /latest_commit":  handler.LatestCommit,
	"GET /neofetch":       handler.NeofetchInfo,
	"POST /sudo":          handler.VerifySudo,
	"POST /update_server": handler.UpdateServer,
	"GET /version":        handler.GetVersion,
	"POST /analytics":     handler.SaveAnalytics,
	"GET /blogs":          handler.GetBlogs,
	"GET /blog":           handler.GetBlog,
	"POST /blog":          handler.PostBlog,
	"GET /blog_file":      handler.GetBlogFile,
}
