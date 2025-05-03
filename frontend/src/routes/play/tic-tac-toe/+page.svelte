<script lang="ts">
    import CreateRoomForm from './CreateRoomForm.svelte';

    function handleRoomCreation(options: { playAs: string, private: boolean }) {
        console.log("Room creation requested with options:", options);
        // Add your logic here to actually create the room,
        // perhaps call an API or update global state.
        // Example: api.createRoom(options);
    }
</script>

<ol class="flex items-center gap-4">
    <li><a class="opacity-60 hover:underline" href="/play">Play</a></li>
    <li class="opacity-50" aria-hidden="true">&rsaquo;</li>
    <li><a class="text-primary-400 hover:underline" href="/play/tic-tac-toe">Tic-Tac-Toe</a></li>
</ol>

<h3 class="h3 lora-700 text-surface-200 my-4">Create a room</h3>

<CreateRoomForm onRoomCreate={handleRoomCreation} />

<h3 class="h3 lora-700 text-surface-200 my-4">Available rooms</h3>
Nothing here...
