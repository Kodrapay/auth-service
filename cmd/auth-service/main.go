package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/auth-service/internal/config"
	"github.com/kodra-pay/auth-service/internal/middleware"
	"github.com/kodra-pay/auth-service/internal/routes"
)

func main() {
	cfg := config.Load("auth-service", "7001")

	app := fiber.New()
	app.Use(middleware.RequestID())

	routes.Register(app, cfg.ServiceName)

	log.Printf("%s listening on :%s", cfg.ServiceName, cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
