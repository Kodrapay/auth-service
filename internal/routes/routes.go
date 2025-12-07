package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/auth-service/internal/handlers"
	"github.com/kodra-pay/auth-service/internal/config"
	"github.com/kodra-pay/auth-service/internal/repositories"
	"github.com/kodra-pay/auth-service/internal/services"
	"github.com/kodra-pay/auth-service/internal/session"
)

func Register(app *fiber.App, serviceName string) {
	health := handlers.NewHealthHandler(serviceName)
	health.Register(app)

	cfg := config.Load(serviceName, "7001")
	repo, err := repositories.NewAuthRepository(cfg.PostgresDSN)
	if err != nil {
		panic(err)
	}

	sessionMgr, err := session.NewRedisSessionManager(cfg.RedisAddr)
	if err != nil {
		log.Printf("Warning: Failed to initialize Redis session manager: %v", err)
		sessionMgr = nil
	}

	authSvc := services.NewAuthService(repo, cfg, sessionMgr)
	authHandler := handlers.NewAuthHandler(cfg, authSvc)

	app.Post("/auth/login", authHandler.Login)
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/refresh", authHandler.Refresh)
	app.Post("/auth/logout", authHandler.Logout)
	app.Post("/auth/validate", authHandler.ValidateSession)
}
