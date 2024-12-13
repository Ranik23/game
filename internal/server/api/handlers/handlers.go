package handlers

import (
	"game/internal/usecase"
	"net/http"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WelcomeHandler(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(g *gin.Context) {
		session := sessions.Default(g)
		session.Set("home_visited", "true")
		session.Save()
		g.HTML(http.StatusOK, "welcome.html", gin.H{})
	}
}

func RoleHandler(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(g *gin.Context) {
		session := sessions.Default(g)
		session.Set("roleSelection_visited", "true")
		session.Save()
		g.HTML(http.StatusOK, "role.html", gin.H{})
	}
}

func LoginHandlerGET(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "login.html", gin.H{})
	}
}

func LoginHandlerPOST(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "admin" && password == "123" {
			session := sessions.Default(c)
			session.Set("login_visited", true)
			session.Save()

			c.JSON(http.StatusOK, gin.H{
				"message":  "login success",
				"redirect": "/role/admin-panel",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_credentials",
			})
		}
	}
}

func MainHandler(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "player-panel.html", gin.H{})
	}
}

func AdminMainHandler(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "admin-panel.html", gin.H{})
	}
}

func LogoutHandler(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/role/login")
	}
}
