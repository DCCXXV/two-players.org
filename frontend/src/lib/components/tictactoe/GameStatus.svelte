<script lang="ts">
	let { gameState, myTurn } = $props<{
		gameState: any;
		myTurn: boolean;
	}>();
</script>

<div class="text-surface-200 lora-700 text-center text-3xl">
	<h3>
		{#if !gameState}
			Cargando...
		{:else if gameState.game.winner}
			{#if gameState.game.winner === 'draw'}
				It's a draw!
			{:else if gameState.game.winner === 'X'}
				{gameState.players[0]} wins as '<span class="text-primary-300 erasdust-400">X</span>'!
			{:else if gameState.game.winner === 'O'}
				{gameState.players[1]} wins as '<span class="text-secondary-300 erasdust-400">O</span>'!
			{/if}
		{:else}
			{gameState.game.currentTurn === 0 ? `${gameState.players[0]}` : `${gameState.players[1]}`}'s
			turn
		{/if}
	</h3>
</div>
