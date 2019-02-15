package src

import (
	"strconv"
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
	for index, _ := range players {
		players[index].Voted = true
		room.AddPlayer(&players[index])
	}
	room.Reset()
	for _, pl := range players {
		if pl.Voted == true {
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
	if room.GetMostVotedPlayer() != room.players[5] {
		t.Error("GetMostVotedPlayer not working!")
	}
}

func TestRoom_CheckIfAllVoted(t *testing.T) {
	room := Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	room.stage = ALLSTAGE
	for index, _ := range players {
		room.AddPlayer(&players[index])
		players[index].CastVote(players[0].Name)
	}
	if room.CheckIfAllVoted() == false {
		t.Error("Check if all voted not working!")
	}
}

func TestRoom_CheckIfMafiaVoted(t *testing.T) {
	room := Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	room.stage = MAFIASTAGE
	for index, _ := range players {
		room.AddPlayer(&players[index])
		players[index].Job = CITIZEN
		if index%2 == 0 {
			players[index].Job = MAFIA
			players[index].CastVote(players[0].Name)
		}
	}
	if room.CheckIfMafiaVoted() == false {
		t.Error("Check if MAFIA voted not working!")
	}
}

func TestRoom_StartGame(t *testing.T) {
	gameRoom := Room{}
	players := [6]Player{}
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	gameRoom.StartGame()
	var doctorSeen, mafiaSeen, citizenSeen bool
	for _, pl := range gameRoom.players {
		if doctorSeen && mafiaSeen && citizenSeen {
			break
		} else if pl.Job == MAFIA {
			mafiaSeen = true
		} else if pl.Job == DOCTOR {
			doctorSeen = true
		} else if pl.Job == CITIZEN {
			citizenSeen = true
		}
	}
	if !(doctorSeen && mafiaSeen && citizenSeen) {
		t.Error("Job Randomizer Fault!")
	}
	if !(gameRoom.playing && gameRoom.stage == MAFIASTAGE) {
		t.Error("Start Game Error!")
	}
}

func TestFindRoom(t *testing.T) {
	rooms := [6]Room{}
	for index, _ := range rooms {
		rooms[index].name = strconv.Itoa(index)
	}
	roomsLice := rooms[:]
	if roomsLice[5].name != FindRoom(&roomsLice, "5").name {
		t.Error("Finding rooms doesn't work!")
	}
}

func TestRoom_CanGoToNextStage(t *testing.T) {
	gameRoom := Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	for index, _ := range gameRoom.players {
		if gameRoom.players[index].Job == MAFIA {
			gameRoom.players[index].CastVote("BadBoy")
		}
	}
	if !(gameRoom.CanGoToNextStage()) {
		t.Error("Can go to next stage doesn't work!")
	}
}

func TestRoom_CheckIfDoctorSaved(t *testing.T) {
	gameRoom := Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	for index, _ := range gameRoom.players {
		if gameRoom.players[index].Job == DOCTOR {
			gameRoom.players[index].CastVote("BadBoy")
		}
	}
	if !(gameRoom.CheckIfDoctorSaved()) {
		t.Error("CheckIfDoctor Saved doesn't work!")
	}
}

func TestRoom_End(t *testing.T) {
	gameRoom := &Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	players[0].Job = DOCTOR
	gameRoom.End()
	if players[0].Job == DOCTOR {
		t.Error("End doesn't work!")
	}
}

func TestRoom_FindChosenPlayerToDie(t *testing.T) {
	gameRoom := &Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	players[0].Chosen = true
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	chosenPlayer := gameRoom.FindChosenPlayerToDie()
	if chosenPlayer != &players[0] {
		t.Error("Find chosen player doesn't work!")
	}
}

func TestRoom_FindPlayer(t *testing.T) {
	gameRoom := &Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	chosenPlayer := gameRoom.FindPlayer("BadBoy")
	if chosenPlayer != &players[0] {
		t.Error("Find player doesn't work!")
	}
}

func TestRoom_GameOver(t *testing.T) {
	gameRoom := &Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	for index, _ := range gameRoom.players {
		if gameRoom.players[index].Job == MAFIA {
			gameRoom.players[index].Dead = true
		}
	}
	if res, winner := gameRoom.GameOver(); res {
		if winner != CITIZEN {
			t.Error("GameOver not working!")
		}
	}
}

func TestRoom_GetName(t *testing.T) {
	gameRoom := &Room{}
	gameRoom.name = "GoTown"
	if gameRoom.GetName() != gameRoom.name {
		t.Error("GetName doesn't work!")
	}
}

func TestRoom_GetPlayers(t *testing.T) {
	gameRoom := &Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	players[0].Chosen = true
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	if players[0].Name != gameRoom.GetPlayers()[0].Name {
		t.Error("GetPlayers doesn't work!")
	}
}

func TestRoom_GetStage(t *testing.T) {
	gameRoom := &Room{}
	if gameRoom.GetStage() != MAFIASTAGE {
		t.Error("GetStage doesn't work!")
	}
}

func TestRoom_HasDoctor(t *testing.T) {
	gameRoom := &Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	players[0].Job = DOCTOR
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	if !(gameRoom.HasDoctor()) {
		t.Error("HasDoctor doesn't work!")
	}
}

func TestRoom_IsPlaying(t *testing.T) {
	gameRoom := &Room{}
	gameRoom.playing = true
	if !(gameRoom.IsPlaying()) {
		t.Error("IsPlaying doesn't work!")
	}
}

func TestRoom_NextStage(t *testing.T) {
	gameRoom := &Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	for index, _ := range gameRoom.players {
		if gameRoom.players[index].Job == MAFIA {
			gameRoom.players[index].CastVote("BadBoy")
		}
	}
	if gameRoom.NextStage(); gameRoom.stage != DOCTORSTAGE {
		t.Error("NextStage not working!")
	}
}

func TestRoom_SetName(t *testing.T) {
	gameRoom := &Room{}
	gameRoom.SetName("GoTown")
	if gameRoom.GetName() != gameRoom.name {
		t.Error("SetName not working!")
	}
}

func TestRoom_KickPlayer(t *testing.T) {
	gameRoom := &Room{}
	players := [6]Player{}
	players[0].Name = "BadBoy"
	for index, _ := range players {
		gameRoom.AddPlayer(&players[index])
	}
	if gameRoom.KickPlayer(&players[5]); len(gameRoom.players) != 5 {
		t.Error("KickPlayer, doesn't remove the player from the players slice!")
	}
	gameRoom.KickPlayer(&players[0])
	if gameRoom.GetOwner().Name == "BadBoy" {
		t.Error("KickPlayer, doesn't assign new owner!")
	}
	for i := 0; i < 4; i++ {
		gameRoom.KickPlayer(gameRoom.players[0])
	}
	if len(gameRoom.players) != 0 {
		t.Error("KickPlayer, for every member must kick everyone!")
	}
	gameRoom.AddPlayer(&players[0])
	if gameRoom.GetOwner() != &players[0] {
		t.Error("AddPlayer, on emptied by KickPlayer room, doesn't assing new owner!")
	}
}
