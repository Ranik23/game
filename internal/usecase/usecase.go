package usecase

import (
	//"game/internal/models"
	"context"
	"game/internal/storage/postgres"
	"game/internal/storage/redis"
	"log/slog"
	"strconv"
	"sync"
)

type UseCase interface {
	AddPlayer() error
	AddAdmin() error
	IsAdminLoggedIn() (bool, error)
	PlayersNumberExceeded() (bool, error)
}


type useCase struct {
	postgresClient postgres.Storage
	redisClient	   redis.Redis
	logger 		   *slog.Logger
	mu sync.Mutex
}

func NewUseCase(postgres *postgres.PostgresClient,
				redis *redis.RedisClient,
				logger *slog.Logger) *useCase {
	return &useCase{
		postgresClient: postgres,
		redisClient:    redis,
		logger: logger,
	}
}


func (u *useCase) AddPlayer() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	number, err := u.redisClient.Get(context.Background(), "players_number")
	if err != nil {
		return err
	}

	n, err := strconv.Atoi(number)
	if err != nil {
		return err
	}

	n += 1

	if err := u.redisClient.Set(context.Background(), "players_number", strconv.Itoa(n)); err != nil {
		return err
	}

	return nil
}


func (u* useCase) AddAdmin() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.redisClient.Set(context.Background(), "admin_logged", "true")
}


func (u *useCase) IsAdminLoggedIn() (bool, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	loggedIn, err := u.redisClient.Get(context.Background(), "admin_logged_in") // TODO: возможно mutex и не нужен, если redis клиент thread-safe
	if err != nil {
		return false, err
	}
	return loggedIn == "true", nil
}

func (u *useCase) PlayersNumberExceeded() (bool, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	
	exceeded, err := u.redisClient.Get(context.Background(), "players_number") // TODO: redis хранить по ключу только строки?
	if err != nil {
		return false, err
	}

	e, err := strconv.Atoi(exceeded) 
	if err != nil {
		return false, err 
	}
	if e >= 9 {
		return true, nil
	}

	return false, nil
}


