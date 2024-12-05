package session

import (
	"game/internal/models"
	"time"
)


type PlayerSession struct {
	Player 			*models.Player
	CreatedAt 		time.Time
	FinishedAt 		time.Time
	Projects		[]models.Project
}

type GameSession struct {
	Admin 		*models.Admin
	Players 	[]*models.Player
	CreatedAt 	time.Time
	FinishedAt 	time.Time	
}

func NewGameSession() *GameSession {
	return &GameSession{}
}