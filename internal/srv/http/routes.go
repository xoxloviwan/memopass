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

	apiRouter := router.Mount("/api/v1")
	userLayer := apiRouter.Mount("/user")
	userLayer.HandleFunc("POST /signup", h.SignUp)
	userLayer.HandleFunc("GET /login", h.Login)

	apiRouter.Use(middleware.CheckAuth)

	formGroup := apiRouter.Mount("/item/add")
	formGroup.Use(middleware.ParseForm(rr.logger))
	formGroup.HandleFunc("POST /pair", h.AddPair)
	formGroup.HandleFunc("POST /card", h.AddCard)
	formGroup.HandleFunc("POST /file", h.AddBinary)
	formGroup.HandleFunc("POST /text", h.AddText)

	getGroup := apiRouter.Mount("/item")
	getGroup.Use(middleware.ParseQueryParams(rr.logger))
	getGroup.HandleFunc("GET /pairs", h.GetPairs)
	getGroup.HandleFunc("GET /cards", h.GetCards)
	getGroup.HandleFunc("GET /files", h.GetBinaries)
	getGroup.HandleFunc("GET /texts", h.GetTexts)

	return router
}
