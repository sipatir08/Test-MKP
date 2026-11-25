package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sipatir08/cinema-ticket-api/models"
)

// Kode error PostgreSQL yang kita terjemahkan ke error models.
const (
	pgForeignKeyViolation = "23503" // id yang direferensikan tidak ada / data masih dipakai
	pgExclusionViolation  = "23P01" // EXCLUDE constraint: jadwal overlap di studio yang sama
)

const showtimeColumns = `id, movie_id, studio_id, start_time, end_time, price, status, created_at, updated_at`

type ShowtimeRepo struct {
	db *pgxpool.Pool
}

func NewShowtimeRepo(db *pgxpool.Pool) *ShowtimeRepo {
	return &ShowtimeRepo{db: db}
}

func (r *ShowtimeRepo) Create(ctx context.Context, s *models.Showtime) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // kalau sudah Commit, Rollback ini tidak melakukan apa-apa

	// end_time dihitung dari durasi film. Kalau movie tidak ada, SELECT kosong -> ErrNoRows.
	err = tx.QueryRow(ctx, `
		INSERT INTO showtimes (movie_id, studio_id, start_time, end_time, price)
		SELECT m.id, $2::bigint, $3::timestamptz, $3::timestamptz + make_interval(mins => m.duration_minutes), $4::int
		FROM movies m WHERE m.id = $1
		RETURNING id, end_time, status, created_at, updated_at`,
		s.MovieID, s.StudioID, s.StartTime, s.Price,
	).Scan(&s.ID, &s.EndTime, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrInvalidReference
	}
	if err != nil {
		return mapPgError(err, models.ErrInvalidReference)
	}

	// Buat stok tiket: 1 baris per kursi di studio ini, semuanya 'available'.
	if _, err := tx.Exec(ctx,
		`INSERT INTO showtime_seats (showtime_id, seat_id) SELECT $1, id FROM seats WHERE studio_id = $2`,
		s.ID, s.StudioID,
	); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *ShowtimeRepo) GetByID(ctx context.Context, id int64) (*models.Showtime, error) {
	var s models.Showtime
	err := scanShowtime(r.db.QueryRow(ctx, `SELECT `+showtimeColumns+` FROM showtimes WHERE id = $1`, id), &s)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// shortcut: tanpa pagination/filter, tambahkan ?page&date&cinema_id kalau data jadwal sudah ribuan.
func (r *ShowtimeRepo) List(ctx context.Context) ([]models.Showtime, error) {
	rows, err := r.db.Query(ctx, `SELECT `+showtimeColumns+` FROM showtimes ORDER BY start_time`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Showtime{} // bukan nil, supaya JSON-nya [] bukan null
	for rows.Next() {
		var s models.Showtime
		if err := scanShowtime(rows, &s); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *ShowtimeRepo) Update(ctx context.Context, s *models.Showtime) error {
	err := scanShowtime(r.db.QueryRow(ctx, `
		UPDATE showtimes s SET
			movie_id   = m.id,
			start_time = $3::timestamptz,
			end_time   = $3::timestamptz + make_interval(mins => m.duration_minutes),
			price      = $4,
			status     = $5,
			updated_at = now()
		FROM movies m
		WHERE s.id = $1 AND m.id = $2
		RETURNING s.id, s.movie_id, s.studio_id, s.start_time, s.end_time, s.price, s.status, s.created_at, s.updated_at`,
		s.ID, s.MovieID, s.StartTime, s.Price, s.Status,
	), s)
	if errors.Is(err, pgx.ErrNoRows) {
		// Tidak ada baris ter-update: entah jadwalnya tidak ada, atau movie_id-nya tidak ada.
		var exists bool
		if err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM showtimes WHERE id = $1)`, s.ID).Scan(&exists); err != nil {
			return err
		}
		if exists {
			return models.ErrInvalidReference
		}
		return models.ErrNotFound
	}
	return mapPgError(err, models.ErrInvalidReference)
}

func (r *ShowtimeRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM showtimes WHERE id = $1`, id)
	if err != nil {
		// FK dari bookings menolak delete kalau jadwal sudah punya transaksi.
		return mapPgError(err, models.ErrShowtimeHasBooking)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func scanShowtime(row pgx.Row, s *models.Showtime) error {
	return row.Scan(&s.ID, &s.MovieID, &s.StudioID, &s.StartTime, &s.EndTime, &s.Price, &s.Status, &s.CreatedAt, &s.UpdatedAt)
}

// mapPgError mengubah error constraint PostgreSQL jadi error models.
// onForeignKey dipilih pemanggil karena arti FK violation beda di insert vs delete.
func mapPgError(err error, onForeignKey error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgForeignKeyViolation:
			return onForeignKey
		case pgExclusionViolation:
			return models.ErrScheduleConflict
		}
	}
	return err
}