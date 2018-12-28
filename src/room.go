package src

import (
	"math/rand"
	"time"
)

type Room struct {
	name    string
	players []Player
}

func CreateRoom(name string, player *Player) *Room {
	resultRoom := &Room{name, nil}
	newPlayers := make([]Player, 1)
	newPlayers = append(newPlayers, *player)
	resultRoom.players = newPlayers
	return resultRoom
}

func (r *Room) AddPlayer(player *Player){
	r.players = append(r.players, *player)
}
var (
	citizenCount uint = 0
	mafiaCount   uint = 0
	doctorCount  uint = 0
	sheriffCount uint = 0
	ClientCount  uint = 0
)

func RandomJob(curRoles *map[Role]uint) Role {
	rand.Seed(time.Now().UnixNano())
	if ClientCount < 4 {
		choice := rand.Intn(ROLECOUNT)
		return Role(choice)
	} else if doctorCount == 0 {
		doctorCount += 1
		return DOCTOR
	} else if sheriffCount == 0 {
		sheriffCount += 1
		return SHERIFF
	} else if citizenCount < mafiaCount {
		citizenCount += 1
		return CITIZEN
	} else {
		mafiaCount += 1
		return MAFIA
	}
}

func (r *Room) StartGame(){

}