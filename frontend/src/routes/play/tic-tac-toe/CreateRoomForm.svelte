<script lang="ts">
	import { Switch } from '@skeletonlabs/skeleton-svelte';

    let {
		onRoomCreate
	}: {
		onRoomCreate: (options: {
			playAs: string;
			private: boolean;
			roomName: string;
		}) => void;
	} = $props();

	const playAs = ['Random', 'X', 'O'];
	let as = $state(playAs[0]);
	let privateRoom = $state(false);
    let roomName = $state('');

	function setPlayAs(selected: string) {
		as = selected;
	}

    function handleCreateRoom() {
		onRoomCreate({ playAs: as, private: privateRoom, roomName: roomName });
	}
</script>

<div class="bg-surface-900 rounded-md p-4 m-2 min-w-120">
	<form class="flex flex-col gap-4" onsubmit={handleCreateRoom}>
        <label class="me-2" for="roomName">Room's name: </label>
        <input name="roomName" class="input" type="text" placeholder="Name of the room" />
		<div class="flex items-center justify-between">
            <div class="flex">
                <label class="me-4" for="private">Private room: </label>
                <Switch
                    name="private"
                    checked={privateRoom}
                    onCheckedChange={(e) => (privateRoom = e.checked)}
                />
            </div>
            <div class="flex">
                <label for="private">Play as: </label>
                <div class="flex gap-2 ms-4">
                    {#each playAs as p}
                        <button
                            type="button"
                            class={`chip capitalize ${as === p ? 'preset-filled' : 'preset-tonal'}`}
                            onclick={() => setPlayAs(p)}
                        >
                            <span>{p}</span>
                        </button>
                    {/each}
                </div>
            </div>
		</div>
		<button type="submit" class="btn btn-lg preset-filled-primary-500 mt-4">Create room</button>
	</form>
</div>
