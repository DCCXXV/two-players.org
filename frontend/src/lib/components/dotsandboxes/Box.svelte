<script lang="ts">
	let { value = '' } = $props<{
		value: string;
	}>();

	function getTextClass() {
		if (value === 'P1') return 'text-blue-400';
		if (value === 'P2') return 'text-red-400';
		return '';
	}
</script>

<div
	class="flex h-full w-full items-center justify-center bg-zinc-950 text-2xl font-bold transition-colors sm:text-3xl md:text-4xl
	{getTextClass()}"
>
	{value}
</div>
