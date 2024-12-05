package main

import (
	"game/internal/config"
	"game/internal/server"
	_"game/internal/storage/postgres"
	_"game/internal/storage/redis"
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


	// redisClient := redis.NewRedisClient(config.Redis.Addr,
	// 				config.Redis.Port,
	// 				config.Redis.Password,
	// 				config.Redis.DB)

	// posstgresClient, err := postgres.NewPostgresClient(config.Postgres.Host, 
	// 					config.Postgres.Port,
	// 					config.Postgres.User,
	// 					config.Postgres.Password,
	// 					config.Postgres.DbName)
	// if err != nil {
	// 	log.Fatalf("failed to create a postgres client :%v", err)
	// }


	Router := gin.New()

	Router.LoadHTMLGlob("/home/anton/game/internal/static/*.html")

	Logger := &slog.Logger{}

	server := server.NewServer(config, Logger, Router, nil)
								
	server.Run()

}