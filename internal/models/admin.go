package models

import (
	"log"

	"github.com/gorilla/websocket"
)



type Admin struct {
	Name string
	Players []*Player
}


func NewAdmin(name string) *Admin {
	return &Admin{
		Name: name,
	}
}


func (a *Admin) Run(connection *websocket.Conn) error {
	for {

		var message Message

		if err := connection.ReadJSON(&message); err != nil {
			log.Printf("Failed to read JSON message: %v", err)
			return err
		}

		switch message.Action {
		case "get_players":
			if err := a.getPlayers(connection); err != nil { // TODO : либо можем не передавать соединение, а просто получать список игроков, а потом уже отсюад отправлять 
				log.Printf("Failed to get players: %v", err)
				return err
			}
		
		default:
			log.Println("Action Not Supported")
			return ErrUnsupportedAction
		}
	}
}


type DataResponse struct {
	Type 	string `json:"type"`
	Content interface{} `json:"content"`
}

func (a *Admin) getPlayers(connection *websocket.Conn) error {
	// TODO: Здесь мы должны получить список игроков из вашей базы данных или другого источника
	// в этом примере я просто создал фиктивный список игроков
	players := []*Player{
		{ID: 1, Name: "Игрок 1", Projects: nil},
		{ID: 2, Name: "Игрок 2", Projects: nil},
	}

	data := DataResponse{
		Type:    "players_list",
		Content: players,
	}

	if err := connection.WriteJSON(data); err != nil {
		log.Println("Failed to send JSON message:", err)
		return err
	}

	return nil
} 