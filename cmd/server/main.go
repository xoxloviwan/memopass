package main

import (
	iHttp "iwakho/gopherkeep/internal/srv/http"
	"iwakho/gopherkeep/internal/srv/http/handlers"
	iLog "iwakho/gopherkeep/internal/srv/log"
	"net/http"
	"os"
)

func main() {

	addr := ":4443"
	logger := iLog.New()

	router := iHttp.NewRouter(http.NewServeMux(), logger)
	mux := router.SetupRoutes(&handlers.Handler{})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	logger.Info("Starting server", "addr", srv.Addr)
	err := srv.ListenAndServeTLS("./certs/localhost+1.pem", "./certs/localhost+1-key.pem")
	if err != nil {
		logger.Error("Server error", "error", err)
		os.Exit(1)
	}
}
