<script lang="ts">
	interface Props {
		sticks: number;
		disabled: boolean;
		onMove: (sticksToTake: number) => void;
	}

	let { sticks, disabled, onMove }: Props = $props();

	let selectedCount = $state(0);

	function selectCount(count: number) {
		if (disabled || count > sticks) return;
		selectedCount = count;
	}

	function handleTakeSticks() {
		if (!disabled && selectedCount > 0 && selectedCount <= 3) {
			onMove(selectedCount);
			selectedCount = 0;
		}
	}

	$effect(() => {
		sticks;
		selectedCount = 0;
	});
</script>

<div class="flex flex-col items-center gap-8 p-4 md:p-8">
	<div
		class="relative flex items-end justify-center"
		style="height: 200px; width: 100%; max-width: min(90vw, 600px);"
	>
		{#each Array(sticks) as _, i (i)}
			{@const maxWidth =
				typeof window !== 'undefined' ? Math.min(window.innerWidth * 0.9, 600) : 600}
			{@const spacing = Math.min(32, (maxWidth - 40) / sticks)}
			<div
				class="absolute bottom-0 h-32 w-2 bg-zinc-700"
				style="transform: translateX({(i - sticks / 2) * spacing}px);"
			></div>
		{/each}
	</div>

	<div class="text-center text-sm text-zinc-400">
		<p class="text-lg font-bold text-zinc-400">{sticks} sticks remaining</p>
	</div>

	<div class="flex flex-col items-center gap-4">
		<p class="text-lg text-zinc-300">How many sticks will you take?</p>

		<div class="flex gap-3">
			{#each [1, 2, 3] as count}
				<button
					type="button"
					class="rounded-0 h-20 w-20 border-b-2 text-2xl font-bold transition-all {selectedCount ===
					count
						? 'border-blue-400 bg-blue-400 text-zinc-950'
						: 'border-zinc-600 bg-zinc-900 text-zinc-600 hover:border-zinc-500 hover:bg-zinc-800'} {disabled ||
					count > sticks
						? 'cursor-not-allowed opacity-30'
						: 'cursor-pointer'}"
					onclick={() => selectCount(count)}
					disabled={disabled || count > sticks}
				>
					{count}
				</button>
			{/each}
		</div>

		<button
			class="cursor-pointer rounded-none border-b-2 border-blue-900 bg-blue-400 px-12 py-4 text-xl font-bold text-zinc-950 transition-all hover:bg-blue-500 disabled:border-zinc-700 disabled:cursor-not-allowed disabled:bg-zinc-900 disabled:text-zinc-400"
			disabled={disabled || selectedCount === 0}
			onclick={handleTakeSticks}
		>
			Take {selectedCount || ''}
			{selectedCount === 1 ? 'Stick' : 'Sticks'}
		</button>
	</div>
</div>
