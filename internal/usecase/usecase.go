package usecase

import (
	"game/internal/storage/postgres"
	"game/internal/storage/redis"
	"log/slog"
)





type UseCase interface {
}




type useCase struct {
	postgresClient postgres.Storage
	redisClient	   redis.Redis
	logger 		   *slog.Logger
}

