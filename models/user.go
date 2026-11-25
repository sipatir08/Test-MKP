package models

import (
	"context"
	"time"
)

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" = ga pernah ikut ke response JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuthClaims = identitas user yang dibawa di dalam token JWT
type AuthClaims struct {
	UserID int64
	Role   string
}

// Kontrak: service butuh fungsi ini dari database.
// Implementasinya nanti di repository/, service ga peduli caranya.
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
}