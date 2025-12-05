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
	resp, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	return c.JSON(resp)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	resp, err := h.svc.Register(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
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
	var req struct {
		SessionID string `json:"session_id"`
	}
	_ = c.BodyParser(&req)
	resp := h.svc.Logout(c.Context(), req.SessionID)
	return c.JSON(resp)
}

func (h *AuthHandler) ValidateSession(c *fiber.Ctx) error {
	var req dto.ValidateSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	resp, err := h.svc.ValidateSession(c.Context(), req)
	if err != nil {
		return c.JSON(dto.ValidateSessionResponse{Valid: false})
	}
	return c.JSON(resp)
}
