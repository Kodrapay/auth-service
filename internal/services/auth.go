package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kodra-pay/auth-service/internal/config"
	"github.com/kodra-pay/auth-service/internal/dto"
	"github.com/kodra-pay/auth-service/internal/repositories"
	"github.com/kodra-pay/auth-service/internal/session"
)

type AuthService struct {
	repo       *repositories.AuthRepository
	cfg        config.Config
	sessionMgr SessionManager
}

type SessionManager interface {
	CreateSession(ctx context.Context, sessionID string, data interface{}) error
	GetSession(ctx context.Context, sessionID string) (interface{}, error)
	DeleteSession(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, sessionID string) error
}

func NewAuthService(repo *repositories.AuthRepository, cfg config.Config, sessionMgr SessionManager) *AuthService {
	return &AuthService{repo: repo, cfg: cfg, sessionMgr: sessionMgr}
}

func (s *AuthService) Login(_ context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	ctx := context.Background()
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if user == nil {
		return dto.LoginResponse{}, fmt.Errorf("invalid credentials")
	}
	if !user.IsActive {
		return dto.LoginResponse{}, fmt.Errorf("account inactive")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return dto.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	accessToken, refreshToken, exp := s.generateTokens(user.ID, user.Role)
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	// Create Redis session if session manager is available
	sessionID := ""
	if s.sessionMgr != nil {
		sessionID = generateSessionID()
		merchantID := ""
		if user.MerchantID != nil {
			merchantID = *user.MerchantID
		}

		sessionData := session.SessionData{
			UserID:     user.ID,
			Role:       user.Role,
			MerchantID: merchantID,
			Email:      req.Email,
		}
		_ = s.sessionMgr.CreateSession(ctx, sessionID, sessionData)
	}

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    exp,
		MerchantID:   user.MerchantID,
		Role:         user.Role,
		SessionID:    sessionID,
		Email:        req.Email,
	}, nil
}

func (s *AuthService) Register(_ context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	ctx := context.Background()
	if req.Email == "" || req.Password == "" {
		return dto.RegisterResponse{}, fmt.Errorf("email and password are required")
	}

	// Prevent duplicate emails
	if existing, _ := s.repo.GetUserByEmail(ctx, req.Email); existing != nil {
		return dto.RegisterResponse{}, fmt.Errorf("email already exists")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	userID, err := s.repo.CreateUser(ctx, req.Email, string(hashed), "merchant", req.MerchantID)
	if err != nil {
		return dto.RegisterResponse{}, err
	}
	accessToken, refreshToken, _ := s.generateTokens(userID, "merchant")
	return dto.RegisterResponse{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		MerchantID:   deref(req.MerchantID),
	}, nil
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
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

func (s *AuthService) Logout(ctx context.Context, sessionID string) map[string]string {
	if s.sessionMgr != nil && sessionID != "" {
		_ = s.sessionMgr.DeleteSession(ctx, sessionID)
	}
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

func (s *AuthService) ValidateSession(ctx context.Context, req dto.ValidateSessionRequest) (dto.ValidateSessionResponse, error) {
	if s.sessionMgr == nil {
		return dto.ValidateSessionResponse{Valid: false}, fmt.Errorf("session manager not available")
	}

	sessionData, err := s.sessionMgr.GetSession(ctx, req.SessionID)
	if err != nil {
		return dto.ValidateSessionResponse{Valid: false}, nil
	}

	// Type assert the session data
	data, ok := sessionData.(session.SessionData)
	if !ok {
		return dto.ValidateSessionResponse{Valid: false}, nil
	}

	return dto.ValidateSessionResponse{
		Valid:      true,
		UserID:     data.UserID,
		Role:       data.Role,
		MerchantID: data.MerchantID,
		Email:      data.Email,
	}, nil
}

func generateSessionID() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("session_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
