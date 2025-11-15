export interface GameConfig {
	id: string; // normalized lowercase with hyphens
	displayName: string;
	path: string;
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
		sounds: {
			move1: '/sounds/moveX.wav',
			move2: '/sounds/moveO.wav',
			gameOver: '/sounds/gameOver.wav'
		},
		playerSymbols: ['X', 'O']
	},
	domineering: {
		id: 'domineering',
		displayName: 'Domineering',
		path: 'domineering',
		sounds: {
			move1: '/sounds/moveX.wav',
			move2: '/sounds/moveO.wav',
			gameOver: '/sounds/gameOver.wav'
		},
		playerSymbols: ['H', 'V']
	},
	'dots-and-boxes': {
		id: 'dots-and-boxes',
		displayName: 'Dots & Boxes',
		path: 'dots-and-boxes',
		sounds: {
			move1: '/sounds/moveX.wav',
			move2: '/sounds/moveO.wav',
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
