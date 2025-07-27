<script lang="ts">
	import Cell from './Cell.svelte';

	let { board, disabled, onMove } = $props<{
		board: ('X' | 'O' | '')[];
		disabled: boolean;
		onMove: (cellIndex: number) => void;
	}>();
</script>

<div class="border-surface-900 blackboard mx-auto grid max-w-lg grid-cols-3 gap-2 border-16 p-2">
	{#each board as cellValue, i}
		<Cell value={cellValue} disabled={disabled || cellValue != ''} onClick={() => onMove(i)} />
	{/each}
</div>
