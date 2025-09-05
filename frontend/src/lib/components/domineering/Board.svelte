<script lang="ts">
	import Cell from './Cell.svelte';

	let { board, disabled, onMove } = $props<{
		board: ('H' | 'V' | '')[];
		disabled: boolean;
		onMove: (cellIndex: number) => void;
	}>();
</script>
