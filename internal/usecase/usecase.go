package usecase

import (
	//"game/internal/models"
	"game/internal/storage/postgres"
	"game/internal/storage/redis"
	"log/slog"
	"sync"
)





type UseCase interface {
	AddPlayer() error
	AddAdmin() error
	IsAdminLoggedIn() bool
}


type useCase struct {
	postgresClient postgres.Storage
	redisClient	   redis.Redis
	logger 		   *slog.Logger
	mu sync.Mutex
}


func (u *useCase) AddPlayer() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	return nil
}


func (u* useCase) AddAdmin() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	return nil
}


