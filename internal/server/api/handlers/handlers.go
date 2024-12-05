package handlers

import (
	//"game/internal/models"
	"game/internal/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WelcomeHandler(use usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(200, "welcome.html", gin.H{})
	}
}

func RoleHandler(use usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusAccepted, "auth.html", gin.H{})
	}
}

func LoginHandlerGET(use usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusOK, "login.html", gin.H{})
	}
}

func LoginHandlerPOST(use usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
	
		username := c.PostForm("username")
		password := c.PostForm("password")

		log.Println("Username:", username, "Password:", password)

		if username == "admin" && password == "123" {
			c.JSON(http.StatusOK, gin.H{
				"message": "login success",
				"redirect": "/home/role/admin-panel",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid_credentials",
			})
		}
	}
}



func MainHandler(use usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(http.StatusAccepted, "main.html", gin.H{})
	}
}

func AdminMainHandler(usecase.UseCase) gin.HandlerFunc {
	return func (g *gin.Context)  {
		g.HTML(http.StatusAccepted, "admin_main.html", gin.H{})
	}
}


func AuthPostHandler(use usecase.UseCase) gin.HandlerFunc {
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



// handler для обработки соединения между админом и сервером
func WebSocketHandlerMain(use usecase.UseCase) gin.HandlerFunc {
	return func (c *gin.Context)  {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}		
		defer conn.Close()

		//admin := &models.Admin{}

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
}

// handler для обработки соединения между клиентом и сервером
func WebSocketHandler(use usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}
		defer conn.Close()

		// вот тут должно идти разделение на сессии

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
}
