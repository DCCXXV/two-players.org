<script lang="ts">
	import { getAllGameConfigs } from '$lib/config/games';

	const games = getAllGameConfigs();
</script>

<div class="mx-auto max-w-6xl">
	<section class="mb-12 text-center">
		<h1 class="lora-700 mb-4 text-4xl md:text-5xl">
			<span class="text-lime-400">two</span><span class="text-stone-500">-</span><span
				class="text-rose-400">players</span><span class="text-stone-500">.org</span
			>
		</h1>
		<p class="mb-6 text-2xl text-stone-400">
			Play combinatorial games vs other players in your browser. No signup.
		</p>
	</section>

	<section class="mb-12">
		<h2 class="mb-4 text-2xl font-bold text-stone-200">Games</h2>
		<div class="grid gap-4 sm:grid-cols-2">
			{#each games.slice(0, 3) as game}
				<a
					href="/play/{game.path}"
					class="rounded-0 group block border-b-1 border-stone-700 bg-stone-900 p-4 transition-all hover:border-lime-400"
				>
					<h3 class="mb-2 text-xl font-bold text-stone-200 group-hover:text-lime-400">
						{game.displayName}
					</h3>
					{#if game.gameplayGif}
						<div class="mb-2 overflow-hidden border-1 border-stone-700">
							<img src={game.gameplayGif} alt="{game.displayName} gameplay" class="h-auto w-full" />
						</div>
					{/if}
				</a>
			{/each}

			<div
				class="rounded-0 flex items-center justify-center border-b-1 border-stone-700 bg-stone-900 p-4"
			>
				<p class="text-stone-500">More coming soon</p>
			</div>
		</div>
	</section>

	<section class="mb-12">
		<h2 class="mb-4 text-2xl font-bold text-stone-200">About</h2>
		<div class="border-b-1 border-stone-700 bg-stone-900 p-6">
			<p class="mb-3 text-stone-400">
				Free, open-source platform for playing combinatorial games online. Real-time multiplayer, no downloads or accounts.
			</p>
			<p class="text-stone-400">
				<a
					href="https://github.com/DCCXXV/two-players.org"
					target="_blank"
					rel="noopener noreferrer"
					class="text-lime-400 hover:underline"
				>
					Open source on GitHub
				</a>
			</p>
		</div>
	</section>
</div>
