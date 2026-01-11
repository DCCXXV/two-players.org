<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import type { GameConfig } from '$lib/config/games';
	import CreateRoomForm from '$lib/components/forms/CreateRoomForm.svelte';
	import RoomCard from '$lib/components/ui/RoomCard.svelte';
	import Collapsible from '$lib/components/ui/Collapsible.svelte';
	import GameInfo from '$lib/components/game/GameInfo.svelte';
	import { displayName, roomListUpdates } from '$lib/socketStore';

	interface Room {
		id: string;
		name: string;
		game_type: string;
		is_private: boolean;
		created_by?: string | null;
		other_player?: string | null;
	}

	interface Props {
		gameConfig: GameConfig;
	}

	let { gameConfig }: Props = $props();

	let availableRooms = $state<Room[]>([]);
	let isLoadingRooms = $state(true);
	let errorLoadingRooms = $state<string | null>(null);
	let isCreatingRoom = $state(false);

	async function loadRooms() {
		isLoadingRooms = true;
		errorLoadingRooms = null;

		try {
			const response = await fetch(
				import.meta.env.VITE_SOCKET_URL + `/api/v1/rooms?game_type=${gameConfig.id}`
			);

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Failed to fetch rooms' }));
				throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
			}

			const data = await response.json();
			availableRooms = data;
		} catch (error) {
			console.error('Error loading rooms:', error);
			errorLoadingRooms = error instanceof Error ? error.message : 'An unknown error occurred.';
			availableRooms = [];
		} finally {
			isLoadingRooms = false;
		}
	}

	async function handleRoomCreation(options: {
		Name: string;
		GameType: string;
		IsPrivate: boolean;
	}) {
		isCreatingRoom = true;

		try {
			const response = await fetch(import.meta.env.VITE_SOCKET_URL + '/api/v1/rooms', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'X-Display-Name': $displayName
				},
				body: JSON.stringify({
					name: options.Name,
					game_type: gameConfig.id,
					is_private: options.IsPrivate
				})
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Failed to create room' }));
				throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
			}

			const newRoom = await response.json();
			sessionStorage.setItem(`room_${newRoom.id}`, JSON.stringify(newRoom));
			goto(`/play/${gameConfig.path}/${newRoom.id}`);
		} catch (error) {
			console.error('Error creating room:', error);
			alert(`Error creating room: ${error instanceof Error ? error.message : 'Unknown error'}`);
		} finally {
			isCreatingRoom = false;
		}
	}

	// Listen for real-time room list updates
	$effect(() => {
		if ($roomListUpdates && $roomListUpdates.game_type === gameConfig.id) {
			availableRooms = $roomListUpdates.rooms;
			isLoadingRooms = false;
		}
	});

	onMount(() => {
		loadRooms();
	});
</script>

<ol class=" mb-6 flex items-center gap-4">
	<li><a class="opacity-60 hover:underline" href="/play">Play</a></li>
	<li class="opacity-50" aria-hidden="true">&rsaquo;</li>
	<li>
		<a class="text-blue-400 hover:underline" href={`/play/${gameConfig.path}`}
			>{gameConfig.displayName}</a
		>
	</li>
</ol>

<div class="flex flex-col gap-4 lg:flex-row">
	<div class="w-full lg:w-3/5">
		{#if $displayName}
			<div class="flex flex-col gap-3 lg:flex-row">
				<CreateRoomForm {gameConfig} displayName={$displayName} onRoomCreate={handleRoomCreation} />
				<div class="w-full bg-zinc-900 text-amber-300">
					<h3 class="font-bold">NOTE</h3>
					<p>
						Currently if you want to play vs yourself you will have to open two-players.org in new
						tab and look for your room or copy and paste the room link when you create it to get in
						directly.
					</p>
				</div>
			</div>
		{:else}
			<div
				class="flex w-full items-center justify-center border-2 border-zinc-500 p-4 md:max-w-120"
			>
				<p class="text-zinc-400">Initializing...</p>
			</div>
		{/if}

		<section class="my-8">
			<Collapsible title="Available rooms">
				{#if isLoadingRooms}
					<div class="w-full p-8 text-center">
						<p class="text-zinc-400">Loading rooms...</p>
					</div>
				{:else if errorLoadingRooms}
					<div class="w-full p-8 text-center">
						<p class="text-error-500">Error: {errorLoadingRooms}</p>
						<button type="button" class="btn preset-outline-blue mt-4" onclick={loadRooms}
							>Try again</button
						>
					</div>
				{:else if availableRooms.length > 0}
					{#each availableRooms as room (room.id)}
						<RoomCard {room} showGameType={false} />
					{/each}
				{:else}
					<div class="w-full p-8 text-center">
						<p class="text-zinc-400">No public rooms available. Create one!</p>
					</div>
				{/if}
			</Collapsible>
		</section>
	</div>

	<div class="w-full lg:w-2/5">
		<GameInfo {gameConfig} />
	</div>
</div>

{#if isCreatingRoom}
	<div class="fixed inset-0 z-10 flex items-center justify-center bg-black/50 opacity-80">
		<p class="lora-700 text-xl">Creating room...</p>
	</div>
{/if}
