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
			{gameState.game.winner === 'draw'
				? "It's a draw!"
				: `The winner is ${gameState.game.winner}!`}
		{:else}
			{gameState.game.currentTurn === 0 ? 'P1' : 'P2'}'s turn
		{/if}
	</h3>
</div>
