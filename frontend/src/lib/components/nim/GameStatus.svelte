<script lang="ts">
	interface Props {
		gameState: any;
		myTurn: boolean;
	}

	let { gameState, myTurn }: Props = $props();
</script>

<div class="text-center text-3xl text-zinc-400">
	<h3>
		{#if !gameState}
			Loading...
		{:else if gameState.game.winner}
			{#if gameState.game.winner === 'P1'}
				{gameState.players[0]} wins as '<span class="text-blue-400">P1</span>'!
			{:else if gameState.game.winner === 'P2'}
				{gameState.players[1]} wins as '<span class="text-red-400">P2</span>'!
			{/if}
		{:else if gameState.players.length < 2}
			Waiting for another player to join...
		{:else}
			{gameState.players[gameState.game.currentTurn]}'s turn
		{/if}
	</h3>
</div>
