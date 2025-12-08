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
	var merchantID sql.NullInt32 // To handle nullable int for MerchantID

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, // This should be int based on model change
		&merchantID, // Scan into sql.NullInt32
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

	// Handle nullable merchantID
	if merchantID.Valid {
		val := int(merchantID.Int32)
		u.MerchantID = &val
	} else {
		u.MerchantID = nil
	}

	return &u, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, email, passwordHash, role string, merchantID *int) (int, error) { // merchantID *int, return int
	query := `
		INSERT INTO users (merchant_id, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, TRUE)
		RETURNING id
	`
	var id int // id is int
	if err := r.db.QueryRowContext(ctx, query, merchantID, email, passwordHash, role).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthRepository) UpdateLastLogin(ctx context.Context, userID int) error { // userID changed to int
	_, err := r.db.ExecContext(ctx, `UPDATE users SET last_login = NOW(), updated_at = NOW() WHERE id = $1`, userID)
	return err
}