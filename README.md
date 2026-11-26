# Backend Development Test: Mitra Kasih Perkasa 2025

| | |
|---|---|
| **Nama Lengkap** | Muhammad Fathir |
| **Date** | 26 November 2026 |
| **Email** | fathir080604@gmail.com |
| **Phone** | +62 856-2448-2521 |
| **Repository** | https://github.com/sipatir08/Test-MKP |

Platform pembelian tiket bioskop online skala nasional.

Alur booking, kunci kursi, dan refund dijawab di bagian A (desain) dan sudah disiapkan tabelnya di bagian B (database). Bagian C mengimplementasikan API sesuai soal: login dan CRUD jadwal tayang.

| Bagian | Lokasi |
|---|---|
| **A. System Design** | [`docs/system-design.md`](docs/system-design.md), [`docs/topology.jpg`](docs/topology.jpg), [`docs/flowchart-booking.jpg`](docs/flowchart-booking.jpg), [`docs/flowchart-refund.jpg`](docs/flowchart-refund.jpg) |
| **B. Database Design** | [`docs/erd.jpg`](docs/erd.jpg), [`db/01_schema.sql`](db/01_schema.sql), [`db/02_seed.sql`](db/02_seed.sql) |
| **C. Skill Test (Golang)** | source code di repo ini |
| **Postman** | [`postman_collection.json`](postman_collection.json) |

## Menjalankan

### 1. Siapkan database PostgreSQL
Buat database baru (PostgreSQL lokal atau Supabase), lalu jalankan:
```bash
psql -d cinema -f db/01_schema.sql -f db/02_seed.sql
```

### 2. Buat file `.env` di root project
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/cinema?sslmode=disable
JWT_SECRET=isi-string-acak-yang-panjang
PORT=8080
```

### 3. Jalankan
```bash
go mod tidy
go run .
```
API berjalan di `http://localhost:8080`.

### 4. Coba lewat Postman
Import `postman_collection.json`, jalankan request **Login Admin** (token otomatis tersimpan), lalu coba endpoint lainnya.

## Akun contoh

| Email | Password | Role | Akses |
|---|---|---|---|
| admin@mkp.id | admin123 | admin | Login + CRUD jadwal tayang |
| customer@mkp.id | customer123 | customer | Login + lihat jadwal (ubah data → 403) |

## Endpoint

Semua endpoint jadwal wajib header `Authorization: Bearer <access_token>` dari login.

| Method | Path | Akses | Keterangan |
|---|---|---|---|
| POST | `/api/v1/auth/login` | publik | Login, mengembalikan JWT (berlaku 24 jam) |
| GET | `/api/v1/showtimes` | user login | List jadwal tayang |
| GET | `/api/v1/showtimes/{id}` | user login | Detail jadwal |
| POST | `/api/v1/showtimes` | admin | Buat jadwal + stok kursi otomatis |
| PUT | `/api/v1/showtimes/{id}` | admin | Ubah film/jam/harga/status (`scheduled` / `cancelled`) |
| DELETE | `/api/v1/showtimes/{id}` | admin | Hapus jadwal (409 jika sudah ada booking) |

### Contoh body create jadwal
```json
{
  "movie_id": 1,
  "studio_id": 1,
  "start_time": "2027-01-10T19:00:00+07:00",
  "price": 50000
}
```
`end_time` dihitung otomatis dari durasi film.

### Format error
Semua error berbentuk `{"error": "pesan"}` dengan status:

| Status | Arti |
|---|---|
| 400 | Validasi gagal (field kosong, tanggal lewat, harga 0, field tidak dikenal) |
| 401 | Token tidak ada / tidak valid, atau email/password salah |
| 403 | Bukan admin |
| 404 | Data tidak ditemukan |
| 409 | Jadwal bentrok di studio yang sama, atau jadwal sudah punya booking |
| 422 | `movie_id` / `studio_id` tidak ada |

## Struktur Project (Clean Architecture)

```
main.go         merakit dependency: config → repository → service → controller → routes
config/         load .env & koneksi database
models/         entity, error bisnis, interface repository (tidak bergantung ke layer lain)
repository/     query PostgreSQL (pgx), satu-satunya layer yang menyentuh database
service/        logika bisnis: login, validasi jadwal
controllers/    terima request HTTP, panggil service, kirim JSON
middleware/     validasi JWT & cek role admin
routes/         daftar endpoint (router bawaan net/http)
utils/          implementasi JWT (HS256) & bcrypt
db/             schema & seed PostgreSQL
docs/           system design, ERD, flowchart, topologi
```

Arah dependency: `controllers → service → models ← repository`

**Tech:** Go 1.26+ (`net/http`), PostgreSQL, pgx v5, golang-jwt v5, bcrypt.