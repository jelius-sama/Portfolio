package parser

import (
	"net/http"
	"strings"
)

var DynamicRoutes []string = []string{"/blog/", "/achievements"}

type RouteData struct {
	FullPath string
	Route    string
}

func PerformSSR(fullPath string) (string, error, int) {
	var route RouteData

	for _, DynamicRoute := range DynamicRoutes {
		if strings.HasPrefix(fullPath, DynamicRoute) {
			route = RouteData{
				FullPath: fullPath,
				Route:    DynamicRoute,
			}
			break
		}
	}

	switch route.Route {
	case "/blog/":
		ssrData, err, status := SSRBlogPage(route.FullPath)
		return ssrData, err, status

	case "/achievements":
		ssrData, err, status := SSRAchievementsPage(route.FullPath)
		return ssrData, err, status

	default:
		return "", nil, http.StatusNoContent
	}
}
