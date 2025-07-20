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
	import RematchButton from '$lib/components/ui/buttons/RematchButton.svelte';
	import Board from '$lib/components/tictactoe/Board.svelte';
	import GameStatus from '$lib/components/tictactoe/GameStatus.svelte';

	import type { PageData } from './$types';

	export let data: PageData;

	let hasJoined = false;

	let moveSoundX: HTMLAudioElement;
	let moveSoundO: HTMLAudioElement;
	let gameOverSound: HTMLAudioElement;

	let prevTurn: number | undefined = undefined;
	let gameHasEnded = false;

	onMount(() => {
		if (moveSoundX) {
			moveSoundX.volume = 0.5;
		}
		if (moveSoundO) {
			moveSoundO.volume = 0.5;
		}
		if (gameOverSound) {
			gameOverSound.volume = 0.4;
		}
	});

	$: if ($gameState && prevTurn === undefined) {
		prevTurn = $gameState.game.currentTurn;
	}

	$: if ($gameState && prevTurn !== undefined && $gameState.game.currentTurn !== prevTurn) {
		if (prevTurn === 0 && moveSoundX) {
			moveSoundX.currentTime = 0;
			moveSoundX.play().catch((e) => console.error('Error playing sound X', e));
		} else if (prevTurn === 1 && moveSoundO) {
			moveSoundO.currentTime = 0;
			moveSoundO.play().catch((e) => console.error('Error playing sound O', e));
		}
		prevTurn = $gameState.game.currentTurn;
	}

	$: if ($gameState && $gameState.game.winner && !gameHasEnded) {
		if (gameOverSound) {
			gameOverSound.currentTime = 0;
			gameOverSound.play().catch((e) => console.error('Error playing game over sound', e));
			gameHasEnded = true;
		}
	}

	$: if ($gameState && $gameState.game.winner === '' && gameHasEnded) {
		gameHasEnded = false;
	}

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

	$: myTurn = $gameState?.players[$gameState?.game.currentTurn] === $displayName;

	function onMove(cellIndex: number) {
		sendWebSocketMessage({
			type: 'make_move',
			payload: {
				cellIndex: cellIndex
			}
		});
	}

	function handleRematch() {
		sendWebSocketMessage({ type: 'rematch_request' });
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
	<div class="flex flex-col justify-between gap-4 lg:flex-row">
		<div class="w-full lg:w-1/5">
			<div>
				{#if $gameState?.players?.length > 0}
					<div class="border-surface-400 mb-4 border-2">
						<div class="flex">
							<div
								class="text-primary-400 bg-surface-900 border-surface-400 w-10 border-e-2 border-b-2 p-2 text-center font-bold"
							>
								<span>X</span>
							</div>
							<div
								class="text-primary-400 bg-surface-800 border-surface-400 w-full border-b-2 p-2 font-bold"
							>
								<span class="text-surface-200">{$gameState.players[0] || 'Waiting...'}</span>
							</div>
						</div>
						<div class="flex">
							<div
								class="text-secondary-400 bg-surface-900 border-surface-400 w-10 border-e-2 p-2 text-center font-bold"
							>
								<span>O</span>
							</div>
							<div class="text-primary-400 bg-surface-800 border-surface-400 w-full p-2 font-bold">
								<span class="text-surface-200">{$gameState.players[1] || 'Waiting...'}</span>
							</div>
						</div>
					</div>
					{#if $gameState.spectatorCount > 0}
						<details class="bg-surface-800 border-surface-400 mb-4 border-2 p-2">
							<summary class="text-surface-200 cursor-pointer font-bold">
								Spectators ({$gameState.spectatorCount})
							</summary>

							<ul class=" text-surface-200 p-1 text-sm">
								{#each $gameState.spectators as spectator}
									<li class="py-1">👁️ {spectator}</li>
								{/each}
							</ul>
						</details>
					{/if}
					{#if $gameState.game.winner != ''}
						<RematchButton
							rematchCount={$gameState.rematchCount}
							maxPlayers={$gameState.maxPlayers}
							onClick={handleRematch}
						/>
					{/if}
					<!-- Debug info temporal
					<details class="mt-4">
						<summary class="text-surface-400 cursor-pointer text-sm">Debug Info</summary>
						<pre class="bg-surface-900 mt-2 overflow-auto rounded p-2 text-xs">
         					{JSON.stringify($gameState, null, 2)}
                        </pre>
					</details>-->
				{:else}
					<p class="text-surface-400">No players yet...</p>
				{/if}
			</div>
		</div>
		<div class="w-full lg:w-4/5">
			{#if $gameState.players.length == 2}
				<Board
					board={$gameState.game.board}
					disabled={!myTurn || !!$gameState.game.winner}
					{onMove}
				/>
			{:else}
				<Board board={$gameState.game.board} disabled={true} {onMove} />
			{/if}
			<div class="mt-2">
				<GameStatus gameState={$gameState} {myTurn}></GameStatus>
			</div>
		</div>
	</div>
{:else}
	<div class="text-center">
		<h1 class="h3 lora-700 text-warning-400">Loading Room...</h1>
		<p class="text-surface-400 mb-4">Connecting to the room, please wait.</p>
		<p class="text-surface-500 text-sm">
			If this is taking too long, you can
			<a href="/play/tic-tac-toe" class="text-primary-400 hover:underline">return to lobby</a>.
		</p>
	</div>
{/if}

<audio src="/sounds/moveX.wav" bind:this={moveSoundX} preload="auto"></audio>
<audio src="/sounds/moveO.wav" bind:this={moveSoundO} preload="auto"></audio>
<audio src="/sounds/gameOver.wav" bind:this={gameOverSound} preload="auto"></audio>
