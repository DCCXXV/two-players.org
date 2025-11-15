<script lang="ts">
	import { onMount } from 'svelte';
	import GameCard from '$lib/components/ui/GameCard.svelte';
	import Collapsible from '$lib/components/ui/Collapsible.svelte';
	import RoomCard from '$lib/components/ui/RoomCard.svelte';
	import { roomListUpdates } from '$lib/socketStore';

	interface Room {
		id: string;
		name: string;
		game_type: string;
		is_private: boolean;
		created_by?: string | null;
		other_player?: string | null;
		created_at?: string;
	}

	let allRooms = $state<Room[]>([]);
	let isLoadingRooms = $state(true);
	let errorLoadingRooms = $state<string | null>(null);

	const gameTypes = ['tic-tac-toe', 'domineering', 'nim', 'dots-and-boxes'];

	async function loadAllRooms() {
		isLoadingRooms = true;
		errorLoadingRooms = null;

		try {
			// Fetch rooms from all game types in parallel
			const responses = await Promise.all(
				gameTypes.map((gameType) =>
					fetch(import.meta.env.VITE_SOCKET_URL + `/api/v1/rooms?game_type=${gameType}`)
				)
			);

			// Check if all responses are ok
			const failedResponse = responses.find((r) => !r.ok);
			if (failedResponse) {
				throw new Error(`HTTP error! status: ${failedResponse.status}`);
			}

			// Parse all JSON responses
			const roomsData = await Promise.all(responses.map((r) => r.json()));

			// Flatten and combine all rooms
			const combinedRooms: Room[] = roomsData.flat().map((room: any) => ({
				id: room.id,
				name: room.name,
				game_type: room.game_type,
				is_private: !!room.is_private,
				created_by: room.created_by,
				other_player: room.other_player,
				created_at: room.created_at
			}));

			// Sort by created_at (oldest first) so people waiting get priority
			allRooms = combinedRooms.sort((a, b) => {
				if (!a.created_at || !b.created_at) return 0;
				return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
			});
		} catch (error) {
			console.error('Error loading rooms:', error);
			errorLoadingRooms = error instanceof Error ? error.message : 'Failed to load rooms';
			allRooms = [];
		} finally {
			isLoadingRooms = false;
		}
	}

	// Listen for real-time room list updates
	$effect(() => {
		if ($roomListUpdates) {
			// Update the specific game type's rooms in our combined list
			const gameType = $roomListUpdates.game_type;
			const updatedRooms = $roomListUpdates.rooms;

			// Remove old rooms of this game type and add new ones
			allRooms = [...allRooms.filter((r) => r.game_type !== gameType), ...updatedRooms].sort(
				(a, b) => {
					if (!a.created_at || !b.created_at) return 0;
					return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
				}
			);

			isLoadingRooms = false;
		}
	});

	onMount(() => {
		loadAllRooms();
	});
</script>

<section class="my-8">
	<Collapsible title="Available games">
		<GameCard title="Tic Tac Toe" path="tic-tac-toe" />
		<GameCard title="Domineering" path="domineering" />
		<GameCard title="Nim" path="nim" />
		<GameCard title="Dots and Boxes" path="dots-and-boxes" />
	</Collapsible>
</section>

<section class="my-8">
	<Collapsible title="Current rooms">
		{#if isLoadingRooms}
			<div class="w-full p-8 text-center">
				<p class="text-stone-400"></p>
			</div>
		{:else if errorLoadingRooms}
			<div class="w-full p-8 text-center">
				<p class="text-error-500">Error: {errorLoadingRooms}</p>
				<button type="button" class="btn preset-outline mt-4" onclick={loadAllRooms}
					>Try again</button
				>
			</div>
		{:else if allRooms.length > 0}
			{#each allRooms as room (room.id)}
				<RoomCard {room} />
			{/each}
		{:else}
			<div class="w-full p-8 text-center">
				<p class="text-stone-400">No public rooms available. Create one from a game page!</p>
			</div>
		{/if}
	</Collapsible>
</section>
