package models

import (
	"context"
	"time"
)

const (
	ShowtimeScheduled = "scheduled"
	ShowtimeCancelled = "cancelled"
)

type Showtime struct {
	ID        int64     `json:"id"`
	MovieID   int64     `json:"movie_id"`
	StudioID  int64     `json:"studio_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"` // dihitung database dari start_time + durasi film
	Price     int64     `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShowtimeRepository interface {
	// Create menyimpan jadwal sekaligus membuat stok kursi (showtime_seats).
	Create(ctx context.Context, s *Showtime) error
	GetByID(ctx context.Context, id int64) (*Showtime, error)
	List(ctx context.Context) ([]Showtime, error)
	Update(ctx context.Context, s *Showtime) error
	Delete(ctx context.Context, id int64) error
}