package db

import (
	"context"
	"log"
	"time"

	"github.com/critiq17/tripListik/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init() {

	cfg := config.LoadConfig()

	dbUrl := cfg.DB_URL

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbUrl)

	if err != nil {
		log.Fatalf("DB unavailable: %v", err)
	}

	DB = pool

	log.Printf("Successful connection to PostgreSQL %s", dbUrl)
}
