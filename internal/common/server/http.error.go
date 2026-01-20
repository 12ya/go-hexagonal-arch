package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func HttpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, message string, status int) {
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
