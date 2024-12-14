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

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNilPtr             = fmt.Errorf("admin is nil")
	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrNoAdminSet         = errors.New("no admin is set")
	ErrPlayerNotFound     = errors.New("player not found")
	ErrLoginNotFound      = errors.New("login not found")
	ErrAdminIsAlreadySet  = errors.New("admin is set already")
)

const (
	RedisKeyAdminLogged = "admin_logged"
	MaxPlayers          = 9
)

type UseCase interface {
	AddPlayer(*models.Player) error
	AddAdmin(*models.Admin) error
	CountPlayers() (int, error)
	IsAdminLoggedIn() (bool, error)
	PlayersNumberExceeded() (bool, error)
	RemovePlayer(playerID int) error
	GetPlayers() ([]models.Player, error)
	AddLoginInfo(login, password string) error
	CheckLoginInfo(login string, password string) error
	CreateTeam() (*models.Team, error)
}

type useCaseImpl struct {
	postgresClient postgres.Storage
	redisClient    redis.Redis
	logger         *slog.Logger
	Admin          *models.Admin
	mutex          sync.Mutex
}

func NewUseCase(
	postgres *postgres.PostgresClient,
	redis *redis.RedisClient,
	logger *slog.Logger) *useCaseImpl {
	return &useCaseImpl{
		postgresClient: postgres,
		redisClient:    redis,
		logger:         logger,
		Admin:          nil,
	}
}

func (uc *useCaseImpl) CountPlayers() (int, error) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()
	players, err := uc.GetPlayers()
	if err != nil {
		uc.logger.Error("")
		return 0, err
	}
	return len(players), nil
}

func (uc *useCaseImpl) AddPlayer(team *models.Team, player *models.Player) error {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	if uc.Admin == nil {
		uc.logger.Error("Cannot add player: no admin is set")
		return ErrNoAdminSet
	}

	team.Players = append(team.Players, player)
	uc.logger.Info("Player added successfully", "player", player.UserName)
	return nil
}

func (uc *useCaseImpl) AddAdmin(admin *models.Admin) error {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	if uc.Admin != nil {
		return ErrAdminIsAlreadySet
	}

	uc.Admin = admin
	uc.logger.Info("Admin added successfully", "admin", admin.UserName)

	return uc.redisClient.Set(context.Background(), RedisKeyAdminLogged, "true")
}

func (uc *useCaseImpl) IsAdminLoggedIn() (bool, error) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()
	ctx := context.Background()

	logged, err := uc.redisClient.Get(ctx, RedisKeyAdminLogged)
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
		return false, ErrNoAdminSet
	}

	maxPlayers := 9

	count := 0

	for _, team := range uc.Admin.Teams {
		count += len(team.Players)
	}

	return count > maxPlayers, nil
}

func (uc *useCaseImpl) RemovePlayer(playerID int) error {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	if uc.Admin == nil {
		return ErrNoAdminSet
	}

	for i, team := range uc.Admin.Teams {
		for _, player := range team.Players {
			if player.ID == playerID {
				uc.Admin.Teams[i].Players = append(uc.Admin.Teams[i].Players[:i], uc.Admin.Teams[i].Players[i+1:]...)
				uc.logger.Info("Player removed successfully", "playerID", playerID)
				return nil
		}
		}
	}

	return ErrPlayerNotFound
}

func (uc *useCaseImpl) GetPlayers() ([]models.Player, error) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	if uc.Admin == nil {
		return nil, ErrNoAdminSet
	}

	var total []models.Player

	for _, team := range uc.Admin.Teams {
		for _, player := range team.Players {
			total = append(total, *player)
		}
	}

	return total, nil
}

func (uc *useCaseImpl) AddLoginInfo(login, password string) error {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	exists, err := uc.postgresClient.CheckLoginExists(login)
	if err != nil {
		uc.logger.Error("Failed to check login existence", "error", err)
		return err
	}

	if exists {
		uc.logger.Warn("Login already exists", "login", login)
		return ErrLoginAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("Failed to hash password", "error", err)
		return err
	}

	loginInfo := models.LoginInfo{
		Login: login,
		Hash:  hash,
	}

	if err := uc.postgresClient.Insert(loginInfo); err != nil {
		uc.logger.Error("Failed to insert login info", "error", err)
		return err
	}

	uc.logger.Info("Login info added successfully", "login", login)
	return nil
}

func (uc *useCaseImpl) CheckLoginInfo(login, password string) error {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	exists, err := uc.postgresClient.CheckLoginExists(login)
	if err != nil {
		uc.logger.Error("Failed to check login existence", "error", err)
		return err
	}

	if !exists {
		uc.logger.Warn("Login does not exists", "login", login)
		return ErrLoginNotFound
	}

	hash, err := uc.postgresClient.GetHash(login)
	if err != nil {
		uc.logger.Error("Failed to get hash", "error", err)
		return err
	}

	return bcrypt.CompareHashAndPassword(hash, []byte(password))
	// TODO: я сделал пока что так, потому что у нас не созданы таблицы до конца
}

func (uc *useCaseImpl) CreateTeam() (*models.Team, error) {
	return nil, nil
}
