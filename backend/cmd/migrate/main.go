package main

import (
	"flag"
	"log"

	"github.com/critiq17/tripListik/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	direction := flag.String("direction", "up", "up or down")
	steps := flag.Int("steps", 0, "number of steps to run (0 = all)")
	flag.Parse()

	cfg := config.LoadConfig()
	if cfg.DSN == "" {
		log.Fatal("DSN is required")
	}

	m, err := migrate.New("file://migrations", cfg.DSN)
	if err != nil {
		log.Fatalf("migrate init error: %v", err)
	}

	switch *direction {
	case "up":
		if *steps > 0 {
			if err := m.Steps(*steps); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("migrate up error: %v", err)
			}
			return
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate up error: %v", err)
		}
	case "down":
		if *steps > 0 {
			if err := m.Steps(-*steps); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("migrate down error: %v", err)
			}
			return
		}
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate down error: %v", err)
		}
	default:
		log.Fatalf("unknown direction: %s", *direction)
	}
}
