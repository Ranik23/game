package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"github.com/gorilla/websocket"
)

type Admin struct {
	Name    string
	Players []*Player
}

type PlayersInfo struct {
	Action  string    `json:"action"`
	Content []*Player `json:"content"`
}

var (
	ErrPlayerNotFound = fmt.Errorf("player not found")
)

func NewAdmin(name string) *Admin {
	return &Admin{
		Name: name,
	}
}

func writeJSONWithLog(conn *websocket.Conn, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return err
	}
	log.Printf("Sending message: %s", data)
	return conn.WriteJSON(message)
}

func (a *Admin) Run(connection *websocket.Conn) error {
	defer connection.Close()

	actionHandlers := map[string]func(a *Admin, conn *websocket.Conn, data string) error{
		"get_players": func(a *Admin, conn *websocket.Conn, data string) error {
			return a.getPlayers(conn)
		},
		"accept_player": func(a *Admin, conn *websocket.Conn, data string) error {
			id, err := strconv.Atoi(data)
			if err != nil {
				return fmt.Errorf("invalid player ID: %v", err)
			}
			return a.acceptPlayer(conn, id)
		},
		"delete_player": func(a *Admin, conn *websocket.Conn, data string) error {
			id, err := strconv.Atoi(data)
			if err != nil {
				return fmt.Errorf("invalid player ID: %v", err)
			}
			return a.rejectPlayer(conn, id)
		},
		"start_game": func(a *Admin, conn *websocket.Conn, data string) error {
			return a.startGame(conn)
		},
	}

	for {

		t, message, err := connection.ReadMessage()
		if err != nil {
			log.Printf("Failed to read the message: %v", err)
			return err
		}

		if t != websocket.TextMessage {
			log.Printf("Unsupported message type: %v", t)
			continue
		}

		var input struct {
			Action string `json:"action"`
			Data   string `json:"data"`
		}
		if err := json.Unmarshal(message, &input); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		if handler, ok := actionHandlers[input.Action]; ok {
			if err := handler(a, connection, input.Data); err != nil {
				log.Printf("Error handling action %s: %v", input.Action, err)
			}
		} else {
			log.Printf("Unsupported action: %s", input.Action)
			connection.WriteMessage(websocket.TextMessage, []byte("error: unsupported action"))
		}
	}
}

func (a *Admin) getPlayers(connection *websocket.Conn) error {
	return writeJSONWithLog(connection, PlayersInfo{Action: "players_list", Content: a.Players})
}

func (a *Admin) handlePlayer(connection *websocket.Conn, id int, action string) error {
	for index, player := range a.Players {
		if player.ID == id {
			switch action {
			case "accept":
				player.Accepted <- struct{}{}
				message := map[string]string{
					"Action":   "player_accepted",
					"UserName": player.UserName,
					"ID":       strconv.Itoa(player.ID),
				}
				return writeJSONWithLog(connection, message)

			case "reject":
				player.Rejected <- struct{}{}
				message := map[string]string{
					"Action":   "player_rejected",
					"UserName": player.UserName,
					"ID":       strconv.Itoa(player.ID),
				}
				if err := writeJSONWithLog(connection, message); err != nil {
					return err
				}
				a.Players = append(a.Players[:index], a.Players[index+1:]...)
				return nil
			}
		}
	}
	return ErrPlayerNotFound
}

func (a *Admin) acceptPlayer(connection *websocket.Conn, id int) error {
	return a.handlePlayer(connection, id, "accept")
}

func (a *Admin) rejectPlayer(connection *websocket.Conn, id int) error {
	return a.handlePlayer(connection, id, "reject")
}

func (a *Admin) startGame(connection *websocket.Conn) error {
	for _, player := range a.Players {
		player.GameStarted <- struct{}{}
	}

	message := map[string]string{
		"type":    "game_started",
		"message": "The game has started!",
	}
	return writeJSONWithLog(connection, message)
}
