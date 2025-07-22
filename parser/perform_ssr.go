package parser

import (
	"strings"
)

var DynamicRoutes []string = []string{"/blog/"}

type RouteData struct {
	FullPath string
	Route    string
}

func PerformSSR(fullPath string) (string, error) {
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
		ssrData, err := SSRBlogPage(route.FullPath)
		return ssrData, err

	default:
		return "", nil
	}
}
