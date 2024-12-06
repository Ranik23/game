package main

import (
	"game/internal/config"
	"game/internal/server"
	"game/internal/storage/postgres"
	"game/internal/storage/redis"
	"game/internal/usecase"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig(filepath.Join(os.Getenv("HOME"), "game/config/config.yaml"))
	if err != nil {
		log.Fatalf("failed to load the config: %v", err)
	}


	redisClient := redis.NewRedisClient(cfg.Redis.Addr,
					cfg.Redis.Port,
					cfg.Redis.Password,
					cfg.Redis.DB)

	posstgresClient, err := postgres.NewPostgresClient(cfg.Postgres.Host, 
						cfg.Postgres.Port,
						cfg.Postgres.User,
						cfg.Postgres.Password,
						cfg.Postgres.DbName)
	if err != nil {
		log.Fatalf("failed to create a postgres client :%v", err)
	}


	Router := gin.Default()

	Router.LoadHTMLGlob(filepath.Join(os.Getenv("HOME"), "game/internal/static/*.html"))

	Logger := slog.Default()

	userOperator := usecase.NewUseCase(posstgresClient, redisClient, Logger)

	server := server.NewServer(cfg, Logger, Router, userOperator)
								
	server.Run()

}