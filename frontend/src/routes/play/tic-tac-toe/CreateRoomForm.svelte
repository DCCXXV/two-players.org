<script lang="ts">
	import { Switch } from '@skeletonlabs/skeleton-svelte';

	let {
		onRoomCreate,
		displayName
	}: {
		onRoomCreate: (options: { Name: string; GameType: string; IsPrivate: boolean }) => void;
		displayName: string;
	} = $props();

	const playAsOptions = ['Random', 'X', 'O'];
	let selectedPlayAs = $state(playAsOptions[0]);
	let isPrivateRoom = $state(false);
	let roomNameValue = $state('');
	let userHasEditedName = false;

	$effect(() => {
		if (displayName && !userHasEditedName) {
			roomNameValue = displayName + "'s room";
		}
	});

	function setPlayAs(selected: string) {
		selectedPlayAs = selected;
	}

	function handleCreateRoom(event: Event) {
		event.preventDefault();
		if (!roomNameValue.trim()) {
			alert('Room name cannot be empty!');
			return;
		}
		const roomDetails = {
			Name: roomNameValue,
			GameType: 'tic-tac-toe',
			IsPrivate: isPrivateRoom
		};
		onRoomCreate(roomDetails);
		userHasEditedName = false;
		roomNameValue = displayName ? displayName + "'s room" : '';
		selectedPlayAs = playAsOptions[0];
		isPrivateRoom = false;
	}
</script>

<div class="border-surface-500 w-full border-2 p-4 md:max-w-120">
	<form class="flex flex-col gap-4" onsubmit={handleCreateRoom}>
		<label for="roomName">Room's name: </label>
		<input
			name="roomName"
			class="input text-primary-400 lora-700 bg-surface-800 rounded-none text-2xl"
			type="text"
			placeholder={'Name of the room'}
			bind:value={roomNameValue}
			oninput={() => (userHasEditedName = true)}
		/>
		<div class="flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
			<div class="flex items-center">
				<label class="me-4" for="private">Private room: </label>
				<Switch
					name="private"
					checked={isPrivateRoom}
					onCheckedChange={(e) => (isPrivateRoom = e.checked)}
				/>
			</div>
			<div class="flex items-center">
				<label>Play as: </label>
				<div class="ms-4 flex flex-wrap gap-2">
					{#each playAsOptions as p}
						<button
							type="button"
							class={`chip capitalize ${selectedPlayAs === p ? 'preset-filled rounded-none' : 'preset-tonal rounded-none'}`}
							onclick={() => setPlayAs(p)}
						>
							<span>{p}</span>
						</button>
					{/each}
				</div>
			</div>
		</div>
		<button
			type="submit"
			class="btn btn-lg preset-filled-primary-500 rounded-none"
			disabled={!displayName || !roomNameValue.trim()}
		>
			Create
		</button>
	</form>
</div>
