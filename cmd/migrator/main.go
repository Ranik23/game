package main

import (
	"flag"
	"game/internal/config"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	upFlag := flag.Bool("up", true, "Apply Migrations")
	downFlag := flag.Bool("down", false, "Rollback Migrations")

	flag.Parse()

	cfg, err := config.LoadConfig(filepath.Join(os.Getenv("HOME"), "game/config/config.yaml"))

	if err != nil {
		log.Fatalf("failed to load the config: %v", err)
	}

	dsn := "postgres://" + cfg.Postgres.User + ":" + cfg.Postgres.Password +
		"@" + cfg.Postgres.Host + ":" + cfg.Postgres.Port +
		"/" + cfg.Postgres.DbName + "?sslmode=disable"

	m, err := migrate.New("file://"+os.Getenv("HOME")+"/game/internal/storage/postgres/migrations/tables", dsn)
	if err != nil {
		log.Fatalf("failed to load the path to migrations: %v", err)
	}

	if *upFlag {
		err := m.Up()
		if err != nil {
			if err == migrate.ErrNoChange {
				log.Println("no new migrations to apply")
			} else {
				log.Fatalf("failed to migrate: %v", err)
			}
		} else {
			log.Println("migrations applied successfully")
		}
	} else {
		if *downFlag {
			err := m.Down()
			if err != nil {
				log.Fatalf("failed to rollback migration: %v", err)
			}
			log.Println("migrations down successfully")
		}
	}
	
}
