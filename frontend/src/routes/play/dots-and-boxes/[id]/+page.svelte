<script lang="ts">
	import { getGameConfig } from '$lib/config/games';
	import GameRoom from '$lib/components/game/GameRoom.svelte';
	import Board from '$lib/components/dotsandboxes/Board.svelte';
	import GameStatus from '$lib/components/dotsandboxes/GameStatus.svelte';
	import { sendWebSocketMessage, gameState, displayName } from '$lib/socketStore';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const gameConfig = getGameConfig('dots-and-boxes');

	let mySymbol = $derived(
		$gameState?.players[0] === $displayName
			? 'P1'
			: $gameState?.players[1] === $displayName
				? 'P2'
				: undefined
	);

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

<svelte:head>
	{#if data.meta}
		<title>{data.meta.title}</title>
		<meta name="description" content={data.meta.description} />

		<meta property="og:type" content="website" />
		<meta property="og:url" content={data.meta.url} />
		<meta property="og:title" content={data.meta.title} />
		<meta property="og:description" content={data.meta.description} />
		<meta property="og:image" content={data.meta.imageUrl} />

		<meta property="twitter:card" content="summary_large_image" />
		<meta property="twitter:url" content={data.meta.url} />
		<meta property="twitter:title" content={data.meta.title} />
		<meta property="twitter:description" content={data.meta.description} />
		<meta property="twitter:image" content={data.meta.imageUrl} />
	{/if}
</svelte:head>

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
