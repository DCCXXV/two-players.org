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
				class="absolute bottom-0 h-32 w-2 bg-stone-700"
				style="transform: translateX({(i - sticks / 2) * spacing}px);"
			></div>
		{/each}
	</div>

	<div class="flex flex-col items-center gap-4">
		<p class="text-lg text-stone-300">How many sticks will you take?</p>

		<div class="flex gap-3">
			{#each [1, 2, 3] as count}
				<button
					type="button"
					class="rounded-0 h-20 w-20 border-b-2 text-2xl font-bold transition-all {selectedCount ===
					count
						? 'border-lime-400 bg-lime-400 text-stone-950'
						: 'border-stone-600 bg-stone-900 text-stone-600 hover:border-stone-500 hover:bg-stone-800'} {disabled ||
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
			class="cursor-pointer rounded-none border-b-2 border-lime-900 bg-lime-400 px-12 py-4 text-xl font-bold text-stone-950 transition-all hover:bg-lime-500 disabled:border-stone-800 disabled:bg-stone-950 disabled:text-stone-400"
			disabled={disabled || selectedCount === 0}
			onclick={handleTakeSticks}
		>
			Take {selectedCount || ''}
			{selectedCount === 1 ? 'Stick' : 'Sticks'}
		</button>
	</div>

	<div class="text-center text-sm text-stone-400">
		<p class="text-lg font-bold text-amber-400">{sticks} sticks remaining</p>
	</div>
</div>
