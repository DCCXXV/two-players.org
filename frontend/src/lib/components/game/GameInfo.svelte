<script lang="ts">
	import type { GameConfig } from '$lib/config/games';
	import Collapsible from '$lib/components/ui/Collapsible.svelte';

	interface Props {
		gameConfig: GameConfig;
	}

	let { gameConfig }: Props = $props();
</script>

<aside class="flex h-full flex-col">
	<Collapsible title="About {gameConfig.displayName}">
		<div class="bg-stone-950 p-4 text-start">
			<p class="mb-4 text-stone-300">{gameConfig.description}</p>

			{#if gameConfig.gameplayGif}
				<div class="mb-4 overflow-hidden w-full">
					<img
						src={gameConfig.gameplayGif}
						alt="{gameConfig.displayName} gameplay"
						class="h-auto w-2/3 border-2 border-stone-700 mx-auto"
					/>
				</div>
			{/if}

			<div class="mt-4">
				<h3 class="mb-2 text-lg font-bold text-lime-400">Rules</h3>
				<ul class="space-y-2">
					{#each gameConfig.rules as rule, index}
						<li class="flex gap-2 text-sm text-stone-300">
							<span class="text-lime-400">{index + 1}.</span>
							<span>{rule}</span>
						</li>
					{/each}
				</ul>
			</div>
		</div>
	</Collapsible>
</aside>
