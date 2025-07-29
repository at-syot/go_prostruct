package httpx

import "net/http"

type MiddlewareChain struct {
	middlewares []func(http.Handler) http.Handler
}

func NewMiddlewares(mws ...func(http.Handler) http.Handler) *MiddlewareChain {
	return &MiddlewareChain{mws}
}

func (c *MiddlewareChain) Handle(h http.Handler) http.Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}
	return h
}

func (c *MiddlewareChain) Append(mws ...func(http.Handler) http.Handler) {
	c.middlewares = append(c.middlewares, mws...)
}

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
