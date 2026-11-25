package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/sipatir08/cinema-ticket-api/config"
	"github.com/sipatir08/cinema-ticket-api/controllers"
	"github.com/sipatir08/cinema-ticket-api/repository"
	"github.com/sipatir08/cinema-ticket-api/routes"
	"github.com/sipatir08/cinema-ticket-api/service"
	"github.com/sipatir08/cinema-ticket-api/utils"
)

// main cuma merakit: config -> repository -> service -> controller -> routes
func main() {
	cfg := config.Load()
	db := config.ConnectDB(cfg.DatabaseURL)
	defer db.Close()

	jwt := utils.NewJWT(cfg.JWTSecret, 24*time.Hour)

	authService := service.NewAuthService(repository.NewUserRepo(db), utils.Bcrypt{}, jwt)
	showtimeService := service.NewShowtimeService(repository.NewShowtimeRepo(db))

	app := fiber.New()
	routes.SetupRoutes(app,
		controllers.NewAuthController(authService),
		controllers.NewShowtimeController(showtimeService),
		jwt,
	)

	log.Fatal(app.Listen(":" + cfg.Port))
}