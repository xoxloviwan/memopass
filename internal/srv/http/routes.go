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
	mware := middleware.New(rr.logger)
	router := routegroup.New(rr.mux)
	router.Use(mware.Logging)

	apiRouter := router.Mount("/api/v1")
	userLayer := apiRouter.Mount("/user")
	userLayer.HandleFunc("POST /signup", h.SignUp)
	userLayer.HandleFunc("GET /login", h.Login)

	apiRouter.Use(mware.CheckAuth)

	formGroup := apiRouter.Mount("/item/add")
	formGroup.Use(mware.ParseForm)
	formGroup.HandleFunc("POST /pair", h.AddPair)
	formGroup.HandleFunc("POST /card", h.AddCard)
	formGroup.HandleFunc("POST /file", h.AddBinary)
	formGroup.HandleFunc("POST /text", h.AddText)

	getGroup := apiRouter.Mount("/item")
	formGroup.Use(mware.ParseQueryParams)
	getGroup.HandleFunc("GET /pair", h.GetPairs)
	getGroup.HandleFunc("GET /card", h.GetCards)
	getGroup.HandleFunc("GET /file", h.GetBinaries)
	getGroup.HandleFunc("GET /text", h.GetTexts)

	return router
}
