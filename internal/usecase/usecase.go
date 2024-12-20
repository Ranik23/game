// TODO: НУЖНО ОКОНЧАТЕЛЬНО РЕШИТЬ ГДЕ МЫ ХРАНИМ ДАННЫЕ ОБ ИГРОКАХ И ТАК ДАЛЕЕ. ЛИБО В СТРУКТУРЕ КАК СЛАЙС. ЛИБО БЕРЕМ ИЗ БД
// ИНАЧЕ ЭТОТ СЛОЙ ПОЧТИ НИКАК НЕ ВЗАИМОДЕЙСТВУЕТ С БД СЛОЕМ.


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
	ErrTeamAlreadyExists  = errors.New("team already exists")
)

const (
	RedisKeyAdminLogged = "admin_logged"
	MaxPlayers          = 9
)

type UseCase interface {
	AddPlayer(*models.Team, *models.Player) error
	AddAdmin(*models.Admin) error
	CountPlayers() (int, error)
	IsAdminLoggedIn() (bool, error)
	PlayersNumberExceeded() (bool, error)
	RemovePlayer(playerID int) error
	GetPlayers() ([]models.Player, error)
	GetTeams()   ([]models.Team, error)
	AddLoginInfo(login, password string) error
	CheckLoginInfo(login string, password string) error
	CreateTeam(teamName string) (*models.Team, error)
	CountTeams() (int, error)
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

// TODO: сюда надо и команду добавить в аргументы
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

	_, err = uc.postgresClient.GetHash(login)
	if err != nil {
		uc.logger.Error("Failed to get hash", "error", err)
		return err
	}

	return nil //bcrypt.CompareHashAndPassword(hash, []byte(password))
	// TODO: я сделал пока что так, потому что у нас не созданы таблицы до конца
	// И ВОЗМОЖНО НУЖНО СПУСТИТЬ ЭТУ ФУНКЦИЮ ВНИЗ И ПОСТАВИТЬ ЗАГЛУШКУ НА ВСЕГДА NIL, ЧТОБЫ НЕ УРОДОВАТЬ USECASE
}

func (uc *useCaseImpl) CreateTeam(teamName string) (*models.Team, error) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	for _, team := range uc.Admin.Teams {
		if team.Name == teamName {
			uc.logger.Warn("Team already exists", "team", teamName)
			return nil, ErrTeamAlreadyExists
		}
	}

	uc.Admin.Teams = append(uc.Admin.Teams, &models.Team{
		Name: teamName,
	})

	return nil, nil
}


func (uc *useCaseImpl) GetTeams() ([]models.Team, error) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	var teams []models.Team

	if uc.Admin == nil {
		uc.logger.Error("Admin not set")
		return nil, ErrNoAdminSet
	}

	for _, team := range uc.Admin.Teams {
		teams = append(teams, *team)
	}

	return teams, nil
}


func (uc *useCaseImpl) CountTeams() (int, error) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()
	return len(uc.Admin.Teams), nil
}
