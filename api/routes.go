package api

import (
	"KazuFolio/api/handler"
	"net/http"
)

var ApiRoutes = map[string]http.HandlerFunc{
	"GET /latest_commit": handler.LatestCommit,
	"GET /neofetch":      handler.NeofetchInfo,
}
