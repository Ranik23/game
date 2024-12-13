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
	Action  string   `json:"action"`
	Content []Player `json:"content"`
}

var (
	ErrPlayerNotFound = fmt.Errorf("player not found")
)

func NewAdmin(name string) *Admin {
	return &Admin{
		Name: name,
	}
}

func writeJSONWithLog(conn *websocket.Conn, message map[string]interface{}) error {
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

		var input map[string]interface{}

		if err := json.Unmarshal(message, &input); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		if action, ok := input["Action"].(string); !ok {
			switch action {
			case "get_players":
				return a.getPlayers(connection)
			case "accept_player":
				id, ok := input["PlayerID"].(int)
				if !ok {
					log.Printf("Invalid player ID: %v", id)
					continue
				}
				return a.acceptPlayer(connection, id)
			case "reject_player":
				id, ok := input["PlayerID"].(int)
				if !ok {
					log.Printf("Invalid player ID: %v", id)
					continue
				}
				return a.rejectPlayer(connection, id)
			case "start_game":
				return a.startGame(connection)
			}
		}
	}
}

func (a *Admin) getPlayers(connection *websocket.Conn) error {
	var new_players []Player

	var message map[string]interface{}

	for _, player := range a.Players {
		new_players = append(new_players, Player{
			ID:       player.ID,
			UserName: player.UserName,
		})
	}
	message = map[string]interface{}{
		"Action": "players_list",
		"Players": new_players,
	}
	return writeJSONWithLog(connection, message)
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
	message := map[string]interface{}{
		"Action":  "game_started",
		"Message": "The game has started!",
	}
	return writeJSONWithLog(connection, message)
}

func (a *Admin) handlePlayer(connection *websocket.Conn, id int, action string) error {
	for index, player := range a.Players {
		if player.ID == id {
			switch action {
			case "accept":
				player.Accepted <- struct{}{}
				message := map[string]interface{}{
					"Action":   "player_accepted",
					"UserName": player.UserName,
					"ID":       strconv.Itoa(player.ID),
				}
				a.Players = append(a.Players, player)
				return writeJSONWithLog(connection, message)

			case "reject":
				player.Rejected <- struct{}{}
				message := map[string]interface{}{
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
