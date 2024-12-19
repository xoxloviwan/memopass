package http

import (
	"iwakho/gopherkeep/internal/srv/http/handlers"
	"iwakho/gopherkeep/internal/srv/http/middleware"
	"iwakho/gopherkeep/internal/srv/log"
	"net/http"

	"github.com/go-pkgz/routegroup"
)

type Router struct {
	mux    *http.ServeMux
	logger log.Log
}

// NewRouter creates a new router
func NewRouter(mux *http.ServeMux, logger log.Log) *Router {
	return &Router{
		mux,
		logger,
	}
}

func (rr *Router) SetupRoutes(h *handlers.Handler) http.Handler {
	router := routegroup.New(rr.mux)
	router.Use(middleware.Logging(rr.logger))

	// create a new group for the /api/user path
	apiRouter := router.Mount("/api/user")
	apiRouter.HandleFunc("POST /signup", h.SignUp)
	apiRouter.HandleFunc("POST /login", h.Auth)

	protectedGroup := router.Mount("/api/v1")
	protectedGroup.Use(middleware.CheckAuth)
	protectedGroup.HandleFunc("GET /protected", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}
