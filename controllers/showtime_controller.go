package controllers

import (
	"net/http"

	"github.com/sipatir08/cinema-ticket-api/service"
)

type ShowtimeController struct {
	showtime *service.ShowtimeService
}

func NewShowtimeController(showtime *service.ShowtimeService) *ShowtimeController {
	return &ShowtimeController{showtime: showtime}
}

func (h *ShowtimeController) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.showtime.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, list)
}

func (h *ShowtimeController) Get(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	s, err := h.showtime.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, s)
}

func (h *ShowtimeController) Create(w http.ResponseWriter, r *http.Request) {
	var in service.CreateShowtimeInput
	if err := decodeJSON(w, r, &in); err != nil {
		writeError(w, err)
		return
	}
	s, err := h.showtime.Create(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusCreated, s)
}

func (h *ShowtimeController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	var in service.UpdateShowtimeInput
	if err := decodeJSON(w, r, &in); err != nil {
		writeError(w, err)
		return
	}
	s, err := h.showtime.Update(r.Context(), id, in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, s)
}

func (h *ShowtimeController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.showtime.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}