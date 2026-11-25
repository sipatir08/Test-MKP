package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sipatir08/cinema-ticket-api/models"
)

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeData(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, map[string]any{"data": data})
}

// writeError: satu-satunya tempat error models diterjemahkan ke HTTP status.
func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, models.ErrValidation):
		status = http.StatusBadRequest
	case errors.Is(err, models.ErrInvalidCredentials), errors.Is(err, models.ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(err, models.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, models.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, models.ErrInvalidReference):
		status = http.StatusUnprocessableEntity
	case errors.Is(err, models.ErrScheduleConflict), errors.Is(err, models.ErrShowtimeHasBooking):
		status = http.StatusConflict
	}

	msg := err.Error()
	if status == http.StatusInternalServerError {
		log.Printf("internal error: %v", err)
		msg = "internal server error" // detail error database tidak dibocorkan ke client
	}
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)) // body maksimal 1MB
	dec.DisallowUnknownFields()                                     // field asing seperti "id" atau "end_time" ditolak
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("%w: invalid JSON body: %v", models.ErrValidation, err)
	}
	return nil
}

func pathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("%w: id must be a positive integer", models.ErrValidation)
	}
	return id, nil
}