<script lang="ts">
	import type { GameState } from '$lib/socketStore';

	let { gameState, myTurn } = $props<{
		gameState: GameState;
		myTurn: boolean;
	}>();
</script>

<div class="text-center text-3xl text-stone-400">
	<h3>
		{#if !gameState}
			Cargando...
		{:else if gameState.game.winner}
			{#if gameState.game.winner === 'draw'}
				It's a draw!
			{:else if gameState.game.winner === 'X'}
				{gameState.players[0]} wins as '<span class="text-rose-400">X</span>'!
			{:else if gameState.game.winner === 'O'}
				{gameState.players[1]} wins as '<span class="text-lime-400">O</span>'!
			{/if}
		{:else}
			{gameState.game.currentTurn === 0 ? `${gameState.players[0]}` : `${gameState.players[1]}`}'s
			turn
		{/if}
	</h3>
</div>
