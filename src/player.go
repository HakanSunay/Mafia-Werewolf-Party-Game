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
	roomOwner    bool
	room 		 Room
	number       uint
	Name         string
	Job          Role
	invulnerable bool
	votes        uint
}

func (pl *Player) Save() {
	pl.invulnerable = true
}

func (pl *Player) Kill() {

}

func (pl *Player) CreateRoom(roomName string){
	CreateRoom(roomName, pl)
}

func (pl* Player) Blame(plName string){

}

func (pl* Player) StartGame(){
	if pl.roomOwner == true && len(pl.room.players) >= 6 {
		pl.room.StartGame()
	}
}