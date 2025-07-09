<script lang="ts">
	import ActiveConnectionsSidebar from '$lib/components/ActiveConnectionsSidebar.svelte';
	export let data;
</script>

<div class="flex flex-col gap-4 lg:flex-row">
	<section class="mt-4 w-full lg:w-4/5">
		<h1 class="lora-700 text-surface-200 h3 mb-6">Current games</h1>
		<section class="flex flex-wrap items-stretch gap-4">
			<!-- Tic-Tac-Toe -->
			<div class="group relative max-w-xs">
				<a
					href="/play/tic-tac-toe"
					class="card preset-filled-surface-100-900 border-warning-200-800 card-hover divide-surface-200-800 z-10 flex h-full max-w-xs flex-col divide-y overflow-hidden rounded-none border-[1px]"
				>
					<article class="flex-1 space-y-4 p-4">
						<div>
							<h2 class="h6 lora-700">Tic-Tac-Toe</h2>
						</div>
						<p class="opacity-60">
							Played on a 3x3 grid, one player uses 'X', the other 'O'. Players take turns marking
							an empty square. The first player to get three of their marks in a horizontal,
							vertical, or diagonal row wins. If the grid is full and no one has won, the game is a
							draw.
						</p>
					</article>
					<footer class="flex items-center gap-2 p-4">
						<span class="badge preset-filled-surface-500 rounded-none">Impartial</span>
						<span class="badge preset-filled-surface-500 rounded-none">Draw</span>
						<span class="badge preset-filled-surface-500 rounded-none">Strongly Solved</span>
					</footer>
				</a>
				<span
					class="bg-warning-400 pointer-events-none absolute top-full left-0 z-0 -mt-2 w-full translate-y-0 px-2 py-1 text-center text-xs whitespace-nowrap text-black opacity-0 transition-all duration-300 group-hover:translate-y-2 group-hover:opacity-100"
					role="tooltip"
				>
					WIP
				</span>
			</div>

			<!-- Domineering -->
			<div class="group relative max-w-xs">
				<a
					href="#"
					class="card preset-filled-surface-100-900 border-secondary-200-800 card-hover divide-surface-200-800 z-10 flex h-full max-w-xs flex-col divide-y overflow-hidden rounded-none border-[1px]"
				>
					<article class="flex-1 space-y-4 p-4">
						<div>
							<h2 class="h6 lora-700">Domineering</h2>
						</div>
						<p class="opacity-60">
							Two players take turns placing dominoes on a grid: one places them vertically, the
							other horizontally. Each domino covers two adjacent squares. The first player unable
							to move loses. No draws are possible.
						</p>
					</article>
					<footer class="flex items-center gap-2 p-4">
						<span class="badge preset-filled-surface-500 rounded-none">Partisan</span>
						<span class="badge preset-filled-surface-500 rounded-none">No Draw</span>
						<span class="badge preset-filled-surface-500 rounded-none">Solved</span>
					</footer>
				</a>
				<span
					class="bg-secondary-400 pointer-events-none absolute top-full left-0 z-0 -mt-2 w-full translate-y-0 px-2 py-1 text-center text-xs whitespace-nowrap text-white opacity-0 transition-all duration-300 group-hover:translate-y-2 group-hover:opacity-100"
					role="tooltip"
				>
					PENDING
				</span>
			</div>

			<!-- Dots and Boxes -->
			<div class="group relative max-w-xs">
				<a
					href="#"
					class="card preset-filled-surface-100-900 border-secondary-200-800 card-hover divide-surface-200-800 z-10 flex h-full max-w-xs flex-col divide-y overflow-hidden rounded-none border-[1px]"
				>
					<article class="flex-1 space-y-4 p-4">
						<div>
							<h2 class="h6 lora-700">Dots and Boxes</h2>
						</div>
						<p class="opacity-60">
							Players take turns drawing lines between adjacent dots on a grid: Completing the
							fourth side of a box claims it and earns another turn. The player with the most boxes
							at the end wins.
						</p>
					</article>
					<footer class="flex items-center gap-2 p-4">
						<span class="badge preset-filled-surface-500 rounded-none">Partisan</span>
						<span class="badge preset-filled-surface-500 rounded-none">No Draw</span>
						<span class="badge preset-filled-surface-500 rounded-none">Solved</span>
					</footer>
				</a>
				<span
					class="bg-secondary-400 pointer-events-none absolute top-full left-0 z-0 -mt-2 w-full translate-y-0 px-2 py-1 text-center text-xs whitespace-nowrap text-white opacity-0 transition-all duration-300 group-hover:translate-y-2 group-hover:opacity-100"
					role="tooltip"
				>
					PENDING
				</span>
			</div>
		</section>
	</section>

	<ActiveConnectionsSidebar connections={data.connections} />
</div>
