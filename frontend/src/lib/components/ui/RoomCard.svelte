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
		class="flex flex-col gap-2 border-b-1 border-zinc-700 bg-zinc-900 p-4 text-start transition-colors hover:border-blue-400"
	>
		{#if showGameType}
			<div class="mb-1 flex justify-between">
				<span class="text-sm font-bold text-zinc-300">{config.displayName}</span>
			</div>
		{/if}
		<h4 class="text-lg text-pretty text-red-400">{room.name}</h4>
		<div class="space-y-1 text-sm text-zinc-400">
			<div>Player 1: <span class="font-bold">{room.created_by}</span></div>
			<div>
				Player 2:
				{#if room.other_player}
					<span class="font-bold">{room.other_player}</span>
				{:else}
					<span class="text-zinc-500 italic">Join to play!</span>
				{/if}
			</div>
		</div>
		<div class="mt-auto">
			<div class="text-xs text-zinc-400">Click to join →</div>
		</div>
	</div>
</a>
