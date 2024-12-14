package handlers

import (
	"game/internal/models"
	"game/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)





const (
	MaxTeams = 9
)




func LeaderWebSocketHanlder(userOperator usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}

		count, err := userOperator.CountTeams()
		if err != nil {
			log.Println("Failed to count teams:", err)
			return
		}

		if count >= MaxTeams {
			log.Println("Max teams reached")
			return
		}

		team, err := userOperator.CreateTeam()
		if err != nil {
			log.Println("Failed to create team:", err)
			return
		}

		leader := models.NewLeader("", team)

		go leader.Run(conn)
	}
}