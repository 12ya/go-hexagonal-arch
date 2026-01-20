package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

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
			httpRespondWithError(err, "port-not-found", w, r, "not found", http.StatusNotFound)
			return
		}
		httpRespondWithError(err, "", w, r, "internal server error", http.StatusInternalServerError)
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

	httpRespond(response, w, r)
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
			httpRespond(map[string]int{"total_ports": portCounter}, w, r)
			return
		case err := <-errChan:
			log.Printf("Error occured while parsing port json: %+v", err)
			httpRespondWithError(err, "invalid-json", w, r, "Invalid json", http.StatusBadRequest)
			return
		case port := <-portChan:
			portCounter++
			log.Printf("[%d] received port: %+v", portCounter, port)
			p, err := httpPortToDomain(&port)
			if err != nil {
				httpRespondWithError(err, "port-to-domain", w, r, "Error parsing Port Http -> Domain", http.StatusBadRequest)
				return
			}
			if err := h.service.Upsert(r.Context(), p); err != nil {
				httpRespondWithError(err, "upsert-error", w, r, "Error upserting in UploadPorts", http.StatusBadRequest)
				return
			}
		}
	}
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, message string, status int) {
	log.Printf("error: %s, slug: %s, message: %s, status: %d, path: %s, method: %s",
		err, slug, message, status, r.URL.Path, r.Method)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{slug, status})
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	httpStatus int
}

func httpRespond(data any, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// reads Ports from provided Reader and sends them to portChan
func readPorts(ctx context.Context, r io.Reader, portChan chan Port) error {
	decoder := json.NewDecoder(r)

	// Read the opening delimiter
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("Failed to read the opening delimiter: %w", err)
	}

	if token != json.Delim('{') {
		return fmt.Errorf("Expected {, got %v", token)
	}

	for decoder.More() {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		token, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("Expected string, got %v", token)
		}

		portID, ok := token.(string)
		if !ok {
			return fmt.Errorf("Expected string, got %v", token)
		}

		var port Port
		if err := decoder.Decode(&port); err != nil {
			return fmt.Errorf("Failed to decode port: %w", err)
		}

		port.ID = portID
		portChan <- port
	}

	return nil
}

func httpPortToDomain(p *Port) (*domain.Port, error) {
	return domain.NewPort(
		p.ID,
		p.Name,
		p.Code,
		p.City,
		p.Country,
		append([]string(nil), p.Alias...),
		append([]string(nil), p.Regions...),
		append([]float64(nil), p.Coordinates...),
		p.Province,
		p.Timezone,
		append([]string(nil), p.Unlocs...),
	)
}
