package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sipatir08/cinema-ticket-api/models"
)

type CreateShowtimeInput struct {
	MovieID   int64     `json:"movie_id"`
	StudioID  int64     `json:"studio_id"`
	StartTime time.Time `json:"start_time"`
	Price     int64     `json:"price"`
}

// studio_id sengaja tidak bisa diubah: stok kursi jadwal terikat ke denah studio.
type UpdateShowtimeInput struct {
	MovieID   int64     `json:"movie_id"`
	StartTime time.Time `json:"start_time"`
	Price     int64     `json:"price"`
	Status    string    `json:"status"`
}

type ShowtimeService struct {
	repo models.ShowtimeRepository
}

func NewShowtimeService(repo models.ShowtimeRepository) *ShowtimeService {
	return &ShowtimeService{repo: repo}
}

func (s *ShowtimeService) Create(ctx context.Context, in CreateShowtimeInput) (*models.Showtime, error) {
	if in.MovieID <= 0 || in.StudioID <= 0 {
		return nil, fmt.Errorf("%w: movie_id and studio_id are required", models.ErrValidation)
	}
	if err := validateTimeAndPrice(in.StartTime, in.Price); err != nil {
		return nil, err
	}

	st := &models.Showtime{
		MovieID:   in.MovieID,
		StudioID:  in.StudioID,
		StartTime: in.StartTime,
		Price:     in.Price,
		Status:    models.ShowtimeScheduled,
	}
	if err := s.repo.Create(ctx, st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *ShowtimeService) Get(ctx context.Context, id int64) (*models.Showtime, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ShowtimeService) List(ctx context.Context) ([]models.Showtime, error) {
	return s.repo.List(ctx)
}

func (s *ShowtimeService) Update(ctx context.Context, id int64, in UpdateShowtimeInput) (*models.Showtime, error) {
	if in.MovieID <= 0 {
		return nil, fmt.Errorf("%w: movie_id is required", models.ErrValidation)
	}
	if in.Status != models.ShowtimeScheduled && in.Status != models.ShowtimeCancelled {
		return nil, fmt.Errorf("%w: status must be scheduled or cancelled", models.ErrValidation)
	}
	if err := validateTimeAndPrice(in.StartTime, in.Price); err != nil {
		return nil, err
	}

	st := &models.Showtime{
		ID:        id,
		MovieID:   in.MovieID,
		StartTime: in.StartTime,
		Price:     in.Price,
		Status:    in.Status,
	}
	if err := s.repo.Update(ctx, st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *ShowtimeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func validateTimeAndPrice(start time.Time, price int64) error {
	if start.IsZero() {
		return fmt.Errorf("%w: start_time is required (RFC3339, e.g. 2026-12-01T19:00:00+07:00)", models.ErrValidation)
	}
	if !start.After(time.Now()) {
		return fmt.Errorf("%w: start_time must be in the future", models.ErrValidation)
	}
	if price <= 0 {
		return fmt.Errorf("%w: price must be greater than 0", models.ErrValidation)
	}
	return nil
}