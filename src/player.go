package src

import "fmt"

const ROLECOUNT = 4

type Role int

const (
	CITIZEN Role = iota
	MAFIA
	DOCTOR
	SHERIFF
)

type Player struct {
	RoomOwner    bool
	Room         *Room
	Number       uint
	Name         string
	Job          Role
	Invulnerable bool
	Votes        uint
}

func (pl Player) String() string {
	return fmt.Sprintf("Name: %v and Job: %v", pl.Name, pl.Job)
}

func (pl *Player) Save() {
	pl.Invulnerable = true
}

func (pl *Player) Kill() {

}

func (pl *Player) CreateRoom(roomName string) *Room {
	return CreateRoom(roomName, pl)
}

func (pl *Player) Blame(plName string) {

}

func (pl *Player) StartGame() {
	if pl.RoomOwner == true && len(pl.Room.players) >= 6 {
		pl.Room.StartGame()
	}
}
