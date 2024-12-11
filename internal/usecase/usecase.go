package usecase

import (
	//"game/internal/models"
	"context"
	"game/internal/models"
	"game/internal/storage/postgres"
	"game/internal/storage/redis"
	"log/slog"
	"sync"
)

// TODO: возможно в структуре useCase хранить указатель на админа и игроков

type UseCase interface {
	AddPlayer(*models.Player) error
	AddAdmin(*models.Admin) error
	CountPlayers() int
	IsAdminLoggedIn() (bool, error)
	PlayersNumberExceeded() (bool, error)
}


type userOperator struct {
	postgresClient postgres.Storage
	redisClient	   redis.Redis
	logger 		   *slog.Logger
	Admin 		   *models.Admin
	mutex sync.Mutex
}

func NewUseCase(postgres *postgres.PostgresClient,
				redis *redis.RedisClient,
				logger *slog.Logger) *userOperator {
	return &userOperator{
		postgresClient: postgres,
		redisClient:    redis,
		logger: logger,
	}
}

func (operator *userOperator) CountPlayers() int {
	operator.mutex.Lock()
	defer operator.mutex.Unlock()
	return len(operator.Admin.Players)
}

func (operator *userOperator) AddPlayer(player *models.Player) error {
	operator.mutex.Lock()
	defer operator.mutex.Unlock()
	operator.Admin.Players = append(operator.Admin.Players, player)
	return nil
}


func (operator* userOperator) AddAdmin(admin *models.Admin) error {
	operator.mutex.Lock()
	defer operator.mutex.Unlock()

	operator.Admin = admin 

	return operator.redisClient.Set(context.Background(), "admin_logged", "true") // TODO: а если ошибка тут, а админа мы уже добавили
}


func (operator *userOperator) IsAdminLoggedIn() (bool, error) {
	operator.mutex.Lock()
	defer operator.mutex.Unlock()

	if operator.Admin != nil {
		return true, nil
	} else {
		return false, nil
	} 
	// TODO: пока убрать ошибку в возвращаемом значении
}

func (operator *userOperator) PlayersNumberExceeded() (bool, error) {
	operator.mutex.Lock()
	defer operator.mutex.Unlock()
	return false, nil
}


