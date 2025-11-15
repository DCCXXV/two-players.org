<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import CreateRoomForm from './CreateRoomForm.svelte';
	import { displayName } from '$lib/socketStore';

	interface Room {
		id: string;
		Name: string;
		IsPrivate: boolean;
		CreatedBy?: string | null;
		OtherPlayer?: string | null;
	}

	let availableRooms = $state<Room[]>([]);
	let isLoadingRooms = $state(true);
	let errorLoadingRooms = $state<string | null>(null);
	let isCreatingRoom = $state(false);

	async function loadRooms() {
		console.group('loadRooms()');
		isLoadingRooms = true;
		errorLoadingRooms = null;
		console.log('Initial state isLoadingRooms=true, errorLoadingRooms=null');
		try {
			console.log('Trying to fetch /api/v1/rooms...');
			const response = await fetch(
				import.meta.env.VITE_SOCKET_URL + '/api/v1/rooms?game_type=domineering'
			);
			console.log('Response received:', {
				ok: response.ok,
				status: response.status,
				statusText: response.statusText
			});
			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Failed to fetch rooms' }));
				throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			console.log('JSON received from /api/v1/rooms:', data);

			availableRooms = data.map((room: any) => ({
				id: room.id,
				Name: room.name,
				IsPrivate: !!room.is_private,
				CreatedBy: room.created_by,
				OtherPlayer: room.other_player
			}));
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
		console.group('handleRoomCreation()');
		isCreatingRoom = true;
		console.log('isCreatingRoom set to true');
		try {
			console.log('Attempting to fetch ' + import.meta.env.VITE_SOCKET_URL + '/api/v1/rooms');
			const response = await fetch(import.meta.env.VITE_SOCKET_URL + '/api/v1/rooms', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'X-Display-Name': $displayName
				},
				body: JSON.stringify({
					name: options.Name,
					game_type: 'domineering',
					is_private: options.IsPrivate
				})
			});
			console.log('Fetch response received:', {
				ok: response.ok,
				status: response.status,
				statusText: response.statusText
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Failed to create room' }));
				throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
			}

			const newRoom = await response.json();
			console.log('Room created successfully:', newRoom);
			sessionStorage.setItem(`room_${newRoom.id}`, JSON.stringify(newRoom));
			console.log('Navigating to:', `/play/domineering/${newRoom.id}`);
			goto(`/play/domineering/${newRoom.id}`);
		} catch (error) {
			console.error('Error creating room:', error);
			alert(`Error creating room: ${error instanceof Error ? error.message : 'Unknown error'}`);
		} finally {
			isCreatingRoom = false;
			console.log('isCreatingRoom set to false');
			console.groupEnd();
		}
	}

	onMount(() => {
		loadRooms();
	});
</script>

<ol class="flex items-center gap-4">
	<li><a class="opacity-60 hover:underline" href="/play">Play</a></li>
	<li class="opacity-50" aria-hidden="true">&rsaquo;</li>
	<li>
		<a class="text-lime-400 hover:underline" href="/play/domineering">Domineering</a>
	</li>
</ol>

{#if $displayName}
	<CreateRoomForm onRoomCreate={handleRoomCreation} displayName={$displayName} />
{:else}
	<div class="flex w-full items-center justify-center border-2 border-stone-500 p-4 md:max-w-120">
		<p class="text-stone-400">Initializing...</p>
	</div>
{/if}
<h3 class="my-4 text-2xl text-stone-400">Available rooms</h3>

{#if isLoadingRooms}
	<p class="text-stone-400">Loading rooms...</p>
{:else if errorLoadingRooms}
	<p class="text-error-500">Error: {errorLoadingRooms}</p>
	<button type="button" class="btn preset-outline-primary" onclick={loadRooms}>Try again</button>
{:else if availableRooms.length > 0}
	<div class="grid grid-cols-1 gap-4 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-6">
		{#each availableRooms as room (room.id)}
			<a href={`/play/domineering/${room.id}`} class="group">
				<div
					class="flex flex-col gap-2 bg-stone-900 p-4 shadow transition-colors group-hover:bg-stone-800 sm:aspect-square"
				>
					<h4 class="text-primary-400 lora-700 text-lg text-pretty">{room.Name}</h4>
					<div class="space-y-1 text-sm text-stone-300">
						<div>Player 1: <span class="font-bold">{room.CreatedBy}</span></div>
						<div>
							Player 2:
							{#if room.OtherPlayer}
								<span class="font-bold">{room.OtherPlayer}</span>
							{:else}
								<span class="text-stone-500 italic">Join to play!</span>
							{/if}
						</div>
					</div>
					<div class="mt-auto">
						<div class="text-xs text-stone-400">Click to join →</div>
					</div>
				</div>
			</a>
		{/each}
	</div>
{:else}
	<p class="text-stone-400">No public rooms available. Create one!</p>
{/if}
{#if isCreatingRoom}
	<div class="fixed inset-0 z-10 flex items-center justify-center bg-black/50 opacity-80">
		<p class="lora-700 text-xl">Creating room...</p>
	</div>
{/if}
