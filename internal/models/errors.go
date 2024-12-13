package models

import "fmt"


var (
	ErrUnsupportedAction = fmt.Errorf("неизвестное действие")
	ErrPlayerNotFound = fmt.Errorf("player not found")
	ErrRejectedByAdmin = fmt.Errorf("a player has been rejected by admin")
)
