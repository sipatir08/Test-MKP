package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sipatir08/cinema-ticket-api/models"
)

type UserRepo struct {
	db *pgxpool.Pool
}

// Koneksi DB dioper dari luar (main.go), bukan ambil variabel global.
func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrNotFound // service ga perlu tau apa itu pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}