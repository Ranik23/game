package server

import (
	"game/internal/config"
	"game/internal/server/api/handlers"
	"game/internal/usecase"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router 			*gin.Engine
	config 			*config.Config
	logger 			*slog.Logger
	usecase			usecase.UseCase
}

func NewServer(config *config.Config,
				logger *slog.Logger,
				router *gin.Engine,
				usecase usecase.UseCase) *Server {

    server := &Server{
        router:         router,
        config:         config,
        logger:         logger,
		usecase: 		usecase,	
    }

    server.setUpMiddleWares()
    server.setUpRoutes()   
	server.setUpHTMLFiles(os.Getenv("HOME") + "/game/internal/static/*.html")
	server.setUpStaticFiles()

    return server
}

func (s *Server) setUpRoutes() {

	s.router.GET("/", func(g *gin.Context) {
		g.Redirect(http.StatusFound, "/welcome")
	})
	s.router.GET("/auth", handlers.AuthHandler(s.usecase))
	s.router.POST("/auth", handlers.AuthPostHandler(s.usecase))
	
	s.router.GET("/welcome", handlers.WelcomeHandler(s.usecase))
	s.router.GET("/main", handlers.MainHandler(s.usecase))
	s.router.GET("/admin_main", handlers.AdminMainHandler(s.usecase))


	s.router.GET("/ws", handlers.WebSocketHandler(s.usecase))
	s.router.GET("/wsmain", handlers.WebSocketHandlerMain(s.usecase))
}


func (s *Server) setUpMiddleWares() {
	s.router.Use(gin.Recovery())
	s.router.Use(gin.Logger())	
}


func (s *Server) setUpHTMLFiles(pattern string) {
	s.router.LoadHTMLGlob(pattern)
}

func (s *Server) setUpStaticFiles() {
	s.router.Static("/static", filepath.Join(os.Getenv("HOME"), "game/internal/static/scripts"))
}

func (s *Server) Run() {
	if err := s.router.Run(s.config.Localhost.Host + ":" + s.config.Localhost.Port); err != nil {
		s.logger.Error("failed to start the server")
		return
	}
	s.logger.Info("server started succesfully")
}

