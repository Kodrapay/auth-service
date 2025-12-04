package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/kodra-pay/auth-service/internal/models"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(dsn string) (*AuthRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return &AuthRepository{db: db}, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, merchant_id, email, password_hash, role, is_active, last_login, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var u models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.MerchantID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.IsActive,
		&u.LastLogin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, email, passwordHash, role string, merchantID *string) (string, error) {
	query := `
		INSERT INTO users (id, merchant_id, email, password_hash, role, is_active)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4, TRUE)
		RETURNING id
	`
	var id string
	if err := r.db.QueryRowContext(ctx, query, merchantID, email, passwordHash, role).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *AuthRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET last_login = NOW(), updated_at = NOW() WHERE id = $1`, userID)
	return err
}
