package handlers

import (
	"game/internal/models"
	"game/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func AdminWebSocketHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}

		logged, err := userOperator.IsAdminLoggedIn()
		if err != nil {
			log.Println("Failed to check the admin:", err)
			sendMessage(conn, "error", "Failed to check the admin")
			return
		}

		if logged {
			log.Println("Admin is already logged in. Redirecting...")
			sendMessage(conn, "redirect", "/role")
			return
		}

		admin := models.NewAdmin("Anton Fedorov")

		if err := userOperator.AddAdmin(admin); err != nil {
			log.Printf("Failed to add an admin: %v", err)
			sendMessage(conn, "error", "Failed to add an admin")
			return
		}

		go admin.Run(conn)
	}
}