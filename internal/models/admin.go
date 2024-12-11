package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)



type Admin struct {
	Name 		string
	Players 	[]*Player
}

type PlayersInfo struct {
	Action  string		`json:"action"`
	Content []*Player 	`json:"content"`
}

type PlayerInfo struct {
	Action   string		`json:"action"`
	UserName string		`json:"username"`
	ID		 int		`json:"id"`
}

func NewAdmin(name string) *Admin {
	return &Admin{
		Name: name,
	}
}

var (
	ErrPlayerNotFound = fmt.Errorf("player not found")
)


func (a *Admin) Run(connection *websocket.Conn) error {
	for {
		var message Message

		// if err := connection.ReadJSON(&message); err != nil {
		// 	log.Printf("Failed to read JSON message: %v", err)
		// 	continue
		// }

		switch message.Action {
		case "get_players":
			if err := a.getPlayers(connection); err != nil { // TODO : либо можем не передавать соединение, а просто получать список игроков, а потом уже отсюад отправлять 
				log.Printf("Failed to get players: %v", err)
				return err
			}
		case "accept_player":
			id, err := strconv.Atoi(message.Data)
			if err != nil {
				log.Printf("Failed to parse player ID: %v", err)
				return err
			}

			if err := a.acceptPlayer(connection, id); err != nil {
				log.Printf("Failed to accept the player: %v", err)
				return err
			}
		case "delete_player":
			id, err := strconv.Atoi(message.Data)
			if err != nil {
				log.Printf("Failed to parse player ID: %v", err)
				return err
			}
			if err := a.rejectPlayer(connection, id); err != nil {
				log.Printf("Failed to delete the player: %v", err)
				return err
			}
		case "start_game":
			a.startGame(connection)
		default:
			log.Println("Action Not Supported")
			return ErrUnsupportedAction
		}
	}
}

func (a *Admin) acceptPlayer(connection *websocket.Conn, id int) error {
	for _, player := range a.Players {
		if player.ID == id {
			player.Accepted <- struct{}{}

			message := map[string]string{
				"Action" : "player_accepted",
				"UserName" : player.UserName,
				"ID" : strconv.Itoa(player.ID),
			}

			if err := connection.WriteJSON(message); err != nil {
				log.Printf("Failed to send JSON message: %v", err)
				return err
			}
			return nil
		}
	}
	return ErrPlayerNotFound
}


func (a *Admin) rejectPlayer(connection *websocket.Conn, id int) error {
	for index, player := range a.Players {
		if player.ID == id {
			player.Rejected <- struct{}{}

			message := map[string]string{
				"Action" : "player_rejected",
				"UserName" : player.UserName,
				"ID" : strconv.Itoa(player.ID),
			}

			if err := connection.WriteJSON(message); err != nil {
				log.Printf("Failed to send JSON message: %v", err)
				return err
			}
			a.Players = append(a.Players[:index], a.Players[index + 1:]...)
			return nil
		}
	}
	return ErrPlayerNotFound
}

func (a *Admin) getPlayers(connection *websocket.Conn) error {
	if err := connection.WriteJSON(PlayersInfo{Action: "players_list", Content: a.Players}); err != nil {
		log.Printf("Failed to send JSON message: %v", err)
		return err
	}
	return nil
} 


func (a *Admin) startGame(connection *websocket.Conn) error {
	for _, player := range a.Players {
		player.GameStarted <- struct{}{}
	}

	message := map[string]string{
		"type":    "game_started",
		"message": "The game has started!",
	}

	if err := connection.WriteJSON(message); err != nil {
		log.Printf("Failed to send JSON message: %v", err)
		return err
	}

	return nil
}