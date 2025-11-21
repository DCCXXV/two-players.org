export interface GameConfig {
	id: string; // normalized lowercase with hyphens
	displayName: string;
	path: string;
	description: string;
	rules: string[];
	gameplayGif?: string;
	sounds: {
		move1: string;
		move2: string;
		gameOver: string;
	};
	playerSymbols: [string, string];
}

export const GAME_CONFIGS: Record<string, GameConfig> = {
	'tic-tac-toe': {
		id: 'tic-tac-toe',
		displayName: 'Tic-Tac-Toe',
		path: 'tic-tac-toe',
		description:
			"The classic 3x3 grid game where players alternate placing X's and O's, aiming to get three in a row.",
		rules: [
			'You must place your symbol in an empty square',
			'You must make a move on your turn',
			'Three in a row wins (horizontal, vertical, or diagonal)',
			'Game ends in a draw if the board is full with no winner'
		],
		gameplayGif: '/img/tic-tac-toe.gif',
		sounds: {
			move1: '/sounds/move1.wav',
			move2: '/sounds/move2.wav',
			gameOver: '/sounds/gameOver.wav'
		},
		playerSymbols: ['X', 'O']
	},
	domineering: {
		id: 'domineering',
		displayName: 'Domineering',
		path: 'domineering',
		description:
			'A strategic blocking game where one player places horizontal dominoes and the other vertical, competing for space on the grid. The player who can no longer make a valid move loses.',
		rules: [
			'Horizontal player can only place dominoes horizontally',
			'Vertical player can only place dominoes vertically',
			'Dominoes must occupy exactly 2 adjacent empty squares',
			'You lose if you cannot make a valid move on your turn'
		],
		gameplayGif: '/img/domineering.gif',
		sounds: {
			move1: '/sounds/move1.wav',
			move2: '/sounds/move2.wav',
			gameOver: '/sounds/gameOver.wav'
		},
		playerSymbols: ['H', 'V']
	},
	'dots-and-boxes': {
		id: 'dots-and-boxes',
		displayName: 'Dots & Boxes',
		path: 'dots-and-boxes',
		description:
			'Connect dots to form boxes and claim them. The player who completes the most boxes wins!',
		rules: [
			'You must draw a line between two adjacent dots',
			'You cannot draw a line that already exists',
			'Completing a box earns you another turn',
			'Most boxes claimed wins the game'
		],
		gameplayGif: '/img/dots-and-boxes.gif',
		sounds: {
			move1: '/sounds/move1.wav',
			move2: '/sounds/move2.wav',
			gameOver: '/sounds/gameOver.wav'
		},
		playerSymbols: ['P1', 'P2']
	},
	nim: {
		id: 'nim',
		displayName: 'Nim',
		path: 'nim',
		description:
			'A strategic game starting with 21 sticks. Players take turns removing 1, 2, or 3 sticks. Force your opponent to take the last stick and win!',
		rules: [
			'You must remove 1, 2, or 3 sticks on your turn',
			'You cannot skip your turn',
			'You cannot remove more than 3 sticks',
			'You lose if you take the last stick'
		],
		gameplayGif: '/img/nim.gif',
		sounds: {
			move1: '/sounds/move1.wav',
			move2: '/sounds/move2.wav',
			gameOver: '/sounds/gameOver.wav'
		},
		playerSymbols: ['P1', 'P2']
	}
};

export function getGameConfig(gameType: string): GameConfig {
	const normalized = gameType.toLowerCase();
	const config = GAME_CONFIGS[normalized];
	if (!config) {
		throw new Error(`Unknown game type: ${gameType}`);
	}
	return config;
}

export function getAllGameConfigs(): GameConfig[] {
	return Object.values(GAME_CONFIGS);
}
