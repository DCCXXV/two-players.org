<script lang="ts">
    import { displayName } from '$lib/socketStore';
    import { Switch } from '@skeletonlabs/skeleton-svelte';

    let {
        onRoomCreate
    }: {
        onRoomCreate: (options: {
            Name: string;
            GameType: string;
            IsPrivate: boolean;
        }) => void;
    } = $props();

    const playAsOptions = ['Random', 'X', 'O'];
    let selectedPlayAs = $state(playAsOptions[0]);
    let isPrivateRoom = $state(false);
    let roomNameValue = $state($displayName + '\'s room');

    function setPlayAs(selected: string) {
        selectedPlayAs = selected;
    }

    function handleCreateRoom(event: Event) {
        event.preventDefault();
        if (!roomNameValue.trim()) {
            alert('Room name cannot be empty!');
            return;
        }
        onRoomCreate({
            Name: roomNameValue,
            GameType: "Tic Tac Toe",
            IsPrivate: isPrivateRoom
        });
        roomNameValue = $displayName + '\'s room';
        selectedPlayAs = playAsOptions[0];
        isPrivateRoom = false;
    }
</script>

<div class="bg-surface-900 p-4 w-full md:max-w-120">
    <form class="flex flex-col gap-4" onsubmit={handleCreateRoom}>
        <label for="roomName">Room's name: </label>
        <input
            name="roomName"
            class="input text-2xl text-primary-400 lora-700 rounded-none border-primary-400 border-1"
            type="text"
            placeholder={"Name of the room"}
            bind:value={roomNameValue}
        />
        <div class="flex items-center justify-between">
            <div class="flex">
                <label class="me-4" for="private">Private room: </label>
                <Switch
                    name="private"
                    checked={isPrivateRoom}
                    onCheckedChange={(e) => (isPrivateRoom = e.checked)}
                />
            </div>
            <div class="flex">
                <label for="private">Play as: </label>
                <div class="flex gap-2 ms-4">
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
        <button type="submit" class="btn btn-lg preset-filled-primary-500 rounded-none">Create</button>
    </form>
</div>
