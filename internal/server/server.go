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

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	config       *config.Config
	logger       *slog.Logger
	UserOperator usecase.UseCase
}

func NewServer(config *config.Config, logger *slog.Logger, router *gin.Engine, usecase usecase.UseCase) *Server {
	server := &Server{
		router:       router,
		config:       config,
		logger:       logger,
		UserOperator: usecase,
	}

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	server.setUpRoutes()
	server.setUpHTMLFiles(os.Getenv("HOME") + "/game/internal/static/*.html")
	server.setUpStaticFiles()

	return server
}

func (s *Server) setUpRoutes() {

	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/home")
	})

	s.router.GET("/contacts", func(g *gin.Context) {
		g.HTML(http.StatusOK, "contacts.html", nil)
	}, middlewares.EnsureHomeVisited())

	s.router.GET("/about", func(g *gin.Context) {
		g.HTML(http.StatusOK, "about.html", nil)
	}, middlewares.EnsureHomeVisited())

	s.router.GET("/home", handlers.WelcomeHandler(s.UserOperator))
	s.router.GET("/role", middlewares.EnsureHomeVisited(), handlers.RoleHandler(s.UserOperator))

	s.router.GET("/role/login",
		middlewares.EnsureHomeVisited(),
		middlewares.EnsureRoleSelectionVisited(),
		handlers.LoginHandlerGET(s.UserOperator))

	s.router.POST("/create-team", handlers.CreateTeamHandler(s.UserOperator))

	s.router.POST("/role/login", handlers.LoginHandlerPOST(s.UserOperator))

	s.router.GET("/role/login-leader", handlers.LoginLeaderHandlerGET(s.UserOperator))

	s.router.GET("/ws/admin", handlers.AdminWebSocketHandler(s.UserOperator))
	s.router.GET("/ws/player", handlers.ClientWebSocketHandler(s.UserOperator))
	s.router.GET("/ws/leader", handlers.LeaderWebSocketHanlder(s.UserOperator))

	protected := s.router.Group("/role")
	protected.Use(
		middlewares.EnsureHomeVisited(),
		middlewares.EnsureRoleSelectionVisited(),
		middlewares.EnsureLoginVisited())
	{
		protected.GET("/player-panel", handlers.PlayerPanelHandler())
		protected.GET("/leader-panel", handlers.LeaderPanelHanlder())
		protected.GET("/logout", handlers.LogoutHandler())
		protected.GET("/admin-panel", handlers.AdminPanelHandler())
	}
}

func (s *Server) setUpHTMLFiles(pattern string) {
	s.router.LoadHTMLGlob(pattern)
}

func (s *Server) setUpStaticFiles() {
	staticPath := filepath.Join(os.Getenv("HOME"), "game/internal/static/scripts")
	s.router.Static("/static", staticPath)
}

func (s *Server) Run() {
	addr := s.config.Localhost.Host + ":" + s.config.Localhost.Port
	if err := s.router.Run(addr); err != nil {
		s.logger.Error("failed to start the server", "error", err)
		return
	}
	s.logger.Info("server started successfully", "address", addr)
}
