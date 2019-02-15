package src

const ROLECOUNT = 3
const MINIMUMPLAYERCOUNT = 4

type Role int

const (
	MAFIA Role = iota
	DOCTOR
	CITIZEN
)

type Player struct {
	RoomOwner bool
	Room      *Room
	Number    uint
	Name      string
	Job       Role
	Votes     uint
	Dead      bool
	Voted     bool
	Chosen    bool
}

// Save is called when a player of Role DOCTOR
// decides to Save the current player from dying.
// Keep in mind that the current player must not be DEAD.
func (pl *Player) Save() {
	if pl.Dead == false {
		pl.Chosen = false
	}
}

// Die is used on the current player if
// they have been chosen to die by the others.
// Keep in mind that the current player must not be DEAD.
func (pl *Player) Die() {
	if pl.Dead == false {
		pl.Dead = true
	}
}

// AssignChosen is used by MAFIA members, when they vote
// to kill a player. This method is necessary, because
// following the MAFIA vote, the room DOCTOR has a chance
// to Save a Chosen Player.
// Keep in mind that the current player must not be DEAD.
func (pl *Player) AssignChosen() {
	if pl.Dead == false {
		pl.Chosen = true
	}
}

// CreateRoom is used to initialize a room using the name argument.
func (pl *Player) CreateRoom(roomName string) *Room {
	return CreateRoom(roomName, pl)
}

// StartGame is a boolean method that starts the game
// after the owner of the room uses the #START_GAME command.
// The method returns false, if any of the 2 conditions are not met.
// 1 : The player calling the method, must be the owner of the room.
// 2 : The count of the players in the room must be higher than MINIMUMPLAYERCOUNT.
func (pl *Player) StartGame() bool {
	if pl.RoomOwner == true && len(pl.Room.players) >= MINIMUMPLAYERCOUNT {
		pl.Room.StartGame()
		return true
	}
	return false
}

// ResetRound is used to reset the stats of the current player, after
// a the room proceeds into the next stage.
func (pl *Player) ResetRound() {
	pl.Votes = 0
	pl.Voted = false
}

// SetVotes is a simple setter for the Votes of the current player.
func (pl *Player) SetVotes(score uint) {
	pl.Votes = score
}

// CastVote is a method used to blame or save the player specified
// with the votedPlayerName string.
// Keep in mind that the current player must not be DEAD.
func (pl *Player) CastVote(votedPlayerName string) {
	if pl.Voted == false && pl.Dead == false {
		votedPlayer := pl.Room.FindPlayer(votedPlayerName)
		if votedPlayer != nil {
			votedPlayer.IncrementVote()
			pl.Voted = true
		}
	}
}

// IncrementVote is used to ++ the votes of the current player.
func (pl *Player) IncrementVote() {
	pl.Votes += 1
}

// IsEligibleToChat is used to determine if the current player
// can chat. This depends on the current stage of the game or
// the alive status of the current player.
func (pl *Player) IsEligibleToChat() bool {
	if pl.Room != nil {
		if pl.Room.playing {
			if pl.Dead == true {
				return false
			} else if pl.Room.stage == ALLSTAGE {
				return true
			} else if pl.Job == MAFIA && pl.Room.stage == MAFIASTAGE {
				return true
			} else if pl.Job == DOCTOR && pl.Room.stage == DOCTORSTAGE {
				return true
			} else {
				return false
			}
		}
	}
	return true
}

// End is a method used to reset all of the stats of the current player,
// after finishing the game. This is used so that the players return to the
// Lobby, where they can continue play the game by chatting or creating a new room and so on.
func (pl *Player) End() {
	pl.RoomOwner = false
	pl.Job = CITIZEN
	pl.Dead = false
	pl.Chosen = false
	pl.Voted = false
	pl.Votes = 0
	pl.Room = nil
}
