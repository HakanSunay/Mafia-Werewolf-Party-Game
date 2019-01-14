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
	Dead         bool
}

func (pl Player) String() string {
	return fmt.Sprintf("Name: %v and Job: %v", pl.Name, pl.Job)
}

// This method is invoked when a player of type DOCTOR
// decides to use his special ability on us
func (pl *Player) Save() {
	pl.Invulnerable = true
}

// Kill the current player
func (pl *Player) Die() {
	pl.Dead = true
}

func (pl *Player) CreateRoom(roomName string) *Room {
	return CreateRoom(roomName, pl)
}

// Blames another player using their name as a parameter
// aka votes to send him to prison
func (pl *Player) Blame(blamee *Player) {
	blamee.Votes++
}

// This method can be invoked only if
// the current object is the owner of the room
func (pl *Player) StartGame() {
	if pl.RoomOwner == true && len(pl.Room.players) >= 6 {
		pl.Room.StartGame()
	}
}

func (pl *Player) ResetRound() {
	pl.Invulnerable = false
	pl.Votes = 0
}
