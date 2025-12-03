package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/kodra-pay/auth-service/internal/config"
	"github.com/kodra-pay/auth-service/internal/dto"
	"github.com/kodra-pay/auth-service/internal/repositories"
)

type AuthService struct {
	repo *repositories.AuthRepository
	cfg  config.Config
}

func NewAuthService(repo *repositories.AuthRepository, cfg config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

func (s *AuthService) Login(_ context.Context, req dto.LoginRequest) dto.LoginResponse {
	return dto.LoginResponse{
		AccessToken:  "stub_access_" + uuid.NewString(),
		RefreshToken: "stub_refresh_" + uuid.NewString(),
		ExpiresIn:    int64(time.Hour.Seconds()),
	}
}

func (s *AuthService) Register(_ context.Context, req dto.RegisterRequest) dto.RegisterResponse {
	return dto.RegisterResponse{
		UserID:       uuid.NewString(),
		AccessToken:  "stub_access_" + uuid.NewString(),
		RefreshToken: "stub_refresh_" + uuid.NewString(),
	}
}

func (s *AuthService) Refresh(_ context.Context, req dto.RefreshRequest) dto.RefreshResponse {
	return dto.RefreshResponse{
		AccessToken:  "stub_access_" + uuid.NewString(),
		RefreshToken: "stub_refresh_" + uuid.NewString(),
		ExpiresIn:    int64(time.Hour.Seconds()),
	}
}

func (s *AuthService) Logout(_ context.Context) map[string]string {
	return map[string]string{"status": "logged_out"}
}
