package main

import (
	iHttp "iwakho/gopherkeep/internal/srv/http"
	"iwakho/gopherkeep/internal/srv/http/handlers"
	iLog "iwakho/gopherkeep/internal/srv/log"
	"iwakho/gopherkeep/internal/srv/store"
	"net/http"
	"os"
)

func main() {

	addr := ":443"
	logger := iLog.New()

	db, err := store.NewStorage("db/memopass.db")
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	hdl := handlers.NewHandler(db, logger)

	router := iHttp.NewRouter(http.NewServeMux(), logger)
	mux := router.SetupRoutes(hdl)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	logger.Info("Starting server", "addr", srv.Addr)
	err = srv.ListenAndServeTLS("./certs/localhost+1.pem", "./certs/localhost+1-key.pem")
	if err != nil {
		logger.Error("Server error", "error", err)
		os.Exit(1)
	}
}
