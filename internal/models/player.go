package models

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"github.com/gorilla/websocket"
)

var (
	ErrRejectedByAdmin = fmt.Errorf("a player has been rejected by admin")
)

type Player struct {
	ID       	int           	`json:"id"`
	UserName 	string        	`json:"username"`
	Role		string			`json:"role"`
	Accepted 	chan struct{} 	`json:"-"`
	Rejected 	chan struct{} 	`json:"-"`
	GameStarted	chan struct{}	`json:"-"`
	Logger   	*slog.Logger  	`json:"-"`
}

func NewPlayer(ID int, Name string) *Player {
	return &Player{
		ID:       ID,
		UserName: Name,
		Accepted: make(chan struct{}),
		Rejected: make(chan struct{}),
		GameStarted: make(chan struct{}),
		Logger: slog.Default(),
	}
}

func (p *Player) Run(connection *websocket.Conn) error {
	defer connection.Close()
	for {

		t, message, err := connection.ReadMessage()

		if err != nil {
			return err
		}
		
		if t != websocket.TextMessage {
			log.Printf("Unsupported message type: %v", t)
			continue
		}

		var input map[string]interface{}

		if err := json.Unmarshal(message, &input); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		if action, ok := input["Action"].(string); ok {
			switch action {
			}
		}
	}
}
