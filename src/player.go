package src

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

/*func (pl Player) String() string {
	return fmt.Sprintf("Name: %v and Job: %v", pl.Name, pl.Job)
}*/

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
func (pl *Player) StartGame() bool {
	if pl.RoomOwner == true && len(pl.Room.players) >= 6 {
		pl.Room.StartGame()
		return true
	}
	return false
}

func (pl *Player) ResetRound() {
	pl.Invulnerable = false
	pl.Votes = 0
}

func (pl *Player) SetVotes(score uint){
	pl.Votes = score
}

func (pl *Player) CastVote(votedPlayerName string){
	votedPlayer := pl.Room.FindPlayer(votedPlayerName)
	if votedPlayer != nil {
		votedPlayer.IncrementVote()
	}
}

func (pl *Player) IncrementVote(){
	pl.Votes+=1
}