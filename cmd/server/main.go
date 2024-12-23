package main

import (
	"flag"
	"fmt"
	iHttp "iwakho/gopherkeep/internal/srv/http"
	"iwakho/gopherkeep/internal/srv/http/handlers"
	iLog "iwakho/gopherkeep/internal/srv/log"
	"iwakho/gopherkeep/internal/srv/store"
	"net/http"
	"os"
)

const (
	defaultAddress  = ":443"
	defaultDBPath   = "db/memopass.db"
	defaultCertPath = "./certs/cert.pem"
	defaultKeyPath  = "./certs/key.pem"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	vers         = flag.Bool("v", false, "version")
	address      = flag.String("a", defaultAddress, "server address")
	dbPath       = flag.String("d", defaultDBPath, "database path")
)

func init() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		fmt.Printf("Too many arguments")
		os.Exit(1)
	}
	if vers != nil && *vers {
		fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
		os.Exit(0)
	}
}

func main() {
	if address == nil {
		address = new(string)
		*address = defaultAddress
	}
	if dbPath == nil {
		dbPath = new(string)
		*dbPath = defaultDBPath
	}
	logger := iLog.New()

	db, err := store.NewStorage(*dbPath)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	hdl := handlers.NewHandler(db, logger)

	router := iHttp.NewRouter(http.NewServeMux(), logger)
	mux := router.SetupRoutes(hdl)

	srv := &http.Server{
		Addr:    *address,
		Handler: mux,
	}
	logger.Info("Starting server", "addr", srv.Addr)
	err = srv.ListenAndServeTLS(defaultCertPath, defaultKeyPath)
	if err != nil {
		logger.Error("Server error", "error", err)
		os.Exit(1)
	}
}
