package session

import (
	"time"
)


type PlayerSession struct {
	ID 		int		`gorm:"primaryKey"`
	Name 	string	`gorm:""`
	Role 	string	`gorm:""`
	CreatedAt 		time.Time
	FinishedAt 		time.Time
}

type GameSession struct {
	ID 				int					`gorm:"primaryKey"`
	Players 		[]*PlayerSession	`gorm:""`
	Finished 		bool				`gorm:""`
	CreatedAt 		time.Time
	FinishedAt 		time.Time
	MaxPlayers  	int
	CurrentPlayers	int
}
