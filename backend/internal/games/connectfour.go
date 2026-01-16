package games

import (
	"fmt"
)

const (
	connectFourRows = 6
	connectFourCols = 7
)

func init() {
	RegisterGame("connect-four", NewConnectFour)
}

type ConnectFour struct {
	Board         [connectFourRows][connectFourCols]string `json:"board"`
	CurrentTurn   int                                      `json:"currentTurn"`
	Winner        string                                   `json:"winner"`
	WinningCells  [][2]int                                 `json:"winningCells"`
	playerSymbols [2]string
	moves         int
}

func NewConnectFour() (Game, error) {
	c := &ConnectFour{}
	c.Reset()
	return c, nil
}

func (c *ConnectFour) HandleMove(playerIndex int, move any) error {
	if c.IsGameOver() {
		return fmt.Errorf("game is already over")
	}

	if playerIndex != c.CurrentTurn {
		return fmt.Errorf("it's not your turn")
	}

	moveData, ok := move.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid move format")
	}

	colFloat, ok := moveData["column"].(float64)
	if !ok {
		return fmt.Errorf("column must be a number")
	}
	col := int(colFloat)

	if col < 0 || col >= connectFourCols {
		return fmt.Errorf("column out of bounds")
	}

	row := c.getLowestEmptyRow(col)
	if row == -1 {
		return fmt.Errorf("column is full")
	}

	symbol := c.playerSymbols[playerIndex]
	c.Board[row][col] = symbol
	c.moves++

	if winningCells := c.checkWinner(row, col, symbol); winningCells != nil {
		c.Winner = symbol
		c.WinningCells = winningCells
	} else if c.moves == connectFourRows*connectFourCols {
		c.Winner = "draw"
	}

	c.CurrentTurn = (c.CurrentTurn + 1) % 2

	return nil
}

func (c *ConnectFour) getLowestEmptyRow(col int) int {
	for row := connectFourRows - 1; row >= 0; row-- {
		if c.Board[row][col] == "" {
			return row
		}
	}
	return -1
}

func (c *ConnectFour) checkWinner(row, col int, symbol string) [][2]int {
	directions := [][2]int{
		{0, 1},
		{1, 0},
		{1, 1},
		{1, -1},
	}

	for _, dir := range directions {
		cells := c.countInDirection(row, col, dir[0], dir[1], symbol)
		if len(cells) >= 4 {
			return cells
		}
	}
	return nil
}

func (c *ConnectFour) countInDirection(row, col, dRow, dCol int, symbol string) [][2]int {
	cells := [][2]int{{row, col}}

	for i := 1; i < 4; i++ {
		r, co := row+dRow*i, col+dCol*i
		if r < 0 || r >= connectFourRows || co < 0 || co >= connectFourCols {
			break
		}
		if c.Board[r][co] != symbol {
			break
		}
		cells = append(cells, [2]int{r, co})
	}

	for i := 1; i < 4; i++ {
		r, co := row-dRow*i, col-dCol*i
		if r < 0 || r >= connectFourRows || co < 0 || co >= connectFourCols {
			break
		}
		if c.Board[r][co] != symbol {
			break
		}
		cells = append(cells, [2]int{r, co})
	}

	return cells
}

func (c *ConnectFour) GetGameState() any {
	return c
}

func (c *ConnectFour) IsGameOver() bool {
	return c.Winner != ""
}

func (c *ConnectFour) GetWinner() string {
	return c.Winner
}

func (c *ConnectFour) Reset() {
	c.Board = [connectFourRows][connectFourCols]string{}
	c.CurrentTurn = 0
	c.Winner = ""
	c.WinningCells = nil
	c.playerSymbols = [2]string{"R", "B"}
	c.moves = 0
}
