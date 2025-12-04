package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

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
	ctx := context.Background()
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return dto.LoginResponse{}
	}
	if !user.IsActive {
		return dto.LoginResponse{}
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return dto.LoginResponse{}
	}

	accessToken, refreshToken, exp := s.generateTokens(user.ID, user.Role)
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    exp,
		MerchantID:   user.MerchantID,
		Role:         user.Role,
	}
}

func (s *AuthService) Register(_ context.Context, req dto.RegisterRequest) dto.RegisterResponse {
	ctx := context.Background()
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	userID, err := s.repo.CreateUser(ctx, req.Email, string(hashed), "merchant", nil)
	if err != nil {
		return dto.RegisterResponse{}
	}
	accessToken, refreshToken, _ := s.generateTokens(userID, "merchant")
	return dto.RegisterResponse{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (s *AuthService) Refresh(_ context.Context, req dto.RefreshRequest) dto.RefreshResponse {
	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return dto.RefreshResponse{}
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return dto.RefreshResponse{}
	}
	sub, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	accessToken, refreshToken, exp := s.generateTokens(sub, role)
	return dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    exp,
	}
}

func (s *AuthService) Logout(_ context.Context) map[string]string {
	return map[string]string{"status": "logged_out"}
}

func (s *AuthService) generateTokens(userID, role string) (string, string, int64) {
	exp := time.Now().Add(time.Hour).Unix()
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  exp,
	})
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"type": "refresh",
	})
	accessToken, _ := access.SignedString([]byte(s.cfg.JWTSecret))
	refreshToken, _ := refresh.SignedString([]byte(s.cfg.JWTSecret))
	return accessToken, refreshToken, time.Now().Add(time.Hour).Unix()
}
