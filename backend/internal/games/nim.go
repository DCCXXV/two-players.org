package games

import (
	"encoding/json"
	"errors"
)

func init() {
	RegisterGame("nim", func() (Game, error) {
		return NewNimGame(), nil
	})
}

type NimGame struct {
	Sticks      int    `json:"sticks"`
	CurrentTurn int    `json:"currentTurn"`
	Winner      string `json:"winner"`
}

func NewNimGame() *NimGame {
	return &NimGame{
		Sticks:      21,
		CurrentTurn: 0,
		Winner:      "",
	}
}

func (g *NimGame) HandleMove(playerIndex int, move any) error {
	if g.Winner != "" {
		return errors.New("game is already over")
	}

	if playerIndex != g.CurrentTurn {
		return errors.New("not your turn")
	}

	moveMap, ok := move.(map[string]interface{})
	if !ok {
		return errors.New("invalid move format")
	}

	sticksToTake, ok := moveMap["sticks"].(float64)
	if !ok {
		return errors.New("invalid sticks value")
	}

	sticks := int(sticksToTake)

	if sticks < 1 || sticks > 3 {
		return errors.New("you can only take 1, 2, or 3 sticks")
	}

	if sticks > g.Sticks {
		return errors.New("not enough sticks remaining")
	}

	g.Sticks -= sticks

	if g.Sticks == 0 {
		// whoever took the last stick loses, so the winner is the other player
		winnerIndex := 1 - playerIndex
		if winnerIndex == 0 {
			g.Winner = "P1"
		} else {
			g.Winner = "P2"
		}
		return nil
	}

	g.CurrentTurn = 1 - g.CurrentTurn
	return nil
}

func (g *NimGame) GetGameState() interface{} {
	return g
}

func (g *NimGame) IsGameOver() bool {
	return g.Winner != ""
}

func (g *NimGame) GetWinner() string {
	return g.Winner
}

func (g *NimGame) Reset() {
	g.Sticks = 21
	g.CurrentTurn = 0
	g.Winner = ""
}

func (g *NimGame) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sticks      int    `json:"sticks"`
		CurrentTurn int    `json:"currentTurn"`
		Winner      string `json:"winner"`
	}{
		Sticks:      g.Sticks,
		CurrentTurn: g.CurrentTurn,
		Winner:      g.Winner,
	})
}
