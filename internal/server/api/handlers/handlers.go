package handlers

import (
	"encoding/json"
	"game/internal/models"
	"game/internal/usecase"
	"log"
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


func sendErrorMessage(conn *websocket.Conn, action, message string) {
    errMsg := map[string]string{
        "action":  action,
        "message": message,
    }
    msgBytes, err := json.Marshal(errMsg)
    if err != nil {
        log.Println("Failed to marshal error message:", err)
        return
    }

    if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
        log.Println("Failed to send error message:", err)
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

		exceeded, err := userOperator.PlayersNumberExceeded()
		if err != nil {
			sendErrorMessage(conn, "internal_error", "Failed to check the numbers of players")
			log.Println("Failed to check the numbers of players:", err)
			return
		}

		logged, err := userOperator.IsAdminLoggedIn()
		if err != nil {
			sendErrorMessage(conn, "internal_error", "Failed to check the admin on redis")
			log.Println("Failed to check the admin on redis:", err)
			return
		}

		if !logged {
			sendErrorMessage(conn, "admin_not_logged", "Admin is not logged yet")
			log.Println("Admin is not logged in yet")
			return
		}

		if exceeded {
			sendErrorMessage(conn, "players_exceeded", "Players Limit Exceeded")
			log.Println("Players Limit Exceeded")
			return
		}

		player := models.NewPlayer(userOperator.CountPlayers() + 1, "Name")

		if err := userOperator.AddPlayer(player); err != nil {
			log.Printf("Failed to add a player: %v", err)
			return
		}
		select {
		case <-player.Accepted:
			log.Println("player accepted")
			conn.WriteMessage(websocket.TextMessage, []byte("player accepted"))
			go player.Run(conn)
			return
		case <-player.Rejected:
			log.Println("player rejected")
			conn.WriteMessage(websocket.TextMessage, []byte("player rejected"))
			conn.Close()
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
