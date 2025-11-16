<script lang="ts">
	let {
		gameState,
		myTurn
	}: {
		gameState: any;
		myTurn: boolean;
	} = $props();
</script>

<div class="text-center text-3xl text-zinc-400">
	<h3>
		{#if !gameState || !gameState.game || !gameState.players}
			Loading...
		{:else if gameState.game.winner}
			{#if gameState.game.winner === 'H'}
				{gameState.players[0]} wins as '<span class="text-blue-400">H</span>'!
			{:else if gameState.game.winner === 'V'}
				{gameState.players[1]} wins as '<span class="text-red-400">V</span>'!
			{/if}
		{:else if gameState.players.length < 2}
			Waiting for another player to join...
		{:else}
			{gameState.players[gameState.game.currentTurn]}'s turn
		{/if}
	</h3>
</div>
