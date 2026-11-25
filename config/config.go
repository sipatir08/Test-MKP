package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func Load() Config {
	// .env opsional: di Docker/server, env biasanya diset langsung tanpa file
	if err := godotenv.Load(); err != nil {
		log.Println("tidak ada .env, pakai environment variable sistem")
	}

	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL wajib diisi")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET wajib diisi")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	return cfg
}