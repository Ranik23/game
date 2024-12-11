package middlewares

import (
	"game/internal/usecase"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		if username == nil {
			c.Redirect(http.StatusFound, "/home/role/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RoleMiddleware(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")
		if role == nil || role == "" {
			log.Println("need to choose the role first")
			c.Redirect(http.StatusFound, "/home/role")
			c.Abort()
			return
		}
		c.Next()
	}
}

func WelcomeMiddleware(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		welcome := session.Get("welcome")
		if welcome == nil {
			c.Redirect(http.StatusFound, "/home")
			c.Abort()
			return
		}
		c.Next()
	}
}
