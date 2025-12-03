package handlers

import "github.com/gofiber/fiber/v2"

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler { return &AuthHandler{} }

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"action": "login", "status": "not_implemented"})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"action": "register", "status": "not_implemented"})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"action": "refresh", "status": "not_implemented"})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"action": "logout", "status": "not_implemented"})
}
