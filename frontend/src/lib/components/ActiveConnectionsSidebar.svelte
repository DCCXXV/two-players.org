<script lang="ts">
	export let connections: { display_name: string; status: string; game_type: string | null }[];
</script>

<aside class="border-surface-400 w-full space-y-4 border-[1px] p-4 lg:w-1/5">
	<h2 class="h5 lora-700 text-surface-200">Active Players: {connections.length}</h2>
	<hr class="text-primary-400" />
	{#if connections && connections.length > 0}
		<ul class="m-0 list-none space-y-2 p-0">
			{#each connections as connection}
				<li class="flex items-center space-x-2">
					<div class="flex-1">
						<p class="text-surface-200 font-bold">{connection.display_name}</p>
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
		<p class="text-surface-400">No active players found.</p>
	{/if}
</aside>
