<script lang="ts">
	import Dot from './Dot.svelte';
	import HLine from './HLine.svelte';
	import VLine from './VLine.svelte';
	import Box from './Box.svelte';

	let {
		hLines,
		vLines,
		boxes,
		disabled = false,
		onMove,
		mySymbol
	}: {
		hLines: string[][];
		vLines: string[][];
		boxes: string[][];
		disabled?: boolean;
		onMove: (type: string, row: number, col: number) => void;
		mySymbol?: string;
	} = $props();

	let hoveredLine = $state<{ type: string; row: number; col: number } | null>(null);

	function handleLineClick(type: string, row: number, col: number) {
		if (disabled) return;
		onMove(type, row, col);
	}

	function handleLineHover(type: string, row: number, col: number) {
		if (disabled || !mySymbol) {
			hoveredLine = null;
			return;
		}
		hoveredLine = { type, row, col };
	}

	function handleLineLeave() {
		hoveredLine = null;
	}

	function isPreviewLine(type: string, row: number, col: number): boolean {
		if (!hoveredLine || !mySymbol) return false;
		return hoveredLine.type === type && hoveredLine.row === row && hoveredLine.col === col;
	}

	function canPlaceHLine(row: number, col: number): boolean {
		return hLines[row][col] === '';
	}

	function canPlaceVLine(row: number, col: number): boolean {
		return vLines[row][col] === '';
	}

	// 9 by 9 grid with dots, hlines, vlines and boxes
	// Pos(gridRow, gridCol):
	//  (even, even) = Dot
	//  (even, odd) = HLine
	//  (odd, even) = VLine
	//  (odd, odd) = Box
</script>

<div
	class="mx-auto grid aspect-square w-full max-w-[min(100vw-2rem,30rem)] gap-0"
	style="grid-template-columns: 1fr 6fr 1fr 6fr 1fr 6fr 1fr 6fr 1fr; grid-template-rows: 1fr 6fr 1fr 6fr 1fr 6fr 1fr 6fr 1fr;"
>
	{#each Array(9) as _, gridRow (gridRow)}
		{#each Array(9) as _, gridCol (gridCol)}
			{@const isRowEven = gridRow % 2 === 0}
			{@const isColEven = gridCol % 2 === 0}

			{#if isRowEven && isColEven}
				<Dot />
			{:else if isRowEven && !isColEven}
				{@const row = gridRow / 2}
				{@const col = (gridCol - 1) / 2}
				<HLine
					value={hLines[row][col]}
					{disabled}
					onclick={() => handleLineClick('h', row, col)}
					onmouseenter={() => handleLineHover('h', row, col)}
					onmouseleave={handleLineLeave}
					isPreview={isPreviewLine('h', row, col)}
					canPlace={canPlaceHLine(row, col)}
				/>
			{:else if !isRowEven && isColEven}
				{@const row = (gridRow - 1) / 2}
				{@const col = gridCol / 2}
				<VLine
					value={vLines[row][col]}
					{disabled}
					onclick={() => handleLineClick('v', row, col)}
					onmouseenter={() => handleLineHover('v', row, col)}
					onmouseleave={handleLineLeave}
					isPreview={isPreviewLine('v', row, col)}
					canPlace={canPlaceVLine(row, col)}
				/>
			{:else}
				{@const row = (gridRow - 1) / 2}
				{@const col = (gridCol - 1) / 2}
				<Box value={boxes[row][col]} />
			{/if}
		{/each}
	{/each}
</div>
