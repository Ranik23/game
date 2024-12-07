package usecase

import (
	//"game/internal/models"
	"context"
	"game/internal/models"
	"game/internal/storage/postgres"
	"game/internal/storage/redis"
	"log/slog"
	"strconv"
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


type useCase struct {
	postgresClient postgres.Storage
	redisClient	   redis.Redis
	logger 		   *slog.Logger
	Admin 		   *models.Admin
	Players		   []*models.Player
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

func (u *useCase) CountPlayers() int {
	return len(u.Players)
}

func (u *useCase) AddPlayer(player *models.Player) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Players = append(u.Players, player)
	return nil
}


func (u* useCase) AddAdmin(admin *models.Admin) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.Admin = admin 

	return u.redisClient.Set(context.Background(), "admin_logged", "true") // TODO: а если ошибка тут, а админа мы уже добавили
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


