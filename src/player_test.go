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
	player := Player{}

	player.Save()

	if player.Invulnerable == false {
		t.Error("Player wasn't saved, Save not working!")
	}
}

func TestPlayer_Blame(t *testing.T) {
	player1 := Player{}
	player2 := Player{}

	player1.Blame(&player2)

	if player2.Votes == 0 {
		t.Error("Vote not counted, Blame not working!")
	}
}

func TestPlayer_ResetRound(t *testing.T) {
	players := [6]Player{}
	for _, pl := range players {
		pl.Votes++
	}
	for _, pl := range players {
		pl.ResetRound()
		if pl.Votes != 0 {
			t.Error("Reset round not working!")
		}
	}

}
