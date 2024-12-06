package handlers

import (
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
		g.HTML(http.StatusOK, "main.html", gin.H{})
	}
}

func AdminMainHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "admin_main.html", gin.H{})
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
		// TODO: добавить проверку на существование админа(админ только один)
		go handleWebSocketConnection(userOperator, conn)
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
		// TODO : добавить проверку на уже существующее количество пользователей(пользователей не должно быть более 9)
		go handleWebSocketConnection(userOperator, conn)
	}
}

// Общая функция для обработки WebSocket-соединения
func handleWebSocketConnection(userOperator usecase.UseCase, conn *websocket.Conn) {
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			return
		}

		log.Printf("Received message: %s", message)

		err = conn.WriteMessage(messageType, []byte("Message received"))
		if err != nil {
			log.Println("Failed to write message:", err)
			break
		}
	}
}