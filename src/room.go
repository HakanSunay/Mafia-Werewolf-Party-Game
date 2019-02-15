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

// CreateRoom creates a new room, by using the 2 parameters
// given as arguments. Name becomes the name of the room
// and the player becomes owner of the room.
func CreateRoom(name string, player *Player) *Room {
	resultRoom := &Room{name, nil, false, 0}
	newPlayers := make([]*Player, 0)
	newPlayers = append(newPlayers, player)
	resultRoom.players = newPlayers
	return resultRoom
}

// AddPlayer is used to add new players to the
// current room. The player given as parameter
// becomes the leader of the room if the room is empty.
func (r *Room) AddPlayer(player *Player) {
	if len(r.players) == 0 {
		player.RoomOwner = true
	}
	player.Room = r
	r.players = append(r.players, player)
}

var ClientCount uint = 0

// StartGame is used by owner/leader of the room
// to initiate the game. The players inside the room
// are given roles as DOCTOR/CITIZEN/MAFIA on a random basis.
// The room starts playing and the stage becomes 0 == MAFIASTAGE.
func (r *Room) StartGame() {
	rand.Seed(time.Now().UnixNano())
	amountOfPlayers := len(r.players)
	doctorChoice := rand.Intn(amountOfPlayers)
	r.players[doctorChoice].Job = DOCTOR
	for index := range r.players {
		if index%2 == 0 && r.players[index].Job != DOCTOR {
			r.players[index].Job = MAFIA
		} else if r.players[index].Job != DOCTOR {
			r.players[index].Job = CITIZEN
		}
	}
	r.playing = true
	r.stage = 0
}

// IsPlaying checks if the current room is playing.
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

// GetOwner finds the owner/leader of the room.
// If there is no leader, returns nil.
func (r *Room) GetOwner() *Player {
	for _, pl := range r.players {
		if pl.RoomOwner == true {
			return pl
		}
	}
	return nil
}

// GetPlayers is a getter for the players of the room.
func (r *Room) GetPlayers() []*Player {
	return r.players
}

// SetName is used to change the name of the current room.
func (r *Room) SetName(name string) {
	r.name = name
}

// GetName returns the name of the room.
func (r *Room) GetName() string {
	return r.name
}

// GetStage returns the current stage of the room.
func (r *Room) GetStage() int {
	return r.stage
}

// CanGoToNextStage performs a couple of check-ups
// that are based on the current stage of the room.
// For example, if the room is at the MAFIASTAGE,
// in order to advance to the DOCTORSTAGE, all of
// the Mafia members that are alive should have casted their votes.
func (r *Room) CanGoToNextStage() bool {
	if (r.stage == MAFIASTAGE && r.CheckIfMafiaVoted()) ||
		(r.stage == ALLSTAGE && r.CheckIfAllVoted()) ||
		(r.stage == DOCTORSTAGE && (r.CheckIfDoctorSaved() || !(r.HasDoctor()))) {
		return true
	}
	return false
}

// NextStage changes the current stage of the room.
func (r *Room) NextStage() {
	if r.CanGoToNextStage() {
		r.stage++
		r.stage %= 3
	}
}

// GetMostVotedPlayer loops through the players
// of the room and finds the one, who has the most votes.
func (r *Room) GetMostVotedPlayer() *Player {
	maxVotedPlayer := r.players[0]
	for _, pl := range r.players {
		if pl.Votes > maxVotedPlayer.Votes {
			maxVotedPlayer = pl
		}
	}
	return maxVotedPlayer
}

// CheckIfAllVoted performs a check-up on the players
// to determine if they have all voted. Keep in mind,
// this method is used during the ALLSTAGE.
func (r *Room) CheckIfAllVoted() bool {
	for _, pl := range r.players {
		if r.stage == ALLSTAGE && pl.Voted == false && pl.Dead == false {
			return false
		}
	}
	return true
}

// CheckIfMafiaVoted performs a check-up on the MAFIA members
// to determine if they have all voted. Keep in mind,
// this method is used during the MAFIASTAGE.
func (r *Room) CheckIfMafiaVoted() bool {
	for _, pl := range r.players {
		if pl.Dead == false && pl.Job == MAFIA && r.stage == MAFIASTAGE && pl.Voted == false {
			return false
		}
	}
	return true
}

// CheckIfDoctorSaved performs a check-up
// to determine if they has voted. Keep in mind,
// this method is used during the DOCTORSTAGE.
func (r *Room) CheckIfDoctorSaved() bool {
	for _, pl := range r.players {
		if pl.Job == DOCTOR && r.stage == DOCTORSTAGE && pl.Voted == false {
			return false
		}
	}
	return true
}

// FindChosenPlayerToDie loops through the alive players in the current
// room to find out who has been selected to Die or Be Imprisoned.
// Returns nil, if there is no such player.
func (r *Room) FindChosenPlayerToDie() *Player {
	for index := range r.players {
		if r.players[index].Dead == false && r.players[index].Chosen == true {
			r.players[index].Dead = true
			return r.players[index]
		}
	}
	return nil
}

// HasDoctor is used to check if the DOCTOR of the room is still alive
func (r *Room) HasDoctor() bool {
	for _, pl := range r.players {
		if pl.Job == DOCTOR && pl.Dead == false {
			return true
		}
	}
	return false
}

// GameOver determines who has won the game.
// The main logic here is that, if MAFIA members become 0, CITIZENS win.
// Otherwise, if CITIZENS(incl the DOCTOR) become 1 or less than 1, MAFIA win.
// You might find it interesting that when CITIZENS(incl Doctor) and MAFIA
// both become 1, MAFIA win, this is because MAFIA will just shoot the alive CITIZEN.
func (r *Room) GameOver() (bool, Role) {
	aliveMafia := 0
	aliveCitizensDocs := 0
	for _, pl := range r.players {
		if pl.Job == MAFIA && pl.Dead == false {
			aliveMafia++
		} else if (pl.Job == CITIZEN || pl.Job == DOCTOR) && pl.Dead == false {
			aliveCitizensDocs++
		}
	}
	if aliveMafia == 0 {
		return true, CITIZEN
	} else if aliveCitizensDocs <= 1 && aliveMafia >= 1 {
		return true, MAFIA
	} else {
		return false, 0
	}
}

// End is used the finish the Game, after either side wins.
// Every player call its End method, and the current room is discarded.
func (r *Room) End() {
	for index := range r.players {
		r.players[index].End()
	}
	r.playing = false
	r = nil
}

// KickPlayer is used to remove the given player from the current room.
func (r *Room) KickPlayer(pl *Player) {
	wasRoomOwner := pl.RoomOwner
	for index := range r.players {
		if r.players[index] == pl {
			r.players[index].Dead = true
			if wasRoomOwner {
				r.players = append(r.players[:index], r.players[(index+1):]...)
				if len(r.players) != 0 {
					r.players[index+1].RoomOwner = true
					return
				}
			} else {
				r.players = append(r.players[:index], r.players[(index+1):]...)
				return
			}
		}
	}
}
