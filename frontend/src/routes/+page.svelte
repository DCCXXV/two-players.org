<script lang="ts">
	const imgSrc =
		'https://dummyimage.com/420x180/5B8139/1F1919&text=Tic-Tac-Toe';
</script>

<section class="w-full mt-4">
	<h1 class="lora-700 text-surface-200 h1 mb-6">Current games</h1>

	<a
		href="/play/tic-tac-toe"
        class="card preset-filled-surface-100-900 border-surface-200-800 card-hover divide-surface-200-800 block max-w-sm divide-y overflow-hidden border-[1px] rounded-md"
    >
		<header>
			<img src={imgSrc} class="aspect-[21/9] w-full" alt="banner" />
		</header>
		<article class="space-y-4 p-4">
			<div>
				<h2 class="h6">Tic-Tac-Toe</h2>
			</div>
			<p class="opacity-60">
                Played on a 3x3 grid, one player uses 'X', the other 'O'. Players take turns marking an empty square. The first player to get three of their marks in a horizontal, vertical, or diagonal row wins. If the grid is full and no one has won, the game is a draw.
			</p>
		</article>
		<footer class="flex items-center gap-2 p-4">
			<span class="badge preset-filled-surface-500">Impartial</span>
            <span class="badge preset-filled-surface-500">Draw</span>
            <span class="badge preset-filled-surface-500">Strongly Solved</span>
		</footer>
	</a>
</section>
