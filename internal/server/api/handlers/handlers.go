package handlers

import (
	"game/internal/models"
	"game/internal/usecase"
	"log"
	"log/slog"
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
		session.Set("welcome", "true")
		session.Save()
		g.HTML(http.StatusOK, "welcome.html", gin.H{})
	}
}

func RoleHandler(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(g *gin.Context) {
		session := sessions.Default(g)
		session.Set("role", "yes")
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

		log.Println("Username:", username, "Password:", password)

		if username == "admin" && password == "123" {
			session := sessions.Default(c)
			session.Set("username", username)
			session.Save()

			c.JSON(http.StatusOK, gin.H{
				"message":  "login success",
				"redirect": "/home/role/admin-panel",
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

func AdminWebSocketHandler(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {

		logged, err := userOperator.IsAdminLoggedIn()
		if err != nil {
			log.Println("Failed to check the admin:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if logged {
			log.Println("Admin is already logged in. Redirecting...")
			c.Redirect(http.StatusFound, "/home/role/login")
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}

		admin := models.NewAdmin("Anton Fedorov")

		if err := userOperator.AddAdmin(admin); err != nil {
			log.Printf("Failed to add an admin: %v", err)
			return
		}

		go admin.Run(conn)
	}
}

func ClientWebSocketHandler(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade WebSocket"})
			return
		}

		defer func() {
			conn.Close()
			log.Println("Connection closed!")
		}()

		exceeded, err := userOperator.PlayersNumberExceeded()
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("internal server error"))
			log.Println("Failed to check the numbers of players:", err)
			return
		}

		logged, err := userOperator.IsAdminLoggedIn()
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("internal server error"))
			log.Println("Failed to check the admin on redis:", err)
			return
		}

		if !logged {
			conn.WriteMessage(websocket.TextMessage, []byte("admin not logged yet"))
			log.Println("Admin is not logged in yet")
			return
		}

		if exceeded {
			conn.WriteMessage(websocket.TextMessage, []byte("players exceeded"))
			log.Println("9 players are already there")
			return
		}

		player := models.NewPlayer(userOperator.CountPlayers() + 1, "Name", slog.Default())

		if err := userOperator.AddPlayer(player); err != nil {
			log.Printf("Failed to add a player: %v", err)
			return
		}
		select {
		case <-player.Accepted:
			log.Println("player accepted")
			conn.WriteMessage(websocket.TextMessage, []byte("player accepted"))
			return
		case <-player.Rejected:
			log.Println("player rejected")
			conn.WriteMessage(websocket.TextMessage, []byte("player rejected"))
			return
		}

	}
}



func LogoutHandler(userOperator usecase.UseCase, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/home/role/login")
	}
}
