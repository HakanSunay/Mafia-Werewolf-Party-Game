package src

import (
	"testing"
)

func TestCreateRoom(t *testing.T) {
	player := Player{}
	player.Name = "Hakan"
	room := CreateRoom("GoRoom", &player)
	for _, pl := range room.players {
		if pl.RoomOwner == true {
			if pl.Name != player.Name {
				t.Error("Room owner names don't match, CreateRoom not working!")
			}
		}
	}
}

func TestRoom_AddPlayer(t *testing.T) {
	room := Room{}
	player := Player{}
	room.AddPlayer(&player)
	if len(room.players) == 0 {
		t.Error("Empty room, AddPlayer not working!")
	}
}

func TestRoom_AddPlayerBecomesRoomOwner(t *testing.T) {
	room := Room{}
	player := Player{}
	room.AddPlayer(&player)
	if room.players[0].RoomOwner == false {
		t.Error("Room has no owner, AddPlayer not working!")
	}
}

func TestRoom_AddPlayer2(t *testing.T) {
	room := Room{}
	player1 := Player{}
	player2 := Player{}
	player1.Name = "Player1"
	player2.Name = "Player2"
	room.AddPlayer(&player1)
	room.AddPlayer(&player2)
	if room.players[0].Name == room.players[1].Name {
		t.Error("Adding multiple players doesn't work!")
	}
}

func TestRoom_Reset(t *testing.T) {
	room := Room{}
	players := [6]Player{}
	for _, pl := range players {
		pl.Invulnerable = true
		room.AddPlayer(&pl)
	}
	room.Reset()
	for _, pl := range players {
		if pl.Invulnerable == true {
			t.Error("Reset not working!")
		}
	}
}

func TestRoom_GetOwner(t *testing.T) {
	room := Room{}
	player := Player{}
	room.AddPlayer(&player)
	if room.GetOwner() != &player {
		t.Error("GetOwner not working!")
	}
}

func TestRoom_GetMostVotedPlayer(t *testing.T) {
	room := Room{}
	players := [6]Player{}
	var num uint = 0
	for index, _ := range players {
		players[index].SetVotes(num)
		num += 1
		room.AddPlayer(&players[index])
	}
	if room.GetMostVotedPlayer() != room.players[5]{
		t.Error("GetMostVotedPlayer not working!")
	}
}