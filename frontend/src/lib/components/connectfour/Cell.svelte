<script lang="ts">
	let {
		value = '',
		isWinning = false,
		isPreview = false
	} = $props<{
		value: 'R' | 'B' | '';
		isWinning?: boolean;
		isPreview?: boolean;
	}>();

	function getCellClass() {
		if (value === 'R') return 'bg-red-400';
		if (value === 'B') return 'bg-blue-400';
		if (isPreview) return 'bg-zinc-800';
		return 'bg-zinc-950';
	}
</script>

<div
	class="aspect-square w-full border border-zinc-700 transition-colors sm:border-2
	{getCellClass()}
	{isWinning
		? value == 'B'
			? 'ring-4 ring-blue-300 ring-inset'
			: 'ring-4 ring-red-300 ring-inset'
		: ''}"
></div>
