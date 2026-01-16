package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/12ya/go-hexagonal-arch/internal/domain"
)

type PortService interface {
	Upsert(context.Context, *domain.Port) error
	GetPort(ctx context.Context, id string) (*domain.Port, error)
}

type HttpServer struct {
	service PortService
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

	fmt.Println(port)
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
