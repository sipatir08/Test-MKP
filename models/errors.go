package models

import "errors"

// Error bisnis. Controller yang nerjemahin ini ke HTTP status code,
// jadi service & repository ga perlu tau soal HTTP sama sekali.
var (
	ErrValidation         = errors.New("validation error")
	ErrNotFound           = errors.New("resource not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidReference   = errors.New("movie or studio does not exist")
	ErrScheduleConflict   = errors.New("studio already has a showtime in that time range")
	ErrShowtimeHasBooking = errors.New("showtime already has bookings, set status to cancelled instead")
)