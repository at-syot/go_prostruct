package httpx

import "net/http"

// ##### grouping routes with wrapped middlewares
type Route struct {
	Pattern string
	Handler http.HandlerFunc
}

// RegisterRoutes does register route(s) with given
// wrapped middleware(s) to *http.ServeMux
func RegisterRoutes(s *http.ServeMux, composer *MiddlewareChain, routes []Route) {
	for _, r := range routes {
		s.Handle(r.Pattern, composer.Handle(r.Handler))
	}
}
