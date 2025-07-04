package games

import "fmt"

// Game define la interfaz que cada juego debe implementar.
type Game interface {
	// HandleMove procesa un movimiento de un jugador.
	// El formato de 'move' es específico para cada juego.
	HandleMove(playerIndex int, move any) error

	// GetGameState devuelve el estado actual del juego para ser enviado a los clientes.
	GetGameState() any

	// IsGameOver comprueba si el juego ha terminado.
	IsGameOver() bool

	// GetWinner devuelve el identificador del ganador, si lo hay.
	GetWinner() string

	// Reset reinicia el juego a su estado inicial.
	Reset()
}

// Factory es una función que crea una nueva instancia de un juego.
type Factory func() (Game, error)

var gameFactories = make(map[string]Factory)

// RegisterGame se usa para registrar un nuevo tipo de juego.
func RegisterGame(gameType string, factory Factory) {
	if factory == nil {
		panic("La factoría de juegos no puede ser nula")
	}
	_, registered := gameFactories[gameType]
	if registered {
		// Opcional: podrías querer que esto sea un panic si ocurre durante el init.
		// Por ahora, lo dejamos como un log o una advertencia silenciosa.
		return
	}
	gameFactories[gameType] = factory
}

// NewGame crea una instancia de un juego basado en su tipo.
func NewGame(gameType string) (Game, error) {
	factory, ok := gameFactories[gameType]
	if !ok {
		return nil, fmt.Errorf("juego no soportado: %s", gameType)
	}
	return factory()
}
