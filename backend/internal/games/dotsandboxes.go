package games

import (
	"fmt"
)

func init() {
	RegisterGame("dots-and-boxes", NewDotsAndBoxes)
}

type DotsAndBoxes struct {
	HLines         [5][4]string `json:"hLines"`
	VLines         [4][5]string `json:"vLines"`
	Boxes          [4][4]string `json:"boxes"`
	CurrentTurn    int          `json:"currentTurn"`
	Winner         string       `json:"winner"`
	BoxesCompleted int          `json:"boxesCompleted"`
	Scores         [2]int       `json:"scores"`
	playerSymbols  [2]string
}

func NewDotsAndBoxes() (Game, error) {
	d := &DotsAndBoxes{}
	d.Reset()
	return d, nil
}

func (d *DotsAndBoxes) HandleMove(playerIndex int, move any) error {
	if d.IsGameOver() {
		return fmt.Errorf("game is already over")
	}

	if playerIndex != d.CurrentTurn {
		return fmt.Errorf("not your turn")
	}

	moveData, ok := move.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid move format")
	}

	lineType, ok := moveData["type"].(string)
	if !ok {
		return fmt.Errorf("type must be a string")
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

	symbol := d.playerSymbols[playerIndex]

	switch lineType {
	case "h":
		if row < 0 || row > 4 || col < 0 || col > 3 {
			return fmt.Errorf("horizontal line out of bounds")
		}
		if d.HLines[row][col] != "" {
			return fmt.Errorf("line is already occupied")
		}
		d.HLines[row][col] = symbol

	case "v":
		if row < 0 || row > 3 || col < 0 || col > 4 {
			return fmt.Errorf("vertical line out of bounds")
		}
		if d.VLines[row][col] != "" {
			return fmt.Errorf("line is already occupied")
		}
		d.VLines[row][col] = symbol

	default:
		return fmt.Errorf("invalid line type: must be 'h' or 'v'")
	}

	newBoxes := d.checkNewBoxes(lineType, row, col, symbol)

	d.BoxesCompleted = newBoxes
	d.Scores[playerIndex] += newBoxes

	if newBoxes == 0 {
		d.CurrentTurn = (d.CurrentTurn + 1) % 2
	}

	if d.checkGameOver() {
		d.determineWinner()
	}

	return nil
}

func (d *DotsAndBoxes) checkNewBoxes(lineType string, row, col int, symbol string) int {
	completed := 0

	if lineType == "h" {
		if row > 0 && d.isBoxComplete(row-1, col) && d.Boxes[row-1][col] == "" {
			d.Boxes[row-1][col] = symbol
			completed++
		}
		if row < 4 && d.isBoxComplete(row, col) && d.Boxes[row][col] == "" {
			d.Boxes[row][col] = symbol
			completed++
		}
	} else {
		if col > 0 && d.isBoxComplete(row, col-1) && d.Boxes[row][col-1] == "" {
			d.Boxes[row][col-1] = symbol
			completed++
		}
		if col < 4 && d.isBoxComplete(row, col) && d.Boxes[row][col] == "" {
			d.Boxes[row][col] = symbol
			completed++
		}
	}

	return completed
}

func (d *DotsAndBoxes) isBoxComplete(row, col int) bool {
	if row < 0 || row > 3 || col < 0 || col > 3 {
		return false
	}
	return d.HLines[row][col] != "" &&
		d.HLines[row+1][col] != "" &&
		d.VLines[row][col] != "" &&
		d.VLines[row][col+1] != ""
}

func (d *DotsAndBoxes) checkGameOver() bool {
	for i := range 4 {
		for j := range 4 {
			if d.Boxes[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func (d *DotsAndBoxes) determineWinner() {
	if d.Scores[0] > d.Scores[1] {
		d.Winner = d.playerSymbols[0]
	} else if d.Scores[1] > d.Scores[0] {
		d.Winner = d.playerSymbols[1]
	} else {
		d.Winner = "draw"
	}
}

func (d *DotsAndBoxes) GetGameState() any {
	return d
}

func (d *DotsAndBoxes) IsGameOver() bool {
	return d.Winner != ""
}

func (d *DotsAndBoxes) GetWinner() string {
	return d.Winner
}

func (d *DotsAndBoxes) Reset() {
	d.HLines = [5][4]string{}
	d.VLines = [4][5]string{}
	d.Boxes = [4][4]string{}
	d.CurrentTurn = 0
	d.Winner = ""
	d.BoxesCompleted = 0
	d.Scores = [2]int{}
	d.playerSymbols = [2]string{"P1", "P2"}
}
