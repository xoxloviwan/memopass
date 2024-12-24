package http

import (
	"fmt"
	"iwakho/gopherkeep/internal/srv/http/handlers"
	"iwakho/gopherkeep/internal/srv/http/middleware"
	"iwakho/gopherkeep/internal/srv/log"
	"iwakho/gopherkeep/internal/model"
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

	apiRouter := router.Mount("/api/v1")
	userLayer := apiRouter.Mount("/user")
	userLayer.HandleFunc("POST /signup", h.SignUp)
	userLayer.HandleFunc("GET /login", h.Login)

	apiRouter.Use(middleware.CheckAuth)
	apiRouter.HandleFunc("GET /protected", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(model.UserIDCtxKey{}).(int)
		w.Write([]byte(fmt.Sprintln(userID)))
		w.WriteHeader(http.StatusOK)
	})
	apiRouter.HandleFunc("POST /item/add", h.AddItem)

	return router
}
