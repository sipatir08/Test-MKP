-- =====================================================================
-- Cinema Ticket Platform - PostgreSQL schema
-- =====================================================================

-- btree_gist dipakai untuk EXCLUDE constraint (cegah jadwal bentrok di studio yang sama)
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- ---------------------------------------------------------------------
-- USERS
-- ---------------------------------------------------------------------
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(150) NOT NULL UNIQUE,
    phone         VARCHAR(20),
    password_hash TEXT         NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'customer' CHECK (role IN ('admin', 'customer')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------
-- MASTER DATA: kota -> bioskop (cabang) -> studio -> kursi
-- ---------------------------------------------------------------------
CREATE TABLE cities (
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE cinemas (
    id         BIGSERIAL PRIMARY KEY,
    city_id    BIGINT       NOT NULL REFERENCES cities (id),
    name       VARCHAR(100) NOT NULL,
    address    TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX idx_cinemas_city ON cinemas (city_id);

CREATE TABLE studios (
    id          BIGSERIAL PRIMARY KEY,
    cinema_id   BIGINT      NOT NULL REFERENCES cinemas (id),
    name        VARCHAR(50) NOT NULL,
    studio_type VARCHAR(20) NOT NULL DEFAULT 'regular' CHECK (studio_type IN ('regular', 'premium', 'imax')),
    capacity    INT         NOT NULL CHECK (capacity > 0),
    UNIQUE (cinema_id, name)
);

-- Denah kursi fisik per studio (tidak berubah per jadwal)
CREATE TABLE seats (
    id          BIGSERIAL PRIMARY KEY,
    studio_id   BIGINT      NOT NULL REFERENCES studios (id) ON DELETE CASCADE,
    row_label   VARCHAR(3)  NOT NULL,
    seat_number INT         NOT NULL CHECK (seat_number > 0),
    seat_type   VARCHAR(20) NOT NULL DEFAULT 'regular' CHECK (seat_type IN ('regular', 'vip', 'couple')),
    UNIQUE (studio_id, row_label, seat_number)
);

CREATE TABLE movies (
    id               BIGSERIAL PRIMARY KEY,
    title            VARCHAR(200) NOT NULL,
    synopsis         TEXT,
    genre            VARCHAR(50),
    duration_minutes INT          NOT NULL CHECK (duration_minutes > 0),
    age_rating       VARCHAR(10)  NOT NULL DEFAULT 'SU', -- SU, 13+, 17+, 21+ (LSF)
    release_date     DATE,
    poster_url       TEXT,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------
-- JADWAL TAYANG
-- ---------------------------------------------------------------------
CREATE TABLE showtimes (
    id         BIGSERIAL PRIMARY KEY,
    movie_id   BIGINT      NOT NULL REFERENCES movies (id),
    studio_id  BIGINT      NOT NULL REFERENCES studios (id),
    start_time TIMESTAMPTZ NOT NULL,
    end_time   TIMESTAMPTZ NOT NULL,
    price      INT         NOT NULL CHECK (price > 0), -- rupiah, integer supaya tidak ada error pembulatan
    status     VARCHAR(20) NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (end_time > start_time),
    -- Satu studio tidak boleh punya dua jadwal aktif yang waktunya overlap
    EXCLUDE USING gist (studio_id WITH =, tstzrange(start_time, end_time) WITH &&)
        WHERE (status = 'scheduled')
);
CREATE INDEX idx_showtimes_movie_start ON showtimes (movie_id, start_time);
CREATE INDEX idx_showtimes_studio_start ON showtimes (studio_id, start_time);

-- Inventori kursi per jadwal: 1 baris = 1 kursi di 1 jadwal.
-- Inilah "stok tiket". Dibuat otomatis saat jadwal dibuat.
CREATE TABLE showtime_seats (
    id                BIGSERIAL PRIMARY KEY,
    showtime_id       BIGINT      NOT NULL REFERENCES showtimes (id) ON DELETE CASCADE,
    seat_id           BIGINT      NOT NULL REFERENCES seats (id),
    status            VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'locked', 'sold')),
    locked_by_user_id BIGINT REFERENCES users (id),
    locked_until      TIMESTAMPTZ,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (showtime_id, seat_id),
    CHECK (status <> 'locked' OR locked_until IS NOT NULL)
);
CREATE INDEX idx_showtime_seats_status ON showtime_seats (showtime_id, status);

-- ---------------------------------------------------------------------
-- TRANSAKSI
-- ---------------------------------------------------------------------
CREATE TABLE bookings (
    id           BIGSERIAL PRIMARY KEY,
    booking_code VARCHAR(20) NOT NULL UNIQUE,
    user_id      BIGINT      NOT NULL REFERENCES users (id),
    showtime_id  BIGINT      NOT NULL REFERENCES showtimes (id),
    status       VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'paid', 'expired', 'cancelled', 'refunded')),
    total_amount INT         NOT NULL CHECK (total_amount >= 0),
    expires_at   TIMESTAMPTZ NOT NULL, -- batas waktu bayar (= locked_until kursi)
    paid_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_bookings_user ON bookings (user_id, created_at DESC);
CREATE INDEX idx_bookings_showtime_status ON bookings (showtime_id, status);
CREATE INDEX idx_bookings_pending_expiry ON bookings (expires_at) WHERE status = 'pending';

-- Riwayat kursi per booking. Tetap tersimpan walau kursi sudah di-restock.
CREATE TABLE booking_seats (
    id               BIGSERIAL PRIMARY KEY,
    booking_id       BIGINT NOT NULL REFERENCES bookings (id) ON DELETE CASCADE,
    showtime_seat_id BIGINT NOT NULL REFERENCES showtime_seats (id),
    price            INT    NOT NULL CHECK (price >= 0),
    UNIQUE (booking_id, showtime_seat_id)
);

CREATE TABLE payments (
    id           BIGSERIAL PRIMARY KEY,
    booking_id   BIGINT       NOT NULL REFERENCES bookings (id),
    method       VARCHAR(30)  NOT NULL, -- qris, va_bca, ewallet_ovo, ...
    amount       INT          NOT NULL CHECK (amount >= 0),
    status       VARCHAR(20)  NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'success', 'failed', 'refunded')),
    provider_ref VARCHAR(100) UNIQUE, -- id transaksi dari payment gateway, untuk idempotensi webhook
    paid_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX idx_payments_booking ON payments (booking_id);

CREATE TABLE refunds (
    id           BIGSERIAL PRIMARY KEY,
    booking_id   BIGINT      NOT NULL UNIQUE REFERENCES bookings (id), -- 1 booking maksimal 1 refund (idempoten)
    payment_id   BIGINT      NOT NULL REFERENCES payments (id),
    amount       INT         NOT NULL CHECK (amount >= 0),
    reason       TEXT        NOT NULL,
    initiated_by VARCHAR(20) NOT NULL CHECK (initiated_by IN ('cinema', 'customer')),
    status       VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'processing', 'success', 'failed')),
    processed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_refunds_status ON refunds (status) WHERE status IN ('pending', 'processing');