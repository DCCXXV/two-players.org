<script lang="ts">
	import Cell from './Cell.svelte';

	let {
		board,
		disabled = false,
		onMove,
		mySymbol
	}: {
		board: string[][];
		disabled?: boolean;
		onMove: (row: number, col: number) => void;
		mySymbol?: string;
	} = $props();

	let hoveredCell = $state<{ row: number; col: number } | null>(null);

	function handleCellClick(row: number, col: number) {
		if (disabled) return;
		onMove(row, col);
	}

	function handleCellHover(row: number, col: number) {
		if (disabled || !mySymbol) {
			hoveredCell = null;
			return;
		}
		hoveredCell = { row, col };
	}

	function handleCellLeave() {
		hoveredCell = null;
	}

	function isPreviewCell(row: number, col: number): boolean {
		if (!hoveredCell || !mySymbol) return false;

		const { row: hRow, col: hCol } = hoveredCell;

		if (mySymbol === 'H') {
			if (hCol === 7) {
				return row === hRow && (col === hCol || col === hCol - 1);
			} else {
				return row === hRow && (col === hCol || col === hCol + 1);
			}
		} else if (mySymbol === 'V') {
			if (hRow === 7) {
				return col === hCol && (row === hRow || row === hRow - 1);
			} else {
				return col === hCol && (row === hRow || row === hRow + 1);
			}
		}

		return false;
	}

	function canPlaceAt(row: number, col: number): boolean {
		if (!mySymbol) return false;

		if (mySymbol === 'H') {
			if (col === 7) {
				return board[row][col] === '' && board[row][col - 1] === '';
			} else {
				return board[row][col] === '' && board[row][col + 1] === '';
			}
		} else if (mySymbol === 'V') {
			if (row === 7) {
				return board[row][col] === '' && board[row - 1][col] === '';
			} else {
				return board[row][col] === '' && board[row + 1][col] === '';
			}
		}

		return false;
	}
</script>

<div class="mx-auto grid w-full max-w-[min(100vw-2rem,40rem)] grid-cols-8 gap-0">
	{#each board as row, rowIndex (rowIndex)}
		{#each row as cell, colIndex (colIndex)}
			<Cell
				value={cell}
				onclick={() => handleCellClick(rowIndex, colIndex)}
				onmouseenter={() => handleCellHover(rowIndex, colIndex)}
				onmouseleave={handleCellLeave}
				{disabled}
				isPreview={isPreviewCell(rowIndex, colIndex)}
				canPlace={canPlaceAt(rowIndex, colIndex)}
			/>
		{/each}
	{/each}
</div>
