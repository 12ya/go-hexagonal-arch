package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/12ya/go-hexagonal-arch/internal/common/errors"
)

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

func InternalError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Internal server error", http.StatusInternalServerError)
}

func Unauthorized(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Unauthorized", http.StatusUnauthorized)
}

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Bad Request", http.StatusBadRequest)
}

func NotFound(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Not Found", http.StatusBadRequest)
}

func RespondWithError(err error, w http.ResponseWriter, r *http.Request) {
	slugError, ok := err.(errors.SlugError)
	if !ok {
		InternalError("intrnal-server-error", err, w, r)
		return
	}

	switch slugError.ErrorType() {
	case errors.ErrorTypeAuthorization:
		Unauthorized(slugError.Slug(), slugError, w, r)
	default:
		InternalError(slugError.Slug(), slugError, w, r)
	}
}
