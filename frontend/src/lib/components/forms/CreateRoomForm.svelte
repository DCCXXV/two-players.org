<script lang="ts">
	import type { GameConfig } from '$lib/config/games';

	interface Props {
		gameConfig: GameConfig;
		displayName: string;
		onRoomCreate: (options: { Name: string; GameType: string; IsPrivate: boolean }) => void;
	}

	let { gameConfig, displayName, onRoomCreate }: Props = $props();

	let isPrivateRoom = $state(false);
	let roomNameValue = $state('');
	let userHasEditedName = false;

	$effect(() => {
		if (displayName && !userHasEditedName) {
			roomNameValue = displayName + "'s room";
		}
	});

	function handleCreateRoom(event: Event) {
		event.preventDefault();
		if (!roomNameValue.trim()) {
			alert('Room name cannot be empty!');
			return;
		}
		const roomDetails = {
			Name: roomNameValue,
			GameType: gameConfig.id,
			IsPrivate: isPrivateRoom
		};
		onRoomCreate(roomDetails);
		userHasEditedName = false;
		roomNameValue = displayName ? displayName + "'s room" : '';
		isPrivateRoom = false;
	}
</script>

<div class="w-full bg-zinc-800 border-b-1 border-zinc-700 p-4 md:max-w-120">
	<form class="flex flex-col gap-4" onsubmit={handleCreateRoom}>
		<input
			name="roomName"
			class="input lora-700 rounded-none border-0 border-b-1 border-zinc-500 bg-zinc-800 text-2xl text-blue-400 hover:border-blue-400 focus:ring-0"
			type="text"
			placeholder="Name of the room"
			bind:value={roomNameValue}
			oninput={() => (userHasEditedName = true)}
		/>
		<div class="flex flex-col items-start justify-between gap-4 md:flex-row md:items-center"></div>
		<button
			type="submit"
			class="ms-auto w-50 cursor-pointer rounded-none bg-blue-400 border-b-2 border-r-2 border-blue-900 p-2 font-bold text-zinc-950"
			disabled={!displayName || !roomNameValue.trim()}
		>
			Create room
		</button>
	</form>
</div>
