package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func Init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file")
	}

	dbUrl := os.Getenv("DATABASE_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbUrl)

	if err != nil {
		log.Fatalf("DB unavailable: %v", err)
	}

	DB = pool

	log.Println("Successful connection to PostgreSQL", dbUrl)
}
