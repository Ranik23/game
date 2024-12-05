package config

import (
	"os"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Redis    RedisConfig    `yaml:"redis"`
	Postgres PostgresConfig `yaml:"postgres"`
	Localhost LocalhostConfig `yaml:"localhost"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Port	 string 	`yaml:"port"`    // Адрес Redis сервера, например: "localhost:6379"
	Password string `yaml:"password"` // Пароль для Redis, если требуется
	DB       int    `yaml:"db"`       // Номер базы данных Redis
}

type PostgresConfig struct {
	Host     string `yaml:"host"`     // Хост PostgreSQL, например: "localhost"
	Port     string    `yaml:"port"`     // Порт PostgreSQL, например: 5432
	User     string `yaml:"user"`     // Имя пользователя для подключения
	Password string `yaml:"password"` // Пароль для подключения
	DbName   string `yaml:"db_name"`  // Имя базы данных
}

type LocalhostConfig struct {
	Host string `yaml:"host"` // Хост для подключения, например: "localhost"
	Port string    `yaml:"port"` // Порт для подключения, например: 8080
}



func LoadConfig(path string) (*Config, error) {
	var config Config

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}