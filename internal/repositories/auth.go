package repositories

import "log"

type AuthRepository struct {
    dsn string
}

func NewAuthRepository(dsn string) *AuthRepository {
    // placeholder for actual DB connection (gorm/sqlx/sqlc)
    log.Printf("AuthRepository using DSN: %s", dsn)
    return &AuthRepository{dsn: dsn}
}

// TODO: add methods for user CRUD, credentials, sessions.
