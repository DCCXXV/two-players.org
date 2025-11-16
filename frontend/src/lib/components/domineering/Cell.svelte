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
		if (value === 'H') return 'bg-blue-400 border-zinc-950';
		if (value === 'V') return 'bg-red-400 border-zinc-950';
		if (isPreview && canPlace) return 'bg-green-400/10';
		if (isPreview && !canPlace) return 'bg-red-400/10';
		return 'bg-zinc-950 hover:bg-zinc-800';
	}
</script>

<button
	{onclick}
	{onmouseenter}
	{onmouseleave}
	{disabled}
	class="aspect-square w-14 border-2 border-zinc-700 text-5xl font-bold text-zinc-950 transition-colors disabled:cursor-not-allowed md:w-20
	{getCellClass()}
	{disabled && !value ? 'opacity-50' : ''}"
>
	{value}
</button>
