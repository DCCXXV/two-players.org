<script lang="ts">
	import Cell from './Cell.svelte';

	let { board, winningCells, disabled, onMove } = $props<{
		board: ('R' | 'B' | '')[][];
		winningCells: [number, number][] | null;
		disabled: boolean;
		onMove: (column: number) => void;
	}>();

	let hoveredCol = $state<number | null>(null);

	function isWinningCell(row: number, col: number): boolean {
		if (!winningCells) return false;
		return winningCells.some(([r, c]: [number, number]) => r === row && c === col);
	}

	function isColumnFull(col: number): boolean {
		return board[0][col] !== '';
	}

	function getLowestEmptyRow(col: number): number {
		for (let row = board.length - 1; row >= 0; row--) {
			if (board[row][col] === '') return row;
		}
		return -1;
	}

	function isPreviewCell(row: number, col: number): boolean {
		if (hoveredCol === null || disabled) return false;
		if (col !== hoveredCol) return false;
		if (isColumnFull(col)) return false;
		return row === getLowestEmptyRow(col);
	}

	function handleColumnHover(col: number) {
		if (disabled || isColumnFull(col)) {
			hoveredCol = null;
			return;
		}
		hoveredCol = col;
	}

	function handleColumnLeave() {
		hoveredCol = null;
	}

	function handleClick(col: number) {
		if (disabled || isColumnFull(col)) return;
		onMove(col);
	}
</script>

<div class="mx-auto w-full max-w-[min(100vw-2rem,40rem)]">
	<div class="grid grid-cols-7 gap-0">
		{#each board as row, rowIndex (rowIndex)}
			{#each row as cellValue, colIndex (colIndex)}
				<button
					type="button"
					onclick={() => handleClick(colIndex)}
					onmouseenter={() => handleColumnHover(colIndex)}
					onmouseleave={handleColumnLeave}
					disabled={disabled || isColumnFull(colIndex)}
					class="cursor-pointer disabled:cursor-not-allowed"
				>
					<Cell
						value={cellValue}
						isWinning={isWinningCell(rowIndex, colIndex)}
						isPreview={isPreviewCell(rowIndex, colIndex)}
					/>
				</button>
			{/each}
		{/each}
	</div>
</div>
