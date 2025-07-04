package games

import (
	"fmt"
)

func init() {
	// Register the Tic-Tac-Toe game type with its factory function.
	RegisterGame("tic-tac-toe", NewTicTacToe)
}

// TicTacToe implements the Game interface for a game of Tic-Tac-Toe.
type TicTacToe struct {
	// The 3x3 board, represented as a 9-element array. Stores "X" or "O".
	Board [9]string `json:"board"`
	// The index of the player whose turn it is (0 or 1).
	CurrentTurn int `json:"currentTurn"`
	// Stores the winner: "X", "O", "draw", or "" if the game is ongoing.
	Winner string `json:"winner"`
	// The symbols assigned to the two players, e.g., ["X", "O"].
	playerSymbols [2]string
	// The number of moves made so far.
	moves int
}

// NewTicTacToe creates a new, ready-to-play instance of a TicTacToe game.
func NewTicTacToe() (Game, error) {
	t := &TicTacToe{}
	t.Reset()
	return t, nil
}

// HandleMove processes a player's move. It's the core logic of the game.
func (t *TicTacToe) HandleMove(playerIndex int, move any) error {
	if t.IsGameOver() {
		return fmt.Errorf("game is already over")
	}

	if playerIndex != t.CurrentTurn {
		return fmt.Errorf("it's not your turn")
	}

	// 1. Validate the move format. The JSON payload is decoded into a map[string]any.
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

	// 2. Validate that the chosen cell is empty.
	if t.Board[cellIndex] != "" {
		return fmt.Errorf("cell is already occupied")
	}

	// 3. Update the board with the player's symbol.
	symbol := t.playerSymbols[playerIndex]
	t.Board[cellIndex] = symbol
	t.moves++

	// 4. Check for a winner or a draw.
	if t.checkWinner(symbol) {
		t.Winner = symbol
	} else if t.moves == 9 {
		t.Winner = "draw"
	}

	// 5. Switch turns.
	t.CurrentTurn = (t.CurrentTurn + 1) % 2

	return nil
}

// GetGameState returns the current state of the game, suitable for JSON serialization.
func (t *TicTacToe) GetGameState() any {
	// The struct is returned directly, as its fields are tagged for JSON.
	return t
}

// IsGameOver checks if the game has concluded.
func (t *TicTacToe) IsGameOver() bool {
	return t.Winner != ""
}

// GetWinner returns the winner of the game ("X", "O", "draw", or "").
func (t *TicTacToe) GetWinner() string {
	return t.Winner
}

// Reset brings the game back to its initial state.
func (t *TicTacToe) Reset() {
	t.Board = [9]string{}
	for i := range t.Board {
		t.Board[i] = ""
	}
	t.CurrentTurn = 0 // Player 0 ("X") always starts.
	t.Winner = ""
	t.playerSymbols = [2]string{"X", "O"}
	t.moves = 0
}

// checkWinner checks all winning conditions for the given symbol.
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
