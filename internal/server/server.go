package server

import (
	"game/internal/config"
	"game/internal/server/api/handlers"
	"game/internal/server/api/middlewares"
	"game/internal/usecase"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	config  *config.Config
	logger  *slog.Logger
	UserOperator usecase.UseCase
}

func NewServer(config *config.Config, logger *slog.Logger, router *gin.Engine, usecase usecase.UseCase) *Server {
	server := &Server{
		router:  router,
		config:  config,
		logger:  logger,
		UserOperator: usecase,
	}

	server.setUpMiddlewares()
	server.setUpRoutes()
	server.setUpHTMLFiles(os.Getenv("HOME") + "/game/internal/static/*.html")
	server.setUpStaticFiles()

	return server
}

func (s *Server) setUpRoutes() {
	s.router.GET("/", func(g *gin.Context) {
		g.Redirect(http.StatusFound, "/home")
	})

	s.router.GET("/home", handlers.WelcomeHandler(s.UserOperator))
	s.router.GET("/home/role", handlers.RoleHandler(s.UserOperator))
	s.router.GET("/home/role/login", handlers.LoginHandlerGET(s.UserOperator))
	s.router.POST("/home/role/login", handlers.LoginHandlerPOST(s.UserOperator))
	s.router.GET("/home/role/guest-panel", handlers.MainHandler(s.UserOperator))
	s.router.GET("/home/role/admin-panel", handlers.AdminMainHandler(s.UserOperator))
	s.router.GET("/ws-guest", handlers.ClientWebSocketHandler(s.UserOperator))
	s.router.GET("/ws-admin", handlers.AdminWebSocketHandler(s.UserOperator))
}

func (s *Server) setUpMiddlewares() {
	s.router.Use(gin.Recovery())
	s.router.Use(gin.Logger())
	s.router.Use(middlewares.Ip())
}

func (s *Server) setUpHTMLFiles(pattern string) {
	s.router.LoadHTMLGlob(pattern)
}

func (s *Server) setUpStaticFiles() {
	s.router.Static("/static", filepath.Join(os.Getenv("HOME"), "game/internal/static/scripts"))
}

func (s *Server) Run() {
	addr := s.config.Localhost.Host + ":" + s.config.Localhost.Port
	if err := s.router.Run(addr); err != nil {
		s.logger.Error("failed to start the server", "error", err)
		return
	}
	s.logger.Info("server started successfully", "address", addr)
}