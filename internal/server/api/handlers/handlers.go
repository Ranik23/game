package handlers

import (
	"game/internal/models"
	"game/internal/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Проверяем источник запроса (в реальном приложении настройте это)
		return true
	},
}

func WelcomeHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "welcome.html", gin.H{})
	}
}

func RoleHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "auth.html", gin.H{})
	}
}

func LoginHandlerGET(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "login.html", gin.H{})
	}
}

func LoginHandlerPOST(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		log.Println("Username:", username, "Password:", password)

		if username == "admin" && password == "123" {
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

func MainHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "player-main.html", gin.H{})
	}
}

func AdminMainHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "admin-main.html", gin.H{})
	}
}

func AuthPostHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		username := g.PostForm("username")
		password := g.PostForm("password")

		if username == "admin" && password == "admin123" {
			g.Redirect(http.StatusFound, "/admin_main")
		} else {
			g.HTML(http.StatusOK, "auth.html", gin.H{
				"error": "Invalid username or password",
			})
		}
	}
}

// Handler для обработки соединения между админом и сервером
func AdminWebSocketHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}

		// TODO: проверить логики, уточнить в случае неудачи, куда мы перенаправялем пользователя
		// logged, err := userOperator.IsAdminLoggedIn()
		// if err != nil {
		// 	log.Println("failed to check the admin on redis")
		// 	return
		// }

		// if logged {
		// 	c.Redirect(http.StatusFound, "/home/role/login")
		// 	return // возможно стоит переправить на другую  html страницу
		// }

		admin := models.NewAdmin("")

		if err := userOperator.AddAdmin(admin); err != nil {
			log.Printf("Failed to add the admin: %v", err)
			return
		}

		go admin.Run(conn)
	}
}

// Handler для обработки соединения между клиентом и сервером
func ClientWebSocketHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}

		exceeded, err := userOperator.PlayersNumberExceeded()
		if err != nil {
			log.Println("Failed to check the numbers of players")
			return
		}

		if exceeded {
			log.Println("9 players are already there")
			return
		}

		player := models.NewPlayer(userOperator.CountPlayers() + 1, "")

		if err := userOperator.AddPlayer(player); err != nil {
			log.Println("Failed to add a player")
			return
		}

		// TODO: тут еще нужно генерить ПРОЕКТНЫЕ ПРЕДЛОЖЕНИЯ, которые будут генериться через его алгоритм
		// и передавать мы их будем в конструктор

 		// TODO: как добавять имя? через запрос?

		// TODO : проверить логику этого метода, точнее дописать ее
		go player.Run(conn)
	}
}
