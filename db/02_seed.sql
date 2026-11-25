-- =====================================================================
-- Data contoh untuk mencoba API
-- admin@mkp.id / admin123        (role admin    -> boleh create/update/delete jadwal)
-- customer@mkp.id / customer123  (role customer -> hanya boleh lihat jadwal)
-- =====================================================================
CREATE EXTENSION IF NOT EXISTS pgcrypto; -- hanya untuk membuat hash bcrypt di seed

INSERT INTO users (name, email, phone, password_hash, role) VALUES
    ('Admin MKP',    'admin@mkp.id',    '081200000001', crypt('admin123', gen_salt('bf', 10)), 'admin'),
    ('Customer MKP', 'customer@mkp.id', '081200000002', crypt('customer123', gen_salt('bf', 10)), 'customer');

INSERT INTO cities (name) VALUES ('Jakarta'), ('Bandung');

INSERT INTO cinemas (city_id, name, address) VALUES
    (1, 'MKP Cinema Grand Indonesia', 'Jl. M.H. Thamrin No.1, Jakarta Pusat'),
    (2, 'MKP Cinema Paris Van Java',  'Jl. Sukajadi No.131-139, Bandung');

INSERT INTO studios (cinema_id, name, studio_type, capacity) VALUES
    (1, 'Studio 1', 'regular', 50),
    (1, 'Studio 2', 'premium', 50),
    (2, 'Studio 1', 'regular', 50);

-- Tiap studio: baris A-E x kursi 1-10 = 50 kursi
INSERT INTO seats (studio_id, row_label, seat_number)
SELECT s.id, r.row_label, n
FROM studios s
CROSS JOIN (VALUES ('A'), ('B'), ('C'), ('D'), ('E')) AS r(row_label)
CROSS JOIN generate_series(1, 10) AS n;

INSERT INTO movies (title, synopsis, genre, duration_minutes, age_rating, release_date) VALUES
    ('Pengabdi Setan 3', 'Teror keluarga Suwono berlanjut.', 'Horror', 120, '17+', '2026-11-01'),
    ('Laskar Pelangi Reborn', 'Kisah sepuluh anak Belitung.', 'Drama', 125, 'SU', '2026-10-15');