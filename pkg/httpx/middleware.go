package httpx

import (
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type MiddlewareChain struct {
	middlewares []func(http.Handler) http.Handler
}

func NewMiddlewareChain(mws ...func(http.Handler) http.Handler) *MiddlewareChain {
	return &MiddlewareChain{mws}
}

// Handle will apply embeded middlewares
// to given handler
func (c *MiddlewareChain) Handle(h http.Handler) http.Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}
	return h
}

func (c *MiddlewareChain) Append(mws ...func(http.Handler) http.Handler) {
	c.middlewares = append(c.middlewares, mws...)
}

func (c *MiddlewareChain) Merge(chain *MiddlewareChain) {
	c.middlewares = append(c.middlewares, chain.middlewares...)
}

// MakeDevMiddlewares create default middlewares
// for development env
func MakeDevMiddlewares() *MiddlewareChain {
	return NewMiddlewareChain(Recovered, cors.AllowAll().Handler, Logger)
}

func Recovered(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Any("panic", r).
					Msg("recovered from panic")
				WriteInternalErrResp(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		//TODO: status, response raw json
		log.Info().
			Str("host", r.Host).
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Dur("latency", time.Since(start)).
			Msg("http request")
	})
}
