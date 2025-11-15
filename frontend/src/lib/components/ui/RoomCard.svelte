<script lang="ts">
	interface Room {
		id: string;
		name: string;
		game_type: string;
		is_private: boolean;
		created_by?: string | null;
		other_player?: string | null;
	}

	let { room, showGameType = true }: { room: Room; showGameType?: boolean } = $props();

	const gameConfig: Record<string, { color: string; displayName: string; path: string }> = {
		'tic-tac-toe': { displayName: 'Tic-Tac-Toe', path: 'tic-tac-toe' },
		domineering: { displayName: 'Domineering', path: 'domineering' },
		'dots-and-boxes': {
			displayName: 'Dots & Boxes',
			path: 'dots-and-boxes'
		}
	};

	const config = gameConfig[room.game_type] || {
		displayName: room.game_type,
		path: room.game_type
	};
</script>

<a href={`/play/${config.path}/${room.id}`} class="group">
	<div
		class="flex flex-col gap-2 border-b-1 border-stone-700 bg-stone-950 p-4 text-start transition-colors hover:border-lime-400"
	>
		{#if showGameType}
			<div class="mb-1 flex justify-between">
				<span class="text-sm font-bold text-stone-300">{config.displayName}</span>
			</div>
		{/if}
		<h4 class="text-lg text-pretty text-rose-400">{room.name}</h4>
		<div class="space-y-1 text-sm text-stone-400">
			<div>Player 1: <span class="font-bold">{room.created_by}</span></div>
			<div>
				Player 2:
				{#if room.other_player}
					<span class="font-bold">{room.other_player}</span>
				{:else}
					<span class="text-stone-500 italic">Join to play!</span>
				{/if}
			</div>
		</div>
		<div class="mt-auto">
			<div class="text-xs text-stone-400">Click to join →</div>
		</div>
	</div>
</a>
