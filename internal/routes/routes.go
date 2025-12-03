package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/auth-service/internal/config"
	"github.com/kodra-pay/auth-service/internal/handlers"
	"github.com/kodra-pay/auth-service/internal/repositories"
	"github.com/kodra-pay/auth-service/internal/services"
)

func Register(app *fiber.App, cfg config.Config, repo *repositories.AuthRepository) {
	health := handlers.NewHealthHandler(cfg.ServiceName)
	health.Register(app)

	svc := services.NewAuthService(repo, cfg)
	h := handlers.NewAuthHandler(cfg, svc)
	api := app.Group("/auth")
	api.Post("/login", h.Login)
	api.Post("/register", h.Register)
	api.Post("/refresh", h.Refresh)
	api.Post("/logout", h.Logout)
}
