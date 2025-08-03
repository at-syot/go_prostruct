package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/simt/dtacc"
	"github.com/simt/pkg/httpx"
	"github.com/simt/pkg/logger"
	"github.com/simt/server/internal/config"
)

func main() {
	config.LoadConfig()
	appEnv := config.AppConfigurations.Env
	var loggerEnv logger.LogEnv
	switch appEnv {
	case config.ConfEnvDev:
		loggerEnv = logger.LogDev
	case config.ConfEnvStaging:
		loggerEnv = logger.LogStaging
	case config.ConfEnvProd:
		loggerEnv = logger.LogProd
	}
	logg := logger.InitLogger(loggerEnv)

	db, err := dtacc.NewDB()
	if err != nil {
		logg.Fatal().Err(err).Msg("failed to initialize database.")
	}

	_ = db

	s := NewAppServer(logg)
	go s.Listen()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	if err := s.Shutdown(); err != nil {
		logg.Fatal().Err(err).Msg("server shutdown error")
	}
}

type AppServer struct {
	serv                *http.Server
	shutdownGracePeriod time.Duration
	logger              zerolog.Logger
}

func NewAppServer(logg zerolog.Logger) *AppServer {
	s := http.NewServeMux()
	s.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server is running.."))
	})

	apiMux := http.NewServeMux()
	registerAuthRoutes(apiMux)
	registerProtectedRoutes(apiMux)
	s.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))

	serv := &http.Server{Addr: ":" + config.AppConfigurations.Port, Handler: httpx.MakeDevMiddlewares().Handle(s)}

	return &AppServer{serv, config.AppConfigurations.ShutdownGracePeriod, logg}
}

// register, auth, token, refresh, forgot-pw, reset-pw, logout
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
	s.logger.Info().Str("addr", s.serv.Addr).Msg("server is listening")
	if err := s.serv.ListenAndServe(); err != nil {
		s.logger.Info().Msg("server is closed")
	}
}

func (s *AppServer) Shutdown() error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, s.shutdownGracePeriod)
	defer cancel()

	return s.serv.Shutdown(ctx)
}
