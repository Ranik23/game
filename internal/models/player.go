package models

import (
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID       	int           	`json:"id"`
	UserName 	string        	`json:"username"`
	Projects 	[]string      	`json:"projects"`
	Accepted 	chan struct{} 	`json:"-"`
	Rejected 	chan struct{} 	`json:"-"`
	GameStarted	chan struct{}	`json:"-"`
	Logger   	*slog.Logger  	`json:"-"`
}

func NewPlayer(id int, name string) *Player {
	return &Player{
		ID:       id,
		UserName: name,
		Accepted: make(chan struct{}),
	}
}

type Message struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

type ConnectRequest struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
}

var (
	ErrRejectedByAdmin = fmt.Errorf("a player has been rejected by admin")
)

func (p *Player) Run(connection *websocket.Conn) error {
	defer connection.Close()



	for {
		
	}
	return nil
}