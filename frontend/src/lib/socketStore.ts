import { writable, type Writable } from 'svelte/store';

interface WebSocketMessage {
	type: string; // e.g: "connection_ready", "game_update", "error", "create_room", "join_room"
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	payload?: any;
}

const httpBackendUrl: string = import.meta.env.VITE_SOCKET_URL || 'http://localhost:8080';
const wsBackendUrl = httpBackendUrl.replace(/^http/, 'ws');

// Stores
export const socket: Writable<WebSocket | null> = writable(null);
export const isConnected: Writable<boolean> = writable(false);
export const displayName: Writable<string> = writable('');

let socketInstance: WebSocket | null = null;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const reconnectInterval = 3000; // ms

export function connectWebSocket(): void {
	if (socketInstance && socketInstance.readyState < 2) { // 0=CONNECTING, 1=OPEN
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
			console.log(`WebSocket closed. Attempting reconnect ${reconnectAttempts}/${maxReconnectAttempts} in ${reconnectInterval / 1000}s...`);
			setTimeout(connectWebSocket, reconnectInterval);
		} else {
			console.error('WebSocket reconnection attempts exhausted.');
		}
	};

	socketInstance.onerror = (event) => {
		console.error('WebSocket error:', event);
	};

	socketInstance.onmessage = (event) => {
		console.log('WebSocket message received:', event.data);
		try {
			const message: WebSocketMessage = JSON.parse(event.data);

			switch (message.type) {
				case 'connection_ready':
					if (message.payload && typeof message.payload === 'object' && 'displayName' in message.payload) {
						const name = message.payload.displayName as string;
						console.log('Received connection_ready, displayName:', name);
						displayName.set(name);
					} else {
						console.warn('Received connection_ready with invalid payload:', message.payload);
					}
					break;
				// case 'game_update':
				//  break;
				// case 'player_joined':
				//  break;
				// case 'error_message':
				//  break;
				default:
					console.warn('Received unknown message type:', message.type);
			}
		} catch (e) {
			console.error('Failed to parse WebSocket message or invalid message format:', e, event.data);
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
