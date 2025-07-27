<script lang="ts">
	import { onMount } from 'svelte';
	import {
		connectWebSocket,
		displayName,
		isConnected,
		socket as socketStore
	} from '$lib/socketStore';
	import ActiveConnectionsSidebar from '$lib/components/ui/sideelements/ActiveConnectionsSidebar.svelte';
	import GameCard from '$lib/components/ui/cards/GameCard.svelte';
	export let data;

	let connections = data.connections || [];

	onMount(() => {
		if (!$socketStore) {
			connectWebSocket();
		}

		const unsubscribe = socketStore.subscribe((socket) => {
			if (socket) {
				socket.addEventListener('message', (event) => {
					const message = JSON.parse(event.data);
					if (message.type === 'connections_update') {
						connections = message.payload;
					}
				});
			}
		});

		return () => {
			unsubscribe();
		};
	});
</script>

<section class="grid w-full grid-cols-1 gap-6 lg:grid-cols-2">
	<GameCard
		title="Tic Tac Toe"
		path="tic-tac-toe"
		src="/images/games/tic-tac-toe.png"
		desc="On an initially empty 3×3 grid, players alternate placing their mark, X or O, in any vacant square. Turns alternate until a mark completes a row, column or diagonal of three, or the grid is full; the player who first completes such a line wins, and if the grid fills with no line the game is a draw."
		tags={['3×3 grid', 'Alignment/Impartial', 'Yes', 'Draw (perfect play)']}
	/>
	<GameCard
		title="Domineering"
		path="domineering"
		src="https://dummyimage.com/480x480/1F1919/5B8139&text=Coming%20Soon!"
		desc="On an initially empty rectangular grid (typically 8x8), players alternate laying a 1×2 domino—one vertically, the other horizontally—on any pair of adjacent empty squares. Turns alternate until a player has no legal placement. The last player to move wins, and draws are impossible."
		tags={['Rectangular grid', 'Partisan', 'Yes', 'Last player wins (perfect play)']}
	/>
	<GameCard
		title="Dots and Boxes"
		path="dots-and-boxes"
		src="https://dummyimage.com/480x480/1F1919/5B8139&text=Coming%20Soon!"
		desc="On a grid of dots, players take turns drawing a single horizontal or vertical line between two adjacent dots. Turns alternate until all possible lines have been drawn. When a player draws the fourth side of a 1x1 box, they claim that box and immediately take an extra turn. The player who has claimed more boxes wins, and if they claim an equal number, the game is a draw."
		tags={['Dot grid', 'Impartial', 'No', 'Highest score (draws possible)']}
	/>
</section>