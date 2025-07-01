<script lang="ts">
	import {
		isConnected,
		displayName,
		sendWebSocketMessage,
		players,
		gameState,
		roomClosedMessage
	} from '$lib/socketStore';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	export let data;

	let hasJoined = false;

	$: if ($roomClosedMessage) {
		alert($roomClosedMessage); // Show an alert when the room is closed
		$roomClosedMessage = null; // Reset the message
		goto('/play/tic-tac-toe'); // Redirect to the lobby
	}

	$: if ($isConnected && $displayName && !hasJoined && data.room) {
		hasJoined = true;
		sendWebSocketMessage({
			type: 'join_room',
			payload: {
				roomId: data.room.id
			}
		});
	}
</script>

{#if data.error}
	<div class="text-center">
		<h1 class="h3 lora-700 text-error-400 mb-8">Room Not Found</h1>
		<!-- <p class="text-surface-400 mb-4">{data.error}</p> -->
		<a href="/play/tic-tac-toe" class="btn bg-primary-400">Go to Lobby</a>
	</div>
{:else if $gameState}
	<h1 class="h3 lora-700 text-primary-400">{data.room.name}</h1>
	<p class="text-surface-500 mb-4">
		ID: {data.room.id} · Status: {$gameState.canStart
			? 'Ready to start!'
			: 'Waiting for players...'}
		({$gameState.playerCount}/{$gameState.maxPlayers} players)
	</p>

	<div>
		{#if $gameState?.players?.length > 0}
			<div class="border-surface-400 border-2">
				<div class="text-primary-400 bg-surface-800 b border-surface-400 border-b-2 p-2 font-bold">
					P1: <span class="text-surface-200">{$gameState.players[0] || 'Waiting...'}</span>
				</div>
				<div class="text-secondary-400 bg-surface-800 p-2 font-bold">
					P2: <span class="text-surface-200">{$gameState.players[1] || 'Waiting...'}</span>
				</div>
			</div>
			<!-- Espectadores si existen -->
			{#if $gameState.spectatorCount > 0}
				<details class="bg-surface-800 border-surface-400 mt-4 border-2 p-2">
					<summary class="h6 lora-700 text-surface-200 cursor-pointer">
						Spectators ({$gameState.spectatorCount})
					</summary>

					<ul class="bg-surface-950 text-surface-200 mt-2 p-1 text-sm">
						{#each $gameState.spectators as spectator}
							<li class="py-1">👁️ {spectator}</li>
						{/each}
					</ul>
				</details>
			{/if}

			<!-- Debug info temporal -->
			<details class="mt-4">
				<summary class="text-surface-400 cursor-pointer text-sm">Debug Info</summary>
				<pre class="bg-surface-900 mt-2 overflow-auto rounded p-2 text-xs">
			    {JSON.stringify($gameState, null, 2)}
            </pre>
			</details>
		{:else}
			<p class="text-surface-400">No players yet...</p>
		{/if}
	</div>
{:else}
	<div class="text-center">
		<h1 class="h3 lora-700 text-warning-400">Loading Room...</h1>
		<p class="text-surface-400 mb-4">Connecting to the room, please wait.</p>
	</div>
{/if}
