<script lang="ts">
	let {
		value = '',
		disabled,
		onclick,
		onmouseenter,
		onmouseleave,
		isPreview = false,
		canPlace = false
	} = $props<{
		value: 'H' | 'V' | '';
		disabled: boolean;
		onclick: () => void;
		onmouseenter?: () => void;
		onmouseleave?: () => void;
		isPreview?: boolean;
		canPlace?: boolean;
	}>();

	function getCellClass() {
		if (value === 'H') return 'bg-lime-400 border-stone-950';
		if (value === 'V') return 'bg-rose-400 border-stone-950';
		if (isPreview && canPlace) return 'bg-cyan-400/20 hover:bg-cyan-400/20';
		if (isPreview && !canPlace) return 'bg-red-500/30';
		return 'bg-stone-950 hover:bg-stone-800';
	}
</script>

<button
	{onclick}
	{onmouseenter}
	{onmouseleave}
	{disabled}
	class="aspect-square w-14 border-2 border-stone-700 text-5xl font-bold text-stone-950 transition-colors disabled:cursor-not-allowed md:w-20
	{getCellClass()}
	{disabled && !value ? 'opacity-50' : ''}"
>
	{value}
</button>
