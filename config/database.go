package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(databaseURL string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("format DATABASE_URL salah: ", err)
	}

	// pgxpool.New belum benar-benar konek, Ping yang memastikan
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("gagal konek ke database: ", err)
	}

	log.Println("database terhubung")
	return db
}