package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	iHttp "iwakho/gopherkeep/internal/srv/http"
	"iwakho/gopherkeep/internal/srv/http/handlers"
	iLog "iwakho/gopherkeep/internal/srv/log"
	"iwakho/gopherkeep/internal/srv/store"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	logger := iLog.New(buildVersion, false)

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

	// Через этот канал сообщим основному потоку, что соединения закрыты.
	idleConnsClosed := make(chan struct{})
	// Создаем канал для сигналов завершения.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-quit
		signal.Stop(quit)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		logger.Info("Shutdown Server...")
		if err := srv.Shutdown(ctx); err != nil {
			// ошибки закрытия Listener
			logger.Error("Server shutdown error", "error", err)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	logger.Info("Starting server", "addr", srv.Addr)
	err = srv.ListenAndServeTLS(defaultCertPath, defaultKeyPath)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server error", "error", err)
		os.Exit(1)
	}
	<-idleConnsClosed
	logger.Info("Server stopped")
}
