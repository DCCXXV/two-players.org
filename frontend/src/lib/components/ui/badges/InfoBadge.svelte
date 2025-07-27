<script lang="ts">
	let { label, value, tooltip }: { label: string; value: string; tooltip?: string } = $props();
</script>

<div class="border-surface-500 flex w-fit border-2">
	<div class="bg-surface-950 border-surface-500 text-surface-200 border-r-2 p-1 px-2 font-semibold">
		{label}
	</div>
	<div class="bg-surface-800 p-1 px-2">
		{value}
	</div>
</div>
