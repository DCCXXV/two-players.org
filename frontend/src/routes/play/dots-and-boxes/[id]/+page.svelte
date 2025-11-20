<script lang="ts">
	import { getGameConfig } from '$lib/config/games';
	import GameRoom from '$lib/components/game/GameRoom.svelte';
	import Board from '$lib/components/dotsandboxes/Board.svelte';
	import GameStatus from '$lib/components/dotsandboxes/GameStatus.svelte';
	import { sendWebSocketMessage, gameState, displayName } from '$lib/socketStore';
	import type { PageData } from './$types';

	export let data: PageData;

	const gameConfig = getGameConfig('dots-and-boxes');

	$: mySymbol =
		$gameState?.players[0] === $displayName
			? 'P1'
			: $gameState?.players[1] === $displayName
				? 'P2'
				: undefined;

	function onMove(type: string, row: number, col: number) {
		sendWebSocketMessage({
			type: 'make_move',
			payload: {
				type: type,
				row: row,
				col: col
			}
		});
	}
</script>

<GameRoom {gameConfig} room={data.room} error={data.error}>
	{#snippet boardComponent({ gameState, myTurn, disabled })}
		{#if gameState.players.length == 2}
			<Board
				hLines={gameState.game.hLines}
				vLines={gameState.game.vLines}
				boxes={gameState.game.boxes}
				{disabled}
				{onMove}
				{mySymbol}
			/>
		{:else}
			<Board
				hLines={gameState.game.hLines}
				vLines={gameState.game.vLines}
				boxes={gameState.game.boxes}
				disabled={true}
				{onMove}
			/>
		{/if}
	{/snippet}

	{#snippet gameStatusComponent({ gameState, myTurn })}
		<GameStatus {gameState} {myTurn} />
	{/snippet}
</GameRoom>
