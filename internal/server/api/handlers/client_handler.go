package handlers

import (
	"errors"
	"game/internal/models"
	"game/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func ClientWebSocketHandler(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}

		exceeded, err := userOperator.PlayersNumberExceeded()
		if err != nil {
			if errors.Is(err, usecase.ErrNoAdminSet) {
				sendMessage(conn, "error", "Admin is not logged in yet")
				log.Println("Admin is not logged in yet")
			} else {
				sendMessage(conn, "error", "Failed to check the numbers of players")
				log.Printf("Failed to check the numbers of players: %v", err)
			}
			return
		}

		logged, err := userOperator.IsAdminLoggedIn()
		if err != nil {
			sendMessage(conn, "error", "Failed to check the admin")
			log.Printf("Failed to check the admin: %v", err)
			return
		}

		if !logged {
			sendMessage(conn, "error", "Admin is not logged in yet")
			log.Println("Admin is not logged in yet")
			return
		}

		if exceeded {
			sendMessage(conn, "error", "Players Limit Exceeded")
			log.Println("Players Limit Exceeded")
			return
		}

		count, err := userOperator.CountPlayers()
		if err != nil {
			sendMessage(conn, "error", "Failed to count the players")
			log.Printf("Failed to count the players: %v", err)
			return
		}

		player := models.NewPlayer(count + 1, "Anton Fedorov")

		if err := userOperator.AddPlayer(nil, player); err != nil {
			sendMessage(conn, "error", "Failed to add a player")
			log.Printf("Failed to add a player: %v", err)
			return
		}
		select {
		case <-player.Accepted:
			log.Println("player accepted")
			sendMessage(conn, "message", "Player accepted")
			go player.Run(conn)
		case <-player.Rejected:
			log.Println("player rejected")
			sendMessage(conn, "message", "Player Rejected")
			conn.Close()
		}
	}
}
