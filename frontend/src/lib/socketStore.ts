import { writable, type Writable } from 'svelte/store';

export interface TicTacToeGameState {
	board: ('' | 'X' | 'O')[];
	winner: string;
	currentTurn: number;
}

export interface GameState {
	game: TicTacToeGameState;
	players: string[];
	playerCount: number;
	maxPlayers: number;
	spectatorCount: number;
	spectators: string[];
	canStart: boolean;
	rematchCount: number;
}

interface ConnectionReadyPayload {
	displayName: string;
}

interface GameStateUpdatePayload extends GameState {}

interface JoinSuccessPayload {
	message?: string;
}

interface RoomClosedPayload {
	message: string;
}

interface ErrorPayload {
	message: string;
}

export interface ChatMessage {
	displayName: string;
	message: string;
	timestamp: string;
}

interface ChatMessagePayload extends ChatMessage {}

export type WebSocketMessage =
	| { type: 'connection_ready'; payload: ConnectionReadyPayload }
	| { type: 'game_state_update'; payload: GameStateUpdatePayload }
	| { type: 'join_success'; payload?: JoinSuccessPayload }
	| { type: 'room_closed'; payload: RoomClosedPayload }
	| { type: 'error'; payload: ErrorPayload }
	| { type: 'chat_message'; payload: ChatMessagePayload }
	| { type: string; payload?: any };

const httpBackendUrl: string = import.meta.env.VITE_SOCKET_URL || 'http://localhost:8080';
const wsBackendUrl = httpBackendUrl.replace(/^http/, 'ws');

export const socket: Writable<WebSocket | null> = writable(null);
export const isConnected: Writable<boolean> = writable(false);
export const displayName: Writable<string> = writable('');
export const gameState: Writable<GameState | null> = writable(null);
export const players: Writable<string[]> = writable([]);
export const errorMessage: Writable<string | null> = writable(null);
export const roomClosedMessage: Writable<string | null> = writable(null);
export const chatMessages: Writable<ChatMessage[]> = writable([]);

let socketInstance: WebSocket | null = null;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const reconnectInterval = 3000;
export function connectWebSocket(): void {
	if (socketInstance && socketInstance.readyState < 2) {
		console.log('WebSocket already connected or connecting.');
		return;
	}
	console.log(`Attempting to connect WebSocket to: ${wsBackendUrl}/ws`);

	socketInstance = new WebSocket(`${wsBackendUrl}/ws`);

	socket.set(socketInstance);
	isConnected.set(false);

	socketInstance.onopen = (event) => {
		console.log('WebSocket connected:', event);
		isConnected.set(true);
		reconnectAttempts = 0;
	};

	socketInstance.onclose = (event) => {
		console.log('WebSocket disconnected:', event.code, event.reason);
		isConnected.set(false);
		displayName.set('');
		socket.set(null);
		socketInstance = null;

		if (reconnectAttempts < maxReconnectAttempts) {
			reconnectAttempts++;
			console.log(
				`WebSocket closed. Attempting reconnect ${reconnectAttempts}/${maxReconnectAttempts} in ${
					reconnectInterval / 1000
				}s...`
			);
			setTimeout(connectWebSocket, reconnectInterval);
		} else {
			console.error('WebSocket reconnection attempts exhausted.');
		}
	};

	socketInstance.onerror = (event) => {
		console.error('WebSocket error:', event);
	};

	socketInstance.onmessage = (event) => {
		console.log('WebSocket RAW message received:', event.data);
		try {
			const message: WebSocketMessage = JSON.parse(event.data);
			console.log('Parsed message - Type:', message.type, 'Payload:', message.payload);

			switch (message.type) {
				case 'connection_ready':
					if (
						message.payload &&
						typeof message.payload === 'object' &&
						'displayName' in message.payload
					) {
						const name = message.payload.displayName as string;
						console.log('Received connection_ready, displayName:', name);
						displayName.set(name);
					} else {
						console.warn('Received connection_ready with invalid payload:', message.payload);
					}
					break;

				case 'game_state_update':
					console.log('Game state update received:', message.payload);
					gameState.set(message.payload);
					if (message.payload && Array.isArray(message.payload.players)) {
						console.log('Updating players:', message.payload.players);
						players.set(message.payload.players);
					}
					break;

				case 'join_success':
					console.log('Successfully joined room:', message.payload);
					break;

				case 'room_closed':
					console.log('Room closed by host:', message.payload);
					roomClosedMessage.set(message.payload.message || 'The host has left the room.');
					gameState.set(null);
					players.set([]);
					break;

				case 'error':
					console.log('Error received:', message.payload);
					if (message.payload.message === 'You are already in another room.') {
						const currentGame = gameState.subscribe((value) => value)();
						if (currentGame && !currentGame.game?.winner) {
							alert('You are already in an active game. Please finish or leave that game first.');
							window.history.back();
						}
					} else {
						errorMessage.set(message.payload.message);
					}
					break;

				case 'chat_message':
					console.log('Chat message received:', message.payload);
					chatMessages.update((messages) => [...messages, message.payload]);
					break;

				default:
					console.warn('Unknown message type:', message.type, 'Full message:', message);
			}
		} catch (e) {
			console.error('Failed to parse WebSocket message:', e, event.data);
		}
	};
}
export function disconnectWebSocket(): void {
	if (socketInstance) {
		reconnectAttempts = maxReconnectAttempts;
		socketInstance.close();
	}
}
export function sendWebSocketMessage(message: WebSocketMessage): void {
	if (socketInstance && socketInstance.readyState === WebSocket.OPEN) {
		try {
			const jsonString = JSON.stringify(message);
			socketInstance.send(jsonString);
			console.log('WebSocket message sent:', jsonString);
		} catch (e) {
			console.error('Failed to stringify or send WebSocket message:', e, message);
		}
	} else {
		console.error('WebSocket not connected, cannot send message:', message.type);
	}
}

export function updateDisplayName(newName: string): void {
	if (socketInstance && socketInstance.readyState === WebSocket.OPEN) {
		try {
			const message = {
				type: 'update_display_name',
				payload: { displayName: newName }
			};
			const jsonString = JSON.stringify(message);
			socketInstance.send(jsonString);
			console.log('Display name update sent:', jsonString);
		} catch (e) {
			console.error('Failed to send display name update:', e);
		}
	} else {
		console.error('WebSocket not connected, cannot update display name');
	}
}
