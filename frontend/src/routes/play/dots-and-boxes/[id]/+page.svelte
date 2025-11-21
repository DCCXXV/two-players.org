<script lang="ts">
	import { getGameConfig } from '$lib/config/games';
	import GameRoom from '$lib/components/game/GameRoom.svelte';
	import Board from '$lib/components/dotsandboxes/Board.svelte';
	import GameStatus from '$lib/components/dotsandboxes/GameStatus.svelte';
	import { sendWebSocketMessage, gameState, displayName } from '$lib/socketStore';
	import { SvelteSeo } from 'svelte-seo';
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

{#if data.meta}
	<SvelteSeo
		title={data.meta.title}
		description={data.meta.description}
		openGraph={{
			type: 'website',
			url: data.meta.url,
			title: data.meta.title,
			description: data.meta.description,
			images: [
				{
					url: data.meta.imageUrl,
					alt: data.meta.title
				}
			]
		}}
		twitter={{
			card: 'summary_large_image',
			site: data.meta.url,
			title: data.meta.title,
			description: data.meta.description,
			image: data.meta.imageUrl
		}}
	/>
{/if}

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
