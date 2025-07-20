package games

import (
	"fmt"
)

func init() {
	RegisterGame("tic-tac-toe", NewTicTacToe)
}

type TicTacToe struct {
	Board         [9]string `json:"board"`
	CurrentTurn   int       `json:"currentTurn"`
	Winner        string    `json:"winner"`
	playerSymbols [2]string
	moves         int
}

func NewTicTacToe() (Game, error) {
	t := &TicTacToe{}
	t.Reset()
	return t, nil
}

func (t *TicTacToe) HandleMove(playerIndex int, move any) error {
	if t.IsGameOver() {
		return fmt.Errorf("game is already over")
	}

	if playerIndex != t.CurrentTurn {
		return fmt.Errorf("it's not your turn")
	}

	moveData, ok := move.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid move format")
	}

	cellIndexFloat, ok := moveData["cellIndex"].(float64)
	if !ok {
		return fmt.Errorf("cellIndex must be a number")
	}
	cellIndex := int(cellIndexFloat)

	if cellIndex < 0 || cellIndex > 8 {
		return fmt.Errorf("cell index out of bounds")
	}

	if t.Board[cellIndex] != "" {
		return fmt.Errorf("cell is already occupied")
	}

	symbol := t.playerSymbols[playerIndex]
	t.Board[cellIndex] = symbol
	t.moves++

	if t.checkWinner(symbol) {
		t.Winner = symbol
	} else if t.moves == 9 {
		t.Winner = "draw"
	}

	t.CurrentTurn = (t.CurrentTurn + 1) % 2

	return nil
}

func (t *TicTacToe) GetGameState() any {
	return t
}

func (t *TicTacToe) IsGameOver() bool {
	return t.Winner != ""
}

func (t *TicTacToe) GetWinner() string {
	return t.Winner
}

func (t *TicTacToe) Reset() {
	t.Board = [9]string{}
	for i := range t.Board {
		t.Board[i] = ""
	}
	t.CurrentTurn = 0
	t.Winner = ""
	t.playerSymbols = [2]string{"X", "O"}
	t.moves = 0
}

func (t *TicTacToe) checkWinner(symbol string) bool {
	// All possible winning lines
	lines := [8][3]int{
		{0, 1, 2}, // Horizontal
		{3, 4, 5},
		{6, 7, 8},
		{0, 3, 6}, // Vertical
		{1, 4, 7},
		{2, 5, 8},
		{0, 4, 8}, // Diagonal
		{2, 4, 6},
	}

	for _, line := range lines {
		if t.Board[line[0]] == symbol && t.Board[line[1]] == symbol && t.Board[line[2]] == symbol {
			return true
		}
	}
	return false
}
