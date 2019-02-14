package src

import (
	"fmt"
	"math/rand"
	"time"
)

type Room struct {
	name    string
	players []*Player
	playing bool
	stage   int
}

const (
	MAFIASTAGE int = iota
	DOCTORSTAGE
	ALLSTAGE
)

func (r Room) String() string {
	return fmt.Sprintf("The name of the room is %v and the players are %v", r.name, r.players)
}

func CreateRoom(name string, player *Player) *Room {
	resultRoom := &Room{name, nil, false, 0}
	newPlayers := make([]*Player, 0)
	newPlayers = append(newPlayers, player)
	resultRoom.players = newPlayers
	return resultRoom
}

func (r *Room) AddPlayer(player *Player) {
	if len(r.players) == 0 {
		player.RoomOwner = true
	}
	player.Room = r
	r.players = append(r.players, player)
}

var (
	citizenCount uint = 0
	mafiaCount   uint = 0
	doctorCount  uint = 0
	ClientCount  uint = 0
)

func (r *Room) StartGame() {
	rand.Seed(time.Now().UnixNano())
	amountOfPlayers := len(r.players)
	doctorChoice := rand.Intn(amountOfPlayers)
	r.players[doctorChoice].Job = DOCTOR
	for index, _ := range r.players {
		if index%2 == 0 && r.players[index].Job != DOCTOR {
			r.players[index].Job = MAFIA
		} else if r.players[index].Job != DOCTOR {
			r.players[index].Job = CITIZEN
		}
	}
	r.playing = true
	r.stage = 0
}

func (r *Room) IsPlaying() bool {
	return r.playing
}

// Finds the sought room, when JOIN room is called
func FindRoom(rooms *[]Room, name string) *Room {
	for _, rm := range *rooms {
		if rm.name == name {
			return &rm
		}
	}
	return nil
}

// Searches for a player using their name in the current room
func (r *Room) FindPlayer(name string) *Player {
	for _, rm := range r.players {
		if rm.Name == name {
			return rm
		}
	}
	return nil
}

// Reset changes after round completion
func (r *Room) Reset() {
	for _, player := range r.players {
		player.ResetRound()
	}
}

func (r *Room) GetOwner() *Player {
	for _, pl := range r.players {
		if pl.RoomOwner == true {
			return pl
		}
	}
	return nil
}

func (r *Room) GetPlayers() []*Player {
	return r.players
}

func (r *Room) SetName(name string) {
	r.name = name
}

func (r *Room) GetName() string {
	return r.name
}

func (r *Room) GetStage() int {
	return r.stage
}

func (r *Room) CanGoToNextStage() bool {
	if (r.stage == MAFIASTAGE && r.CheckIfMafiaVoted()) ||
		(r.stage == ALLSTAGE && r.CheckIfAllVoted()) ||
		(r.stage == DOCTORSTAGE && (r.CheckIfDoctorSaved() || !(r.HasDoctor()))){
		return true
	}
	return false
}

func (r *Room) NextStage() {
	if r.CanGoToNextStage() {
		r.stage++
		r.stage %= 3
	}
}

func (r *Room) GetMostVotedPlayer() *Player {
	maxVotedPlayer := r.players[0]
	for _, pl := range r.players {
		if pl.Votes > maxVotedPlayer.Votes {
			maxVotedPlayer = pl
		}
	}
	return maxVotedPlayer
}

func (r *Room) CheckIfAllVoted() bool {
	for _, pl := range r.players {
		if r.stage == ALLSTAGE && pl.Voted == false && pl.Dead == false {
			return false
		}
	}
	return true
}

func (r *Room) CheckIfMafiaVoted() bool {
	for _, pl := range r.players {
		if pl.Dead == false && pl.Job == MAFIA && r.stage == MAFIASTAGE && pl.Voted == false {
			return false
		}
	}
	return true
}

func (r *Room) CheckIfDoctorSaved() bool {
	for _, pl := range r.players {
		if pl.Job == DOCTOR && r.stage == DOCTORSTAGE && pl.Voted == false {
			return false
		}
	}
	return true
}

func (r *Room) FindChosenPlayerToDie() *Player {
	for index, _ := range r.players {
		if r.players[index].Dead == false && r.players[index].Chosen == true {
			r.players[index].Dead = true
			return r.players[index]
		}
	}
	return nil
}

func (r *Room) HasDoctor() bool {
	for _, pl := range r.players {
		if pl.Job == DOCTOR && pl.Dead == false {
			return true
		}
	}
	return false
}

func (r* Room) GameOver() (bool, Role) {
	aliveMafia := 0
	aliveCitizensDocs := 0
	for _, pl := range r.players{
		if pl.Job == MAFIA && pl.Dead == false{
			aliveMafia++
		} else if (pl.Job == CITIZEN || pl.Job == DOCTOR) && pl.Dead == false{
			aliveCitizensDocs++
		}
	}
	if aliveMafia == 0 {
		return true, CITIZEN
	} else if aliveCitizensDocs <= 1 && aliveMafia >= 1{
		return true, MAFIA
	} else {
		return false, 0
	}
}
func (r *Room) End() {
	for index, _ := range r.players{
		r.players[index].End()
	}
	r.playing = false
	r = nil
}