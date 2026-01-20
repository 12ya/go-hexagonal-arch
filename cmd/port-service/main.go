package main

import (
	"log"
	"net/http"

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

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("http server ListenAndServe Error: %v", err)
	}

	log.Printf("Bye!")

	return nil
}
