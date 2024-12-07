package models

import (
	"log"
	"github.com/gorilla/websocket"
)




type Player struct {
	ID   int
	Name string
	Projects []*Project
}


func NewPlayer(id int, name string) *Player {
	return &Player{
		ID: id,
		Name: name,
	}
}

// TODO: унифицировать формат сообщения
type Message struct {
	Action string	`json:"action"`
	Data   string	`json:"data"`
}


func (p *Player) Run(connection *websocket.Conn) error {
	for {
		var message Message
		
		if err := connection.ReadJSON(&message); err != nil {
			log.Println("Failed to read JSON message")
			return err
		}

		action := message.Action

		switch action {



		default:
			log.Println("Action Not Supported")
			return ErrUnsupportedAction
		}
	}
}