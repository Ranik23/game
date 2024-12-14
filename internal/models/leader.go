package models

import (
	"errors"

	"github.com/gorilla/websocket"
)

type Leader struct {
	UserName 	string
	Team 		*Team
	playersCh	chan *Player
}

const (
	MaxPlayers = 6
)

func NewLeader(userName string, team *Team) *Leader {
	return &Leader{
		UserName: userName,
		Team: team,
		playersCh: make(chan *Player, MaxPlayers),
	}
}

func (l *Leader) Run(connection *websocket.Conn) error {
	defer connection.Close()
	for {
		_, _, err := connection.ReadMessage()
		if err != nil {
			return err
		}
	}
}

func (l *Leader) acceptPlayerIntoTeam(connection *websocket.Conn, id int) error {
	player, ok := <- l.playersCh
	if !ok {
		return errors.New("no players in the channel")
	}
	l.Team.AddPlayer(player)
	return nil
}

func (l *Leader) removePlayerFromTeam(connection *websocket.Conn, id int) error {
	return nil
}

func (l *Leader) setRoleForPlayer(connection *websocket.Conn, id int, role string) error {
	return nil
}