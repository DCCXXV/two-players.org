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

	symbol := d.playerSymbols[playerIndex]

	if symbol == "H" {
		if col == 7 {
			if d.Board[row][col] != "" || d.Board[row][col-1] != "" {
				return fmt.Errorf("cell is already occupied")
			}
			d.Board[row][col] = symbol
			d.Board[row][col-1] = symbol
		} else {
			if d.Board[row][col] != "" || d.Board[row][col+1] != "" {
				return fmt.Errorf("cell is already occupied")
			}
			d.Board[row][col] = symbol
			d.Board[row][col+1] = symbol
		}
	} else if symbol == "V" {
		if row == 7 {
			if d.Board[row][col] != "" || d.Board[row-1][col] != "" {
				return fmt.Errorf("cell is already occupied")
			}
			d.Board[row][col] = symbol
			d.Board[row-1][col] = symbol
		} else {
			if d.Board[row][col] != "" || d.Board[row+1][col] != "" {
				return fmt.Errorf("cell is already occupied")
			}
			d.Board[row][col] = symbol
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
	d.playerSymbols = [2]string{"H", "V"}
	d.moves = 0
}

func (d *Domineering) checkWinner(symbol string) bool {
	otherPlayerSymbol := "V"
	if symbol == "V" {
		otherPlayerSymbol = "H"
	}

	return !d.hasValidMoves(otherPlayerSymbol)
}

func (d *Domineering) hasValidMoves(symbol string) bool {
	for r := range d.Board {
		for c := range d.Board[r] {
			if symbol == "H" {
				if c+1 < 8 && d.Board[r][c] == "" && d.Board[r][c+1] == "" {
					return true
				}
			} else {
				if r+1 < 8 && d.Board[r][c] == "" && d.Board[r+1][c] == "" {
					return true
				}
			}
		}
	}
	return false
}
