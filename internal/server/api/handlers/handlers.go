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

func HomeHandler(use usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.JSON(200, gin.H{"Message": "Hello, World!"})
	}
}

func AuthHandler(use usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		role := g.Query("role")
		if role == "guest" {
			g.Redirect(http.StatusFound, "/main")
		} else {
			g.HTML(http.StatusOK, "auth.html", gin.H{})
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

func WelcomeHandler(use usecase.UseCase) gin.HandlerFunc {
	return func(g *gin.Context) {
		g.HTML(200, "welcome.html", gin.H{})
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
