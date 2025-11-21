<script lang="ts">
	import { onMount } from 'svelte';
	import GameCard from '$lib/components/ui/GameCard.svelte';
	import Collapsible from '$lib/components/ui/Collapsible.svelte';
	import RoomCard from '$lib/components/ui/RoomCard.svelte';
	import { roomListUpdates } from '$lib/socketStore';
	import { getAllGameConfigs } from '$lib/config/games';
	import { page } from '$app/stores';

	const games = getAllGameConfigs();

	const title = 'Play Games | Two Players';
	const description =
		'Browse and join available game rooms or create your own. Play Tic-Tac-Toe, Domineering, Nim, Dots & Boxes and more!';
	const imageUrl = `${$page.url.origin}/img/tic-tac-toe-preview.png`;
	const url = $page.url.href;

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

<svelte:head>
	<title>{title}</title>
	<meta name="description" content={description} />
	<meta property="og:type" content="website" />
	<meta property="og:url" content={url} />
	<meta property="og:title" content={title} />
	<meta property="og:description" content={description} />
	<meta property="og:image" content={imageUrl} />
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:url" content={url} />
	<meta name="twitter:title" content={title} />
	<meta name="twitter:description" content={description} />
	<meta name="twitter:image" content={imageUrl} />
</svelte:head>

<section class="my-8">
	<Collapsible title="Available games">
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
			{#each games as game}
				<a
					href="/play/{game.path}"
					class="rounded-0 group block border-b-1 border-zinc-700 bg-transparent p-4 transition-all hover:border-blue-400"
				>
					<h3 class="mb-1 text-lg font-bold text-zinc-300 group-hover:text-blue-400">
						{game.displayName}
					</h3>
				</a>
			{/each}
		</div>
	</Collapsible>
</section>

<section class="my-8">
	<Collapsible title="Current rooms">
		{#if isLoadingRooms}
			<div class="w-full p-8 text-center">
				<p class="text-zinc-400"></p>
			</div>
		{:else if errorLoadingRooms}
			<div class="w-full p-8 text-center">
				<p class="text-error-500">Error: {errorLoadingRooms}</p>
				<button type="button" class="btn preset-outline mt-4" onclick={loadAllRooms}
					>Try again</button
				>
			</div>
		{:else if allRooms.length > 0}
			<div class="flex gap-4">
				{#each allRooms as room (room.id)}
					<RoomCard {room} />
				{/each}
			</div>
		{:else}
			<div class="w-full p-8 text-center">
				<p class="text-zinc-400">No public rooms available. Create one from a game page!</p>
			</div>
		{/if}
	</Collapsible>
</section>
