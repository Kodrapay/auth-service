package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/auth-service/internal/handlers"
	"github.com/kodra-pay/auth-service/internal/config"
	"github.com/kodra-pay/auth-service/internal/repositories"
	"github.com/kodra-pay/auth-service/internal/services"
)

func Register(app *fiber.App, serviceName string) {
	health := handlers.NewHealthHandler(serviceName)
	health.Register(app)

	cfg := config.Load(serviceName, "7001")
	repo, err := repositories.NewAuthRepository(cfg.PostgresDSN)
	if err != nil {
		panic(err)
	}
	authSvc := services.NewAuthService(repo, cfg)
	authHandler := handlers.NewAuthHandler(cfg, authSvc)

	app.Post("/auth/login", authHandler.Login)
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/refresh", authHandler.Refresh)
}
