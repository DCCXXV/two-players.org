<script lang="ts">
	const imgSrc =
		'https://dummyimage.com/420x180/5B8139/1F1919&text=Tic-Tac-Toe';
</script>

<section class="w-full mt-4">
	<h1 class="lora-700 text-surface-200 h3 mb-6">Current games</h1>
	<section class="flex gap-4 items-stretch">
		<!-- Tic-Tac-Toe -->
		<div class="relative group max-w-xs">
			<a
				href="/play/tic-tac-toe"
				class="card flex flex-col h-full preset-filled-surface-100-900 border-primary-200-800 card-hover divide-surface-200-800 max-w-xs divide-y overflow-hidden border-[1px] rounded-none z-10"
			>
				<article class="space-y-4 p-4 flex-1">
				<div>
					<h2 class="h6 lora-700">Tic-Tac-Toe</h2>
				</div>
				<p class="opacity-60">
					Played on a 3x3 grid, one player uses 'X', the other 'O'. Players take turns marking an empty square. The first player to get three of their marks in a horizontal, vertical, or diagonal row wins. If the grid is full and no one has won, the game is a draw.
				</p>
				</article>
				<footer class="flex items-center gap-2 p-4">
					<span class="badge preset-filled-surface-500 rounded-none">Impartial</span>
					<span class="badge preset-filled-surface-500 rounded-none">Draw</span>
					<span class="badge preset-filled-surface-500 rounded-none">Strongly Solved</span>
				</footer>
			</a>
			<span
			class="pointer-events-none absolute w-full left-0 top-full -mt-2 px-2 py-1 bg-primary-400 text-white text-xs opacity-0 group-hover:opacity-100 transition-all duration-300 whitespace-nowrap z-0 text-center translate-y-0 group-hover:translate-y-2"
			role="tooltip"
			>
			FINISHED
			</span>
		</div>

		<!-- Domineering -->
		<div class="relative group max-w-xs">
			<a
				href="#"
				class="card flex flex-col h-full preset-filled-surface-100-900 border-secondary-200-800 card-hover divide-surface-200-800 max-w-xs divide-y overflow-hidden border-[1px] rounded-none z-10"
			>
				<article class="space-y-4 p-4 flex-1">
				<div>
					<h2 class="h6 lora-700">Domineering</h2>
				</div>
				<p class="opacity-60">
					Two players take turns placing dominoes on a grid: one places them vertically, the other horizontally. Each domino covers two adjacent squares. The first player unable to move loses. No draws are possible.
				</p>
				</article>
				<footer class="flex items-center gap-2 p-4">
					<span class="badge preset-filled-surface-500 rounded-none">Partisan</span>
					<span class="badge preset-filled-surface-500 rounded-none">No Draw</span>
					<span class="badge preset-filled-surface-500 rounded-none">Solved</span>
				</footer>
			</a>
			<span
			class="pointer-events-none absolute w-full left-0 top-full -mt-2 px-2 py-1 bg-secondary-400 text-white text-xs opacity-0 group-hover:opacity-100 transition-all duration-300 whitespace-nowrap z-0 text-center translate-y-0 group-hover:translate-y-2"
			role="tooltip"
			>
			PENDING
			</span>
		</div>

		<!-- Dots and Boxes -->
		<div class="relative group max-w-xs">
			<a
				href="#"
				class="card flex flex-col h-full preset-filled-surface-100-900 border-secondary-200-800 card-hover divide-surface-200-800 max-w-xs divide-y overflow-hidden border-[1px] rounded-none z-10"
			>
				<article class="space-y-4 p-4 flex-1">
				<div>
					<h2 class="h6 lora-700">Dots and Boxes</h2>
				</div>
				<p class="opacity-60">
					Players take turns drawing lines between adjacent dots on a grid. Completing the fourth side of a box claims it and earns another turn. The player with the most boxes at the end wins.
				</p>
				</article>
				<footer class="flex items-center gap-2 p-4">
					<span class="badge preset-filled-surface-500 rounded-none">Partisan</span>
					<span class="badge preset-filled-surface-500 rounded-none">No Draw</span>
					<span class="badge preset-filled-surface-500 rounded-none">Solved</span>
				</footer>
			</a>
			<span
			class="pointer-events-none absolute w-full left-0 top-full -mt-2 px-2 py-1 bg-secondary-400 text-white text-xs opacity-0 group-hover:opacity-100 transition-all duration-300 whitespace-nowrap z-0 text-center translate-y-0 group-hover:translate-y-2"
			role="tooltip"
			>
			PENDING
			</span>
		</div>
	</section>

</section>
