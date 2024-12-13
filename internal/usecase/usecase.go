package usecase

import (
	"context"
	"errors"
	"fmt"
	"game/internal/models"
	"game/internal/storage/postgres"
	"game/internal/storage/redis"
	"log/slog"
	"sync"
)

// TODO: возможно в структуре useCase хранить указатель на админа и игроков


var (
	ErrNilPtr = fmt.Errorf("admin is nil")
)

type UseCase interface {
    AddPlayer(*models.Player) 				error
    AddAdmin(*models.Admin) 				error
    CountPlayers() 							int
    IsAdminLoggedIn() 						(bool, error)
    PlayersNumberExceeded() 				(bool, error)
	RemovePlayer(playerID int) 				error
	GetPlayers() 							([]*models.Player, error)
}

type useCaseImpl struct {
    postgresClient postgres.Storage
    redisClient    redis.Redis
    logger         *slog.Logger
    Admin          *models.Admin
    mutex          sync.Mutex
}

func NewUseCase(postgres *postgres.PostgresClient,
    redis *redis.RedisClient,
    logger *slog.Logger) *useCaseImpl {
    return &useCaseImpl{
        postgresClient: postgres,
        redisClient:    redis,
        logger:         logger,
        Admin:          nil,
    }
}

func (uc *useCaseImpl) CountPlayers() int {
    uc.mutex.Lock()
    defer uc.mutex.Unlock()

    if uc.Admin == nil {
        return 0
    }
    return len(uc.Admin.Players)
}

func (uc *useCaseImpl) AddPlayer(player *models.Player) error {
    uc.mutex.Lock()
    defer uc.mutex.Unlock()

    if uc.Admin == nil {
        uc.logger.Error("Cannot add player: no admin is set")
        return errors.New("no admin is set")
    }

    uc.Admin.Players = append(uc.Admin.Players, player)
    uc.logger.Info("Player added successfully", "player", player.UserName)
    return nil
}

func (uc *useCaseImpl) AddAdmin(admin *models.Admin) error {
    uc.mutex.Lock()
    defer uc.mutex.Unlock()

    uc.Admin = admin
    uc.logger.Info("Admin added successfully", "admin", admin.Name)

    return uc.redisClient.Set(context.Background(), "admin_logged", "true")
}

func (uc *useCaseImpl) IsAdminLoggedIn() (bool, error) {
    ctx := context.Background()

    logged, err := uc.redisClient.Get(ctx, "admin_logged")
    if err != nil {
        uc.logger.Error("Failed to fetch admin_logged status from Redis", "error", err)
        return false, err
    }

    return logged == "true", nil
}

func (uc *useCaseImpl) PlayersNumberExceeded() (bool, error) {
    uc.mutex.Lock()
    defer uc.mutex.Unlock()

    if uc.Admin == nil {
        return false, errors.New("no admin is set")
    }

	maxPlayers := 9

    return len(uc.Admin.Players) > maxPlayers, nil
}


func (uc *useCaseImpl) RemovePlayer(playerID int) error {
    uc.mutex.Lock()
    defer uc.mutex.Unlock()

    if uc.Admin == nil {
        return errors.New("no admin is set")
    }

    for i, player := range uc.Admin.Players {
        if player.ID == playerID {
            uc.Admin.Players = append(uc.Admin.Players[:i], uc.Admin.Players[i+1:]...)
            uc.logger.Info("Player removed successfully", "playerID", playerID)
            return nil
        }
    }

    return errors.New("player not found")
}

func (uc *useCaseImpl) GetPlayers() ([]*models.Player, error) {
    uc.mutex.Lock()
    defer uc.mutex.Unlock()

    if uc.Admin == nil {
        return nil, errors.New("no admin is set")
    }

    return uc.Admin.Players, nil
}
