package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodra-pay/auth-service/internal/config"
	"github.com/kodra-pay/auth-service/internal/dto"
	"github.com/kodra-pay/auth-service/internal/services"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(cfg config.Config, svc *services.AuthService) *AuthHandler {
	_ = cfg
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	resp := h.svc.Login(c.Context(), req)
	return c.JSON(resp)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	resp := h.svc.Register(c.Context(), req)
	return c.JSON(resp)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	resp := h.svc.Refresh(c.Context(), req)
	return c.JSON(resp)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	resp := h.svc.Logout(c.Context())
	return c.JSON(resp)
}
