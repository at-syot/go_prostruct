package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/simt/dtacc"
)

func main() {
	db, err := dtacc.NewDB()
	if err != nil {
		panic(err)
	}

	_ = db

	s := NewXServer()
	go s.Listen()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	if err := s.Shutdown(); err != nil {
		log.Fatalf("server shutdown err - %v", err)
	}
}

type XServer struct {
	serv *http.Server
}

// register, auth, token, refresh, forgot-pw, reset-pw, logout
func NewXServer() *XServer {
	s := http.NewServeMux()
	s.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is running.."))
	})

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("regis")) })
	apiMux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("auth")) })
	apiMux.HandleFunc("/refresh-token", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("rt")) })
	// apiMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("api/v1 uri - %s %s", r.RequestURI, r.URL)
	// 	w.Write([]byte("this is api/v1 root uri"))
	// })

	// work around how to chain middleware
	wrappedRoutes := []Route{
		{pattern: "/protected", handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("protected"))
		}},
		{pattern: "/dashboard", handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("dashboard"))
		}},
	}
	RegisterRoutes(apiMux, NewMiddlewares(firstMidleware, logginMiddleware), wrappedRoutes)

	s.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))

	serv := &http.Server{Addr: ":3000", Handler: s}
	return &XServer{serv}
}

// ######
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

// #######

// ##### grouping routes with wrapped middlewares
type Route struct {
	pattern string
	handler http.HandlerFunc
}

func RegisterRoutes(s *http.ServeMux, composer *MiddlewareChain, routes []Route) {
	for _, r := range routes {
		s.Handle(r.pattern, composer.Handle(r.handler))
	}
}

// #####

func logginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("before:loggingMiddleware")
		next.ServeHTTP(w, r)
		log.Println("after:loggingMiddleware")
	})
}

func firstMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("before:first middleware")
		next.ServeHTTP(w, r)
		log.Println("after:first middleware")
	})
}

func test() http.Handler {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { log.Println("ok") })
	mws := NewMiddlewares(firstMidleware, logginMiddleware)
	return mws.Handle(h)
}

func (s *XServer) Listen() {
	log.Printf("server is listening on - %s", s.serv.Addr)
	if err := s.serv.ListenAndServe(); err != nil {
		log.Print("server is closed")
	}
}

func (s *XServer) Shutdown() error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.serv.Shutdown(ctx)
}
