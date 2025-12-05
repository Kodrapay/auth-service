package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID     string `json:"user_id"`
	Role       string `json:"role"`
	MerchantID string `json:"merchant_id,omitempty"`
	Email      string `json:"email"`
	CreatedAt  int64  `json:"created_at"`
	ExpiresAt  int64  `json:"expires_at"`
}

type RedisSessionManager struct {
	client *redis.Client
}

func NewRedisSessionManager(redisAddr string) (*RedisSessionManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisSessionManager{client: client}, nil
}

func (m *RedisSessionManager) CreateSession(ctx context.Context, sessionID string, data interface{}) error {
	sessionData, ok := data.(SessionData)
	if !ok {
		return fmt.Errorf("invalid session data type")
	}

	sessionData.CreatedAt = time.Now().Unix()
	sessionData.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	key := fmt.Sprintf("session:%s", sessionID)
	return m.client.Set(ctx, key, jsonData, 24*time.Hour).Err()
}

func (m *RedisSessionManager) GetSession(ctx context.Context, sessionID string) (interface{}, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	data, err := m.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session SessionData
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	// Check if session is expired
	if session.ExpiresAt < time.Now().Unix() {
		_ = m.DeleteSession(ctx, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return session, nil
}

func (m *RedisSessionManager) GetSessionData(ctx context.Context, sessionID string) (*SessionData, error) {
	data, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	session, ok := data.(SessionData)
	if !ok {
		return nil, fmt.Errorf("invalid session data type")
	}
	return &session, nil
}

func (m *RedisSessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return m.client.Del(ctx, key).Err()
}

func (m *RedisSessionManager) RefreshSession(ctx context.Context, sessionID string) error {
	session, err := m.GetSessionData(ctx, sessionID)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()
	return m.CreateSession(ctx, sessionID, *session)
}

func (m *RedisSessionManager) Close() error {
	return m.client.Close()
}
