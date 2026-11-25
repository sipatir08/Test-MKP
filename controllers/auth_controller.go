package controllers

import (
	"net/http"

	"github.com/sipatir08/cinema-ticket-api/service"
)

type AuthController struct {
	auth *service.AuthService
}

func NewAuthController(auth *service.AuthService) *AuthController {
	return &AuthController{auth: auth}
}

func (h *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var in service.LoginInput
	if err := decodeJSON(w, r, &in); err != nil {
		writeError(w, err)
		return
	}
	out, err := h.auth.Login(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	writeData(w, http.StatusOK, out)
}