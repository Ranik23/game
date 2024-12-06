package main

import (
	"game/internal/config"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)




func main() {
	cfg, err := config.LoadConfig(filepath.Join(os.Getenv("HOME"), "game/config/config.yaml"))

	if err != nil {
		log.Fatalf("failed to load the config: %v", err)
	}

	dsn := "host=" + cfg.Postgres.Host + " port=" + cfg.Localhost.Port + " user=" + cfg.Postgres.User + " password=" + cfg.Postgres.Password + " dbname=" + cfg.Postgres.DbName + " sslmode=disable"


    m, err := migrate.New("file://" + os.Getenv("HOME") + "/game/internal/storage/postgres/migrations", dsn)
    if err != nil {
        log.Fatalf("failed to load the path to migrations: %v", err)
    }
    if err := m.Up(); err != nil {
        log.Fatalf("failed to migrate: %v", err)
    }
}