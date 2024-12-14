package models

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

type Admin struct {
	UserName string
	Teams    []*Team
}

type PlayersInfo struct {
	Action  string   `json:"action"`
	Content []Player `json:"content"`
}

func NewAdmin(name string) *Admin {
	return &Admin{
		UserName: name,
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

		action, ok := input["Action"].(string)
		if !ok {
			log.Printf("Invalid action: %v", input["Action"])
			continue
		}

		switch action {
		case "get_players":
			if err := a.getPlayers(connection); err != nil {
				log.Printf("Failed to get players: %v", err)
			}
		case "accept_player":
			id, ok := input["PlayerID"].(float64) // Исправлено на float64
			if !ok {
				log.Printf("Invalid player ID: %v", input["PlayerID"])
				continue
			}
			if err := a.acceptPlayer(connection, int(id)); err != nil {
				log.Printf("Failed to accept player: %v", err)
			}
		case "reject_player":
			id, ok := input["PlayerID"].(float64) // Исправлено на float64
			if !ok {
				log.Printf("Invalid player ID: %v", input["PlayerID"])
				continue
			}
			if err := a.rejectPlayer(connection, int(id)); err != nil {
				log.Printf("Failed to reject player: %v", err)
			}
		case "start_game":
			if err := a.startGame(connection); err != nil {
				log.Printf("Failed to start game: %v", err)
			}
		default:
			log.Printf("Unknown action: %v", action)
		}
	}
}

func (a *Admin) getPlayers(connection *websocket.Conn) error {
	var newPlayers []Player

	for _, team := range a.Teams {
		for _, player := range team.Players {
			newPlayers = append(newPlayers, Player{
				ID:       player.ID,
				UserName: player.UserName,
			})
		}
	}

	message := map[string]interface{}{
		"Action":  "players_list",
		"Players": newPlayers,
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
	for _, team := range a.Teams {
		for _, player := range team.Players {
			player.GameStarted <- struct{}{}
		}
	}

	message := map[string]interface{}{
		"Action":  "game_started",
		"Message": "The game has started!",
	}
	return writeJSONWithLog(connection, message)
}

func (a *Admin) handlePlayer(connection *websocket.Conn, id int, action string) error {
	for _, team := range a.Teams {
		for index, player := range team.Players {
			if player.ID == id {
				switch action {
				case "accept":
					player.Accepted <- struct{}{}
					message := map[string]interface{}{
						"Action":   "player_accepted",
						"UserName": player.UserName,
						"ID":       strconv.Itoa(player.ID),
					}
					if err := writeJSONWithLog(connection, message); err != nil {
						return err
					}
					return nil

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
					team.Players = append(team.Players[:index], team.Players[index+1:]...)
					return nil
				}
			}
		}
	}
	return ErrPlayerNotFound
}