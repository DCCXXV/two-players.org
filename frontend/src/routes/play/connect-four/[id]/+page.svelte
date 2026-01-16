<script lang="ts">
	import { getGameConfig } from '$lib/config/games';
	import GameRoom from '$lib/components/game/GameRoom.svelte';
	import Board from '$lib/components/connectfour/Board.svelte';
	import GameStatus from '$lib/components/connectfour/GameStatus.svelte';
	import { sendWebSocketMessage } from '$lib/socketStore';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const gameConfig = getGameConfig('connect-four');

	function onMove(column: number) {
		sendWebSocketMessage({
			type: 'make_move',
			payload: {
				column: column
			}
		});
	}
</script>

<svelte:head>
	{#if data.meta}
		<title>{data.meta.title}</title>
		<meta name="description" content={data.meta.description} />
		<meta property="og:type" content="website" />
		<meta property="og:url" content={data.meta.url} />
		<meta property="og:title" content={data.meta.title} />
		<meta property="og:description" content={data.meta.description} />
		<meta property="og:image" content={data.meta.imageUrl} />
		<meta name="twitter:card" content="summary_large_image" />
		<meta name="twitter:url" content={data.meta.url} />
		<meta name="twitter:title" content={data.meta.title} />
		<meta name="twitter:description" content={data.meta.description} />
		<meta name="twitter:image" content={data.meta.imageUrl} />
	{/if}
</svelte:head>

<GameRoom {gameConfig} room={data.room} error={data.error}>
	{#snippet boardComponent({ gameState, myTurn, disabled })}
		{#if gameState.players.length == 2}
			<Board
				board={gameState.game.board}
				winningCells={gameState.game.winningCells}
				{disabled}
				{onMove}
			/>
		{:else}
			<Board
				board={gameState.game.board}
				winningCells={gameState.game.winningCells}
				disabled={true}
				{onMove}
			/>
		{/if}
	{/snippet}

	{#snippet gameStatusComponent({ gameState, myTurn })}
		<GameStatus {gameState} {myTurn} />
	{/snippet}
</GameRoom>
