<script lang="ts">
	let {
		gameState,
		myTurn
	}: {
		gameState: any;
		myTurn: boolean;
	} = $props();
</script>

<div class="mt-8 text-center text-zinc-400">
	{#if !gameState || !gameState.game || !gameState.players}
		<h3 class="text-3xl">Loading...</h3>
	{:else if gameState.game.winner}
		{#if gameState.game.winner === 'draw'}
			<h3 class="text-3xl">It's a draw!</h3>
		{:else if gameState.game.winner === 'P1'}
			<h3 class="text-3xl">
				{gameState.players[0]} wins as '<span class="text-blue-400">P1</span>'!
			</h3>
		{:else if gameState.game.winner === 'P2'}
			<h3 class="text-3xl">
				{gameState.players[1]} wins as '<span class="text-red-400">P2</span>'!
			</h3>
		{/if}
	{:else if gameState.players.length < 2}
		<h3 class="text-3xl">Waiting for another player to join...</h3>
	{:else}
		<h3 class="text-3xl">
			{gameState.players[gameState.game.currentTurn]}'s turn
			{#if gameState.game.boxesCompleted > 0}
				<span class="text-green-400">+{gameState.game.boxesCompleted}</span>
			{/if}
		</h3>

		<div class="mt-4 flex justify-center gap-8 text-xl">
			<div>
				<span class="text-blue-400">{gameState.players[0]}</span>: {gameState.game.scores[0]}
			</div>
			<div>
				<span class="text-red-400">{gameState.players[1]}</span>: {gameState.game.scores[1]}
			</div>
		</div>
	{/if}
</div>
