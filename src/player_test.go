package src

import (
	"testing"
)

func TestPlayer_Die(t *testing.T) {
	player := Player{}

	player.Die()

	if player.Dead == false {
		t.Error("Player didn't die, Die not working!")
	}
}

func TestPlayer_Save(t *testing.T) {
	player := &Player{}
	player.Save()
	if player.Chosen == true {
		t.Error("Player wasn't saved, Save not working!")
	}
}

func TestPlayer_ResetRound(t *testing.T) {
	players := [6]Player{}
	for _, pl := range players {
		pl.Votes++
	}
	for _, pl := range players {
		pl.ResetRound()
		if pl.Votes != 0 || pl.Voted != false {
			t.Error("Reset round not working!")
		}
	}

}

func TestPlayer_CastVote(t *testing.T) {
	room := Room{}
	player1 := Player{}
	votedPl := Player{}
	votedPl.Name = "BadBoy"
	room.AddPlayer(&player1)
	room.AddPlayer(&votedPl)
	player1.CastVote(votedPl.Name)
	if votedPl.Votes != 1 {
		t.Error("Cast vote not working!")
	}
}

func TestPlayer_AssignChosen(t *testing.T) {
	coolPlayer := &Player{}
	if coolPlayer.AssignChosen(); !(coolPlayer.Chosen) {
		t.Error("AssignChosen doesn't work!")
	}
}

func TestPlayer_End(t *testing.T) {
	coolPlayer := &Player{}
	coolPlayer.Job = DOCTOR
	if coolPlayer.End(); coolPlayer.Job != CITIZEN {
		t.Error("End for player doesn't work!")
	}
}

func TestPlayer_IncrementVote(t *testing.T) {
	coolPlayer := &Player{}
	if coolPlayer.IncrementVote(); coolPlayer.Votes != 1 {
		t.Error("IncrementVote for player doesn't work!")
	}
}

func TestPlayer_IsEligibleToChat(t *testing.T) {
	coolPlayer := &Player{}
	coolRoom := &Room{}
	coolPlayer.Room = coolRoom
	if coolPlayer.Job = MAFIA; !(coolPlayer.IsEligibleToChat()) {
		t.Error("Is eligible to chat doesn't work!")
	}
}

func TestPlayer_SetVotes(t *testing.T) {
	coolPlayer := &Player{}
	if coolPlayer.SetVotes(10); coolPlayer.Votes != 10 {
		t.Error("SetVotes doesn't work!")
	}
}

func TestPlayer_StartGame(t *testing.T) {
	coolPlayer := &Player{}
	coolRoom := Room{}
	players := [6]Player{}
	coolRoom.AddPlayer(coolPlayer)
	for index, _ := range players {
		coolRoom.AddPlayer(&players[index])
	}
	if !(coolPlayer.StartGame()) {
		t.Error("StartGame for player as room owner doesn't work!")
	}
}
