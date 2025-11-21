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
		value: string;
		disabled: boolean;
		onclick: () => void;
		onmouseenter?: () => void;
		onmouseleave?: () => void;
		isPreview?: boolean;
		canPlace?: boolean;
	}>();

	function getLineClass() {
		if (value === 'P1') return 'bg-zinc-950 border-x-5 border-blue-400';
		if (value === 'P2') return 'bg-zinc-950 border-x-5 border-dotted border-red-400';
		if (isPreview && canPlace) return 'bg-green-400/30';
		if (isPreview && !canPlace) return 'bg-red-400/30';
		if (disabled && !value) return 'bg-zinc-800';
		return 'bg-zinc-700';
	}
</script>

<button
	{onclick}
	{onmouseenter}
	{onmouseleave}
	{disabled}
	class="h-full w-full transition-colors disabled:cursor-not-allowed
	{getLineClass()}"
>
</button>
