package models

import "time"

type User struct {
	ID        int
	MerchantID *int
	Email     string
	PasswordHash string
	Role      string
	IsActive  bool
	LastLogin *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
