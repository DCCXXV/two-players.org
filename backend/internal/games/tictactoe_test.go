package games

import (
	"testing"
)

func makeTicTacToeMovePayload(cellIndex int) any {
	return map[string]any{"cellIndex": float64(cellIndex)}
}

func TestTicTacToe_HandleMove_ValidMove(t *testing.T) {
	game, _ := NewTicTacToe()

	err := game.HandleMove(0, makeTicTacToeMovePayload(0))
	if err != nil {
		t.Fatalf("Expected no error for a valid move, but got %v", err)
	}

	state := game.GetGameState().(*TicTacToe)
	if state.Board[0] != "X" {
		t.Errorf("Expected cell 0 to be 'X', but got '%s'", state.Board[0])
	}
	if state.CurrentTurn != 1 {
		t.Errorf("Expected CurrentTurn to be 1, but got %d", state.CurrentTurn)
	}
}

func TestTicTacToe_HandleMove_OccupiedCell(t *testing.T) {
	game, _ := NewTicTacToe()

	game.HandleMove(0, makeTicTacToeMovePayload(0)) // Player 0 moves to 0
	err := game.HandleMove(1, makeTicTacToeMovePayload(0))

	if err == nil {
		t.Fatal("Expected an error when moving to an occupied cell, but got nil")
	}

	expectedError := "cell is already occupied"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestTicTacToe_HandleMove_NotYourTurn(t *testing.T) {
	game, _ := NewTicTacToe()

	game.HandleMove(0, makeTicTacToeMovePayload(0)) // Player 0 moves
	err := game.HandleMove(0, makeTicTacToeMovePayload(1))

	if err == nil {
		t.Fatal("Expected an error when moving out of turn, but got nil")
	}

	expectedError := "it's not your turn"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestTicTacToe_WinCondition(t *testing.T) {
	game, _ := NewTicTacToe()
	moves := []struct {
		playerIndex int
		cellIndex   int
	}{
		{0, 0}, // X
		{1, 4}, // O
		{0, 1}, // X
		{1, 5}, // O
		{0, 2}, // X wins
	}

	for _, move := range moves {
		err := game.HandleMove(move.playerIndex, makeTicTacToeMovePayload(move.cellIndex))
		if err != nil {
			t.Fatalf("Move sequence failed at player %d, cell %d: %v", move.playerIndex, move.cellIndex, err)
		}
	}

	if !game.IsGameOver() {
		t.Fatal("Expected game to be over, but it wasn't")
	}

	if game.GetWinner() != "X" {
		t.Errorf("Expected winner to be 'X', but got '%s'", game.GetWinner())
	}
}

func TestTicTacToe_DrawCondition(t *testing.T) {
	game, _ := NewTicTacToe()
	moves := []struct {
		playerIndex int
		cellIndex   int
	}{
		{0, 0}, {1, 1}, {0, 2},
		{1, 3}, {0, 5}, {1, 4},
		{0, 6}, {1, 8}, {0, 7},
	}

	for _, move := range moves {
		game.HandleMove(move.playerIndex, makeTicTacToeMovePayload(move.cellIndex))
	}

	if !game.IsGameOver() {
		t.Fatal("Expected game to be over after a draw, but it wasn't")
	}

	if game.GetWinner() != "draw" {
		t.Errorf("Expected winner to be 'draw', but got '%s'", game.GetWinner())
	}
}

func TestTicTacToe_Reset(t *testing.T) {
	game, _ := NewTicTacToe()
	game.HandleMove(0, makeTicTacToeMovePayload(0))
	game.HandleMove(1, makeTicTacToeMovePayload(1))

	game.Reset()

	state := game.GetGameState().(*TicTacToe)

	if state.moves != 0 {
		t.Errorf("Expected moves to be 0 after reset, but got %d", state.moves)
	}
	if state.Winner != "" {
		t.Errorf("Expected winner to be empty after reset, but got '%s'", state.Winner)
	}
	if state.CurrentTurn != 0 {
		t.Errorf("Expected current turn to be 0 after reset, but got %d", state.CurrentTurn)
	}

	for i, cell := range state.Board {
		if cell != "" {
			t.Errorf("Expected board cell %d to be empty after reset, but got '%s'", i, cell)
		}
	}
}
