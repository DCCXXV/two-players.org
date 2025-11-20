<script lang="ts">
	import { getGameConfig } from '$lib/config/games';
	import GameRoom from '$lib/components/game/GameRoom.svelte';
	import Board from '$lib/components/tictactoe/Board.svelte';
	import GameStatus from '$lib/components/tictactoe/GameStatus.svelte';
	import { sendWebSocketMessage } from '$lib/socketStore';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const gameConfig = getGameConfig('tic-tac-toe');

	function onMove(cellIndex: number) {
		sendWebSocketMessage({
			type: 'make_move',
			payload: {
				cellIndex: cellIndex
			}
		});
	}
</script>

<GameRoom {gameConfig} room={data.room} error={data.error}>
	{#snippet boardComponent({ gameState, myTurn, disabled })}
		{#if gameState.players.length == 2}
			<Board board={gameState.game.board} {disabled} {onMove} />
		{:else}
			<Board board={gameState.game.board} disabled={true} {onMove} />
		{/if}
	{/snippet}

	{#snippet gameStatusComponent({ gameState, myTurn })}
		<GameStatus {gameState} {myTurn} />
	{/snippet}
</GameRoom>
