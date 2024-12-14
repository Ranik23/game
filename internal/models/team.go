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