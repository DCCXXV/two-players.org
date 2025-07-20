package games

import (
	"fmt"
)

func init() {
	RegisterGame("domineering", NewDomineering)
}

type Domineering struct {
	Board         [8][8]string `json:"board"`
	CurrentTurn   int          `json:"currentTurn"`
	Winner        string       `json:"winner"`
	playerSymbols [2]string
	moves         int
}

func NewDomineering() (Game, error) {
	d := &Domineering{}
	d.Reset()
	return d, nil
}

func (d *Domineering) HandleMove(playerIndex int, move any) error {
	if d.IsGameOver() {
		return fmt.Errorf("game is already over")
	}

	if playerIndex != d.CurrentTurn {
		return fmt.Errorf("it's not your turn")
	}

	moveData, ok := move.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid move format")
	}

	rowFloat, ok := moveData["row"].(float64)

	if !ok {
		return fmt.Errorf("row must be a number")
	}

	colFloat, ok := moveData["col"].(float64)

	if !ok {
		return fmt.Errorf("col must be a number")
	}

	row := int(rowFloat)
	col := int(colFloat)

	if row < 0 || row > 7 || col < 0 || col > 7 {
		return fmt.Errorf("out of bounds")
	}

	if d.Board[row][col] != "" {
		return fmt.Errorf("cell is already occupied")
	}

	symbol := d.playerSymbols[playerIndex]

	if symbol == "h" && col == 7 {
		d.Board[row][col] = symbol
		d.Board[row][col-1] = symbol
	} else if symbol == "v" && row == 7 {
		d.Board[row][col] = symbol
		d.Board[row-1][col] = symbol
	} else {
		d.Board[row][col] = symbol
		if symbol == "h" {
			d.Board[row][col+1] = symbol
		} else if symbol == "v" {
			d.Board[row+1][col] = symbol
		}
	}

	d.moves++

	if d.checkWinner(symbol) {
		d.Winner = symbol
	}

	d.CurrentTurn = (d.CurrentTurn + 1) % 2

	return nil
}

func (d *Domineering) GetGameState() any {
	return d
}

func (d *Domineering) IsGameOver() bool {
	return d.Winner != ""
}

func (d *Domineering) GetWinner() string {
	return d.Winner
}

func (d *Domineering) Reset() {
	d.Board = [8][8]string{}
	for i := range d.Board {
		for j := range d.Board[i] {
			d.Board[i][j] = ""
		}
	}
	d.CurrentTurn = 0
	d.Winner = ""
	d.playerSymbols = [2]string{"h", "v"}
	d.moves = 0
}

func (d *Domineering) checkWinner(symbol string) bool {
	otherPlayerSymbol := "v"
	if symbol == "v" {
		otherPlayerSymbol = "h"
	}

	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			if otherPlayerSymbol == "h" {
				if c+1 < 8 && d.Board[r][c] == "" && d.Board[r][c+1] == "" {
					return false
				}
			} else {
				if r+1 < 8 && d.Board[r][c] == "" && d.Board[r+1][c] == "" {
					return false
				}
			}
		}
	}

	return true
}
