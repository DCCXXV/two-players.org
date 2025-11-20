<script lang="ts">
	import { getGameConfig } from '$lib/config/games';
	import GameRoom from '$lib/components/game/GameRoom.svelte';
	import Board from '$lib/components/nim/Board.svelte';
	import GameStatus from '$lib/components/nim/GameStatus.svelte';
	import { sendWebSocketMessage } from '$lib/socketStore';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const gameConfig = getGameConfig('nim');

	function onMove(sticksToTake: number) {
		sendWebSocketMessage({
			type: 'make_move',
			payload: {
				sticks: sticksToTake
			}
		});
	}
</script>

<svelte:head>
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
</svelte:head>

<GameRoom {gameConfig} room={data.room} error={data.error}>
	{#snippet boardComponent({ gameState, myTurn, disabled })}
		<Board sticks={gameState.game.sticks} {disabled} {onMove} />
	{/snippet}

	{#snippet gameStatusComponent({ gameState, myTurn })}
		<GameStatus {gameState} {myTurn} />
	{/snippet}
</GameRoom>
