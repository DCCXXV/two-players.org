<script lang="ts">
	import { getAllGameConfigs } from '$lib/config/games';
	import { MoveRight } from 'lucide-svelte';

	const games = getAllGameConfigs();
</script>

<div class="mx-auto max-w-5xl">
	<section class="mb-8 text-center">
		<h1 class="lora-700 mb-3 text-4xl md:text-5xl">
			<span class="text-blue-400">two</span><span class="text-zinc-500">-</span><span
				class="text-red-400">players</span><span class="text-zinc-500">.org</span
			>
		</h1>
		<p class="mb-6 text-2xl text-zinc-400">
			Play combinatorial games vs other players in your browser. No signup.
		</p>
	</section>

	<section class="mb-8">
		<div class="mb-3 flex items-baseline justify-between border-b-1 border-zinc-700 pb-1">
			<h2 class="text-xl font-semibold text-zinc-200">Games</h2>
			<a href="/play" class="inline-flex items-center gap-1 text-sm text-blue-400 hover:underline">View all <MoveRight size={14}/></a>
		</div>
		<div class="grid gap-3 sm:grid-cols-3">
			{#each games.slice(0, 3) as game}
				<a
					href="/play/{game.path}"
					class="rounded-0 group block border-b-1 border-zinc-700 bg-zinc-800 p-3 transition-all hover:border-blue-400"
				>
					<h3 class="mb-2 text-lg font-semibold text-zinc-200 group-hover:text-blue-400">
						{game.displayName}
					</h3>
					{#if game.gameplayGif}
						<div class="overflow-hidden border-1 border-zinc-700">
							<img src={game.gameplayGif} alt="{game.displayName} gameplay" class="h-auto w-full" />
						</div>
					{/if}
				</a>
			{/each}
		</div>
	</section>

	<section class="mb-8">
		<div class="mb-3 border-b-1 border-zinc-700 pb-1">
			<h2 class="text-xl font-semibold text-zinc-200">About</h2>
		</div>
		<div class="border-b-1 border-zinc-700 bg-zinc-800 p-4">
			<p class="mb-2 text-zinc-400">
				Free, open-source platform for playing combinatorial games online. Real-time multiplayer, no downloads or accounts.
			</p>
			<a
				href="https://github.com/DCCXXV/two-players.org"
				target="_blank"
				rel="noopener noreferrer"
				class="inline-flex items-center gap-1 text-blue-400 hover:underline"
			>
				Open source on GitHub <MoveRight size={14}/>
			</a>
		</div>
	</section>
</div>
