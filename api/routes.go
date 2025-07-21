package api

import (
	"KazuFolio/api/handler"
	"net/http"
)

var ApiRoutes = map[string]http.HandlerFunc{
	"GET /latest_commit":  handler.LatestCommit,
	"GET /neofetch":       handler.NeofetchInfo,
	"POST /sudo":          handler.VerifySudo,
	"POST /update_server": handler.UpdateServer,
	"GET /version":        handler.GetVersion,
	"POST /analytics":     handler.SaveAnalytics,
	"GET /blogs":          handler.GetBlogs,
	"Get /blog":           handler.GetBlog,
	"POST /blog":          handler.PostBlog,
}
