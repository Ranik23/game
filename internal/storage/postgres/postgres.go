package postgres

import (
	// "context"
	// "game/internal/session"
	// "time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type Storage interface {
	
}


type PostgresClient struct {
	db *gorm.DB
}

func NewPostgresClient(host, port, user, password, dbname string) (*PostgresClient, error) {
	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PostgresClient{db: db}, nil
}

