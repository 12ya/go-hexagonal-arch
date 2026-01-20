package transport

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/12ya/go-hexagonal-arch/internal/common/server"
	"github.com/12ya/go-hexagonal-arch/internal/domain"
)

type PortService interface {
	Upsert(context.Context, *domain.Port) error
	GetPort(ctx context.Context, id string) (*domain.Port, error)
	CountPorts(context.Context) (int, error)
}

type HttpServer struct {
	service PortService
}

func NewHttpServer(service PortService) HttpServer {
	return HttpServer{service: service}
}

func (h *HttpServer) GetPort(w http.ResponseWriter, r *http.Request) {
	port, err := h.service.GetPort(r.Context(), r.URL.Query().Get("id"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.HttpRespondWithError(err, "port-not-found", w, r, "not found", http.StatusNotFound)
			return
		}
		server.HttpRespondWithError(err, "", w, r, "internal server error", http.StatusInternalServerError)
		return
	}

	response := Port{
		ID:          port.ID(),
		Name:        port.Name(),
		City:        port.City(),
		Country:     port.Country(),
		Alias:       port.Alias(),
		Regions:     port.Regions(),
		Coordinates: port.Coordinates(),
		Province:    port.Province(),
		Timezone:    port.Timezone(),
		Unlocs:      port.Unlocs(),
	}

	server.RespondOK(response, w, r)
}

// reads ports from JSON and creates/updates them in DB
func (h *HttpServer) UploadPorts(w http.ResponseWriter, r *http.Request) {
	log.Println("Uploading ports")

	portChan := make(chan Port)
	errChan := make(chan error)
	doneChan := make(chan struct{})

	go func() {
		if err := readPorts(r.Context(), r.Body, portChan); err != nil {
			errChan <- err
		} else {
			doneChan <- struct{}{}
		}
	}()

	portCounter := 0

	for {
		select {
		case <-r.Context().Done():
			log.Printf("Request cancelled")
			return
		case <-doneChan:
			log.Printf("Finished reading ports")
			server.RespondOK(map[string]int{"total_ports": portCounter}, w, r)
			return
		case err := <-errChan:
			log.Printf("Error occured while parsing port json: %+v", err)
			server.HttpRespondWithError(err, "invalid-json", w, r, "Invalid json", http.StatusBadRequest)
			return
		case port := <-portChan:
			portCounter++
			log.Printf("[%d] received port: %+v", portCounter, port)
			p, err := httpPortToDomain(&port)
			if err != nil {
				server.HttpRespondWithError(err, "port-to-domain", w, r, "Error parsing Port Http -> Domain", http.StatusBadRequest)
				return
			}
			if err := h.service.Upsert(r.Context(), p); err != nil {
				server.HttpRespondWithError(err, "upsert-error", w, r, "Error upserting in UploadPorts", http.StatusBadRequest)
				return
			}
		}
	}
}
