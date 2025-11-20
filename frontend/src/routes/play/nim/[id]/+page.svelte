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

<GameRoom {gameConfig} room={data.room} error={data.error}>
	{#snippet boardComponent({ gameState, myTurn, disabled })}
		<Board sticks={gameState.game.sticks} {disabled} {onMove} />
	{/snippet}

	{#snippet gameStatusComponent({ gameState, myTurn })}
		<GameStatus {gameState} {myTurn} />
	{/snippet}
</GameRoom>
