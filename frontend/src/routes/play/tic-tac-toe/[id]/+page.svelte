<script lang="ts">
	import { getGameConfig } from '$lib/config/games';
	import GameRoom from '$lib/components/game/GameRoom.svelte';
	import Board from '$lib/components/tictactoe/Board.svelte';
	import GameStatus from '$lib/components/tictactoe/GameStatus.svelte';
	import { sendWebSocketMessage } from '$lib/socketStore';
	import { SvelteSeo } from 'svelte-seo';
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
			<Board board={gameState.game.board} {disabled} {onMove} />
		{:else}
			<Board board={gameState.game.board} disabled={true} {onMove} />
		{/if}
	{/snippet}

	{#snippet gameStatusComponent({ gameState, myTurn })}
		<GameStatus {gameState} {myTurn} />
	{/snippet}
</GameRoom>
