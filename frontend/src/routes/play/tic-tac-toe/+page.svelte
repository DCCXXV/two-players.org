<script lang="ts">
    import { onMount } from 'svelte';
    import CreateRoomForm from './CreateRoomForm.svelte';

    interface Room {
        id: string;
        Name: string;
        IsPrivate: boolean;
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
			console.log('Intentando fetch a /api/v1/rooms...');
			const response = await fetch('/api/v1/rooms');
			console.log('Respuesta recibida:', {
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
				IsPrivate: !!room.is_private
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
        console.log('Room creation requested with options:', options);
        isCreatingRoom = true;
        try {
            const response = await fetch('/api/v1/rooms', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
				body: JSON.stringify({
					name: options.Name,
					game_type: "Tic Tac Toe",
					is_private: options.IsPrivate
				})
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({ message: 'Failed to create room' }));
                throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
            }

            const newRoom = await response.json();
            console.log('Room created successfully:', newRoom);

            await loadRooms();
        } catch (error) {
            console.error('Error creating room:', error);
            alert(`Error creating room: ${error instanceof Error ? error.message : 'Unknown error'}`);
        } finally {
            isCreatingRoom = false;
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
        <a class="text-primary-400 hover:underline" href="/play/tic-tac-toe"
            >Tic-Tac-Toe</a
        >
    </li>
</ol>

<h3 class="h3 lora-700 text-surface-200 my-4">Create a room</h3>

<CreateRoomForm onRoomCreate={handleRoomCreation} />

<h3 class="h3 lora-700 text-surface-200 my-4">Available rooms</h3>

{#if isLoadingRooms}
    <p class="text-surface-400">Loading rooms...</p>
{:else if errorLoadingRooms}
    <p class="text-error-500">Error: {errorLoadingRooms}</p>
    <button type="button" class="btn preset-outline-primary" onclick={loadRooms}>Try again</button>
{:else if availableRooms.length > 0}
<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-7 gap-4">
    {#each availableRooms as room (room.id)}
        <div class="bg-surface-900 p-4 shadow flex flex-col gap-2 aspect-square">
            <h4 class="text-lg text-primary-400 lora-700 text-pretty">{room.Name}</h4>
            <span class="text-xs text-surface-400">
                {room.IsPrivate ? 'Private' : 'Public'}
            </span>
            <div class="text-surface-300 text-sm">
				<!--TO DO-->
                <div>Player 1: <span class="font-bold">{room.CreatedBy}</span></div>
                <div>
                    {#if room.OtherPlayer}
                        Player 2: <span class="font-bold">{room.OtherPlayer}</span>
                    {:else}
                        Player 2:
                    {/if}
                </div>
            </div>
        </div>
    {/each}
</div>

{:else}
    <p class="text-surface-400">No public rooms available. Create one!</p>
{/if}

{#if isCreatingRoom}
    <div class="fixed inset-0 bg-black/50 flex items-center justify-center opacity-80">
        <p class="text-white text-xl lora-700">Creating room...</p>
    </div>
{/if}
