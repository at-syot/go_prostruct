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
	"github.com/simt/pkg/httpx"
	"github.com/simt/server/internal/config"
)

func main() {
	// test loading config
	log.Printf("app config: %v", config.AppConfigurations)

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

type AppServer struct {
	serv                *http.Server
	shutdownGracePeriod time.Duration
}

// register, auth, token, refresh, forgot-pw, reset-pw, logout
func NewXServer() *AppServer {
	s := http.NewServeMux()
	s.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is running.."))
	})

	apiMux := http.NewServeMux()
	registerAuthRoutes(apiMux)
	registerProtectedRoutes(apiMux)

	s.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))

	serv := &http.Server{Addr: ":3000", Handler: httpx.MakeDevMiddlewares().Handle(s)}

	return &AppServer{serv, 10 * time.Second}
}

func registerAuthRoutes(s *http.ServeMux) {
	s.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("regis")) })
	s.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("auth")) })
	s.HandleFunc("/refresh-token", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("rt")) })
}

func registerProtectedRoutes(s *http.ServeMux) {
	protectedHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("protected")) }
	dashboardHandler := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("dashboard")) }
	wrappedRoutes := []httpx.Route{
		{Pattern: "GET /protected", Handler: protectedHandler},
		{Pattern: "GET /dashboard", Handler: dashboardHandler},
	}
	httpx.RegisterRoutes(s, httpx.MakeDevMiddlewares(), wrappedRoutes)
}

func (s *AppServer) Listen() {
	log.Printf("server is listening on - %s", s.serv.Addr)
	if err := s.serv.ListenAndServe(); err != nil {
		log.Print("server is closed")
	}
}

func (s *AppServer) Shutdown() error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, s.shutdownGracePeriod)
	defer cancel()

	return s.serv.Shutdown(ctx)
}
