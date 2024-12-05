package usecase

import (
	"game/internal/models"
	"game/internal/session"
	"game/internal/storage/postgres"
	"game/internal/storage/redis"
	"log/slog"
)





type UseCase interface {
	AddAdminToSession(admin *models.Admin) error
	AddPlayerToSession(player *models.Player) error
}


type useCase struct {
	postgresClient postgres.Storage
	redisClient	   redis.Redis
	logger 		   *slog.Logger
	gameSession    *session.GameSession
}

