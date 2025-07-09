<script lang="ts">
	export let connections: { display_name: string; status: string; game_type: string | null }[];
</script>

<details class="bg-surface-900 border-surface-400 mb-4 h-full w-full border-2 lg:w-1/5">
	<summary
		class="text-surface-200 lora-700 hover:bg-surface-950 cursor-pointer p-2 text-lg font-bold"
	>
		Active Players: {connections.length}
	</summary>

	{#if connections && connections.length > 0}
		<ul class="list-none p-2">
			{#each connections as connection}
				<li class="flex items-center space-x-2">
					<div class="flex-1">
						<p class="text-primary-400">{connection.display_name}</p>
						<p class="text-surface-400 text-sm">
							Status: {connection.status}
							{#if connection.game_type}
								· Game: {connection.game_type}
							{/if}
						</p>
					</div>
				</li>
			{/each}
		</ul>
	{:else}
		<p class="text-surface-400 mt-2">No active players found.</p>
	{/if}
</details>

<style>
	details summary {
		list-style: none;
	}
	details summary::-webkit-details-marker {
		display: none;
	}
</style>
