package games

import "fmt"

// Game defines the interface that each game must implement.
type Game interface {
	// HandleMove processes a move from a player.
	// The format of 'move' is specific to each game.
	HandleMove(playerIndex int, move any) error

	// GetGameState returns the current state of the game to be sent to clients.
	GetGameState() any

	// IsGameOver checks if the game has finished.
	IsGameOver() bool

	// GetWinner returns the identifier of the winner, if any.
	GetWinner() string

	// Reset resets the game to its initial state.
	Reset()
}

// Factory is a function that creates a new instance of a game.
type Factory func() (Game, error)

var gameFactories = make(map[string]Factory)

// RegisterGame is used to register a new type of game.
func RegisterGame(gameType string, factory Factory) {
	if factory == nil {
		panic("game factory cannot be nil")
	}
	_, registered := gameFactories[gameType]
	if registered {
		// Optional: you might want this to be a panic if it occurs during init.
		// For now, we leave it as a log or a silent warning.
		return
	}
	gameFactories[gameType] = factory
}

// NewGame creates an instance of a game based on its type.
func NewGame(gameType string) (Game, error) {
	factory, ok := gameFactories[gameType]
	if !ok {
		return nil, fmt.Errorf("unsupported game: %s", gameType)
	}
	return factory()
}
