package auth

import (
	"net/http"

	"github.com/uptrace/bun"
)

// Features
// - register || login by [outlook, gmail, email] (impl login by email first)

// Auth
//   +handlers
//    service
//    repo

// Auth have their own data schema and migrations
// Expose all handlers

// Init Auth struct with *sql.DB && *bun.DB connections
//   do some migrations if not have

type Auth struct {
	s  *http.ServeMux
	db *bun.DB
}

func NewAuth(s *http.ServeMux, db *bun.DB) *Auth {
	return &Auth{s, db}
}

// register, auth, token, refresh, forgot-pw, reset-pw, logout
func (a *Auth) RegisterHandlers() {
	a.s.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("regis")) })
	a.s.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("auth")) })
	a.s.HandleFunc("/token/refresh", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("rt")) })

	// a.s.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("register endpoint"))
	// })
	// a.s.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("auth endpoint"))
	// })
	// a.s.HandleFunc("/token/refresh", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("token refresh endpoint"))
	// })
	// a.s.HandleFunc("/password/forgot", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("forgot password endpoint"))
	// })
	// a.s.HandleFunc("/password/reset", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("reset password endpoint"))
	// })
	// a.s.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("logout endpoint"))
	// })
}
