package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/github.com/steevehook/account-api/middleware"
)

// AuthService represents the authentication service
type AuthService interface {
	loginner
	logoutter
	signupper
}

// RouterConfig represents the application router config
type RouterConfig struct {
	AuthSvc AuthService
}

// NewRouter creates a new application HTTP router
func NewRouter(cfg RouterConfig) http.Handler {
	chain := alice.New(
		middleware.HTTPLogger,
	)
	jsonBodyChain := chain.Append(
		middleware.JSONBody,
	)
	route := func(h http.Handler) http.Handler {
		return chain.Then(h)
	}
	routeWithBody := func(h http.Handler) http.Handler {
		return jsonBodyChain.Then(h)
	}

	router := httprouter.New()
	router.Handler(http.MethodPost, "/login", routeWithBody(login(cfg.AuthSvc)))
	router.Handler(http.MethodPost, "/signup", routeWithBody(signup(cfg.AuthSvc)))
	router.Handler(http.MethodPost, "/logout", routeWithBody(logout(cfg.AuthSvc)))
	router.NotFound = route(NotFound())

	return router
}
