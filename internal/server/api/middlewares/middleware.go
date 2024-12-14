package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func EnsureHomeVisited() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		home_visited := session.Get("home_visited")
		if home_visited != true {
			log.Println("/home is not visited")
			c.Redirect(http.StatusFound, "/home")
			c.Abort()
			return
		}
		c.Next()
	}
}

func EnsureRoleSelectionVisited() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		roleSelection_visited := session.Get("roleSelection_visited")
		if roleSelection_visited != true {
			log.Println("/role is not visited")
			c.Redirect(http.StatusFound, "/role")
			c.Abort()
			return
		}
		c.Next()
	}
}


func EnsureLoginVisited() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		login_visited := session.Get("login_visited")
		if login_visited != true {
			log.Println("/role/login is not visited")
			c.Redirect(http.StatusFound, "/role/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

