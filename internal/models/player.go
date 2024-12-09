package models

import (
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID       int           `json:"id"`
	UserName string        `json:"username"`
	Projects []string      `json:"projects"`
	Accepted chan struct{} `json:"-"`
	Rejected chan struct{} `json:"-"`
	Logger   *slog.Logger  `json:"-"`
}

func NewPlayer(id int, name string, logger *slog.Logger) *Player {
	return &Player{
		ID:       id,
		UserName: name,
		Accepted: make(chan struct{}),
		Logger: logger,
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
	select {
	case <-p.Accepted:
		p.Logger.Info("A player has been accepted", slog.Int("ID", p.ID))
	case <-p.Rejected:
		p.Logger.Info("A player has been rejected", slog.Int("ID", p.ID))
		return ErrRejectedByAdmin
	}

	for {
	}
}