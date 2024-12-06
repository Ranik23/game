package main

import (
	"game/internal/config"
	"game/internal/server"
	"game/internal/storage/postgres"
	_ "game/internal/storage/postgres"
	"game/internal/storage/redis"
	_ "game/internal/storage/redis"
	"game/internal/usecase"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {

	config, err := config.LoadConfig(filepath.Join(os.Getenv("HOME"), "game/config/config.yaml"))
	if err != nil {
		log.Fatalf("failed to load the config: %v", err)
	}


	redisClient := redis.NewRedisClient(config.Redis.Addr,
					config.Redis.Port,
					config.Redis.Password,
					config.Redis.DB)

	posstgresClient, err := postgres.NewPostgresClient(config.Postgres.Host, 
						config.Postgres.Port,
						config.Postgres.User,
						config.Postgres.Password,
						config.Postgres.DbName)
	if err != nil {
		log.Fatalf("failed to create a postgres client :%v", err)
	}


	Router := gin.Default()

	Router.LoadHTMLGlob(filepath.Join(os.Getenv("HOME"), "home/anton/game/internal/static/*.html"))

	Logger := slog.Default()

	userOperator := usecase.NewUseCase(posstgresClient, redisClient, Logger)

	server := server.NewServer(config, Logger, Router, userOperator)
								
	server.Run()

}