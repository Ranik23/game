package postgres

import (
	"errors"
	"game/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage interface {
	Insert(data interface{}) error
	CheckLoginExists(login string) (bool, error)
	GetHash(login string) ([]byte, error)
	//		Get(key string) (interface{}, error)
	//		Update(key string, data interface{}) error
	//		Delete(key string) error
	//	}
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

func (s *PostgresClient) Insert(data interface{}) error {
	if err := s.db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func (s *PostgresClient) CheckLoginExists(login string) (bool, error) {

	var login_info models.LoginInfo

	if err := s.db.Where("login = ?", login).First(&login_info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		return false, err
	}

	return true, nil
}

func (s *PostgresClient) GetHash(login string) ([]byte, error) {
	var login_info models.LoginInfo
	if err := s.db.Where("login = ?", login).First(&login_info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return login_info.Hash, nil
}
