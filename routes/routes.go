package routes

import (
	"net/http"

	"github.com/sipatir08/cinema-ticket-api/controllers"
	"github.com/sipatir08/cinema-ticket-api/middleware"
	"github.com/sipatir08/cinema-ticket-api/models"
)

func SetupRoutes(auth *controllers.AuthController, showtime *controllers.ShowtimeController, v middleware.TokenVerifier) http.Handler {
	authed := func(h http.HandlerFunc) http.HandlerFunc { return middleware.Auth(v, h) }
	admin := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.Auth(v, middleware.RequireRole(models.RoleAdmin, h))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/login", auth.Login)

	mux.HandleFunc("GET /api/v1/showtimes", authed(showtime.List))
	mux.HandleFunc("GET /api/v1/showtimes/{id}", authed(showtime.Get))
	mux.HandleFunc("POST /api/v1/showtimes", admin(showtime.Create))
	mux.HandleFunc("PUT /api/v1/showtimes/{id}", admin(showtime.Update))
	mux.HandleFunc("DELETE /api/v1/showtimes/{id}", admin(showtime.Delete))
	return mux
}