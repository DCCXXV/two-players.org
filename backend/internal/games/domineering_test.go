package games

import (
	"testing"
)

func makeDomineeringMovePayload(col int, row int) any {
	return map[string]any{"row": float64(row), "col": float64(col)}
}

func TestDomineering_HandleMove_ValidMove(t *testing.T) {
	game, _ := NewDomineering()

	err := game.HandleMove(0, makeDomineeringMovePayload(0, 0))
	if err != nil {
		t.Fatalf("Expected no error for a valid move but gor %v", err)
	}

	state := game.GetGameState().(*Domineering)
	if state.Board[0][0] != "h" || state.Board[0][1] != "h" {
		t.Errorf("Expected cell 0,0 and 0,1 to be 'h' but got 0,0:'%s' and 0,1:'%s'", state.Board[0][0], state.Board[0][1])
	}
	if state.CurrentTurn != 1 {
		t.Errorf("Expected CurrentTurn to be 1, but got %d", state.CurrentTurn)
	}

	err = game.HandleMove(1, makeDomineeringMovePayload(7, 7))
	if err != nil {
		t.Fatalf("Expected no error for a valid move but gor %v", err)
	}

	state = game.GetGameState().(*Domineering)
	if state.Board[7][7] != "v" || state.Board[6][7] != "v" {
		t.Errorf("Expected cell 7,7 to be 'v' and 6,7 to be 'v' but got 7,7:'%s' and 6,7:'%s'", state.Board[7][7], state.Board[6][7])
	}
	if state.CurrentTurn != 0 {
		t.Errorf("Expected CurrentTurn to be 1, but got %d", state.CurrentTurn)
	}
}

func TestDomineering_HandleMove_OccupiedCell(t *testing.T) {
	game, _ := NewDomineering()

	game.HandleMove(0, makeDomineeringMovePayload(0, 0))
	err := game.HandleMove(1, makeDomineeringMovePayload(0, 0))

	if err == nil {
		t.Fatal("Expected an error when moving to an occupied cell, but got nil")
	}

	expectedError := "cell is already occupied"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestDomineering_HandleMove_NotYourTurn(t *testing.T) {
	game, _ := NewDomineering()

	game.HandleMove(0, makeDomineeringMovePayload(0, 0))
	err := game.HandleMove(0, makeDomineeringMovePayload(1, 1))

	if err == nil {
		t.Fatal("Expected an error when moving out of turn, but got nil")
	}

	expectedError := "it's not your turn"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestDomineering_WinCondition(t *testing.T) {
	game, _ := NewDomineering()
	d := game.GetGameState().(*Domineering)

	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c += 2 {
			d.Board[r][c] = "h"
			d.Board[r][c+1] = "h"
		}
	}

	d.Board[7][6] = ""
	d.Board[7][7] = ""

	d.CurrentTurn = 0

	err := game.HandleMove(0, makeDomineeringMovePayload(6, 7))
	if err != nil {
		t.Fatalf("Expected no error for the winning move, but got %v", err)
	}

	if winner := game.GetWinner(); winner != "h" {
		t.Errorf("Expected winner to be 'h', but got '%s'", winner)
	}

	if !game.IsGameOver() {
		t.Fatal("Expected game to be over, but it wasn't")
	}
}
