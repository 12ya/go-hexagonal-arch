package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/12ya/go-hexagonal-arch/internal/config"
	"github.com/12ya/go-hexagonal-arch/internal/repository/inmemory"
	"github.com/12ya/go-hexagonal-arch/internal/services"
	"github.com/12ya/go-hexagonal-arch/internal/transport"
	"github.com/gorilla/mux"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.Read()

	portStoreRepo := inmemory.NewPortStore()

	portService := services.NewPortService(portStoreRepo)

	httpServer := transport.NewHttpServer(portService)

	router := mux.NewRouter()
	router.HandleFunc("/port", httpServer.GetPort).Methods("GET")
	router.HandleFunc("/count", httpServer.CountPorts).Methods("GET")
	router.HandleFunc("/ports", httpServer.GetPort).Methods("POST")

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("http server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("http server ListenAndServe Error: %v", err)
	}

	<-stopped
	log.Printf("Bye!")

	return nil
}
