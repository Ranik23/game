package models


type Team struct {
	ID		int
	Name 	string
	Leader 	*Player
	Players []*Player
}

func NewTeam(Leader *Player) *Team {
	return &Team{
		Leader: Leader,
	}
}

func (t *Team) AddPlayer(player *Player) {
	t.Players = append(t.Players, player)
}