<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import type { GameConfig } from '$lib/config/games';
	import {
		isConnected,
		displayName,
		sendWebSocketMessage,
		players,
		gameState,
		roomClosedMessage,
		leftRoomData,
		playerLeftMessage,
		leaveRoom
	} from '$lib/socketStore';
	import RematchButton from '$lib/components/ui/RematchButton.svelte';
	import Collapsible from '$lib/components/ui/Collapsible.svelte';
	import GameChat from '$lib/components/game/GameChat.svelte';
	import type { Snippet } from 'svelte';

	interface Room {
		id: string;
		name: string;
	}

	interface Props {
		gameConfig: GameConfig;
		room: Room;
		error?: string;
		boardComponent: Snippet<[any]>;
		gameStatusComponent: Snippet<[any]>;
	}

	let { gameConfig, room, error, boardComponent, gameStatusComponent }: Props = $props();

	let hasJoined = $state(false);
	let moveSoundPlayer1: HTMLAudioElement;
	let moveSoundPlayer2: HTMLAudioElement;
	let gameOverSound: HTMLAudioElement;
	let prevTurn = $state<number | undefined>(undefined);
	let gameHasEnded = $state(false);

	// Sound effects logic
	onMount(() => {
		if (moveSoundPlayer1) moveSoundPlayer1.volume = 0.5;
		if (moveSoundPlayer2) moveSoundPlayer2.volume = 0.5;
		if (gameOverSound) gameOverSound.volume = 0.4;
	});

	// Initialize prevTurn
	$effect(() => {
		if ($gameState && prevTurn === undefined) {
			prevTurn = $gameState.game.currentTurn;
		}
	});

	// Handle turn changes and play sounds
	$effect(() => {
		if ($gameState && prevTurn !== undefined && $gameState.game.currentTurn !== prevTurn) {
			if (prevTurn === 0 && moveSoundPlayer1) {
				moveSoundPlayer1.currentTime = 0;
				moveSoundPlayer1.play().catch((e) => console.error('Error playing sound 1', e));
			} else if (prevTurn === 1 && moveSoundPlayer2) {
				moveSoundPlayer2.currentTime = 0;
				moveSoundPlayer2.play().catch((e) => console.error('Error playing sound 2', e));
			}
			prevTurn = $gameState.game.currentTurn;
		}
	});

	// Handle game over sound
	$effect(() => {
		if ($gameState && $gameState.game.winner && !gameHasEnded) {
			if (gameOverSound) {
				gameOverSound.currentTime = 0;
				gameOverSound.play().catch((e) => console.error('Error playing game over sound', e));
				gameHasEnded = true;
			}
		}
	});

	// Reset gameHasEnded when game resets
	$effect(() => {
		if ($gameState && $gameState.game.winner === '' && gameHasEnded) {
			gameHasEnded = false;
		}
	});

	// Room closed/left handling
	$effect(() => {
		if ($roomClosedMessage) {
			alert($roomClosedMessage);
			$roomClosedMessage = null;
			goto(`/play/${gameConfig.path}`);
		}
	});

	$effect(() => {
		if ($leftRoomData) {
			const gameType = $leftRoomData.gameType;
			$leftRoomData = null;
			const config =
				gameType === 'tic-tac-toe'
					? 'tic-tac-toe'
					: gameType === 'domineering'
						? 'domineering'
						: 'dots-and-boxes';
			goto(`/play/${config}`);
		}
	});

	$effect(() => {
		if ($playerLeftMessage) {
			alert($playerLeftMessage);
			$playerLeftMessage = null;
		}
	});

	// Join room
	$effect(() => {
		if ($isConnected && $displayName && !hasJoined && room) {
			hasJoined = true;
			sendWebSocketMessage({
				type: 'join_room',
				payload: {
					roomId: room.id
				}
			});
		}
	});

	let myTurn = $derived($gameState?.players[$gameState?.game.currentTurn] === $displayName);

	function handleRematch() {
		sendWebSocketMessage({ type: 'rematch_request' });
	}

	function handleLeaveRoom() {
		if (confirm('Are you sure you want to leave the room?')) {
			leaveRoom();
		}
	}
</script>

{#if error}
	<div class="text-center">
		<h1 class="text-error-400 mb-8 text-3xl">Room Not Found</h1>
		<a href={`/play/${gameConfig.path}`} class="bg-lime-400">Go to Lobby</a>
	</div>
{:else if $gameState}
	<div class="mb-4 flex items-center justify-between">
		<h1 class="text-3xl text-lime-400">{room.name}</h1>
		<button
			on:click={handleLeaveRoom}
			class="rounded-0 cursor-pointer border-r-2 border-b-2 border-rose-900 bg-rose-400 px-4 py-2 text-sm font-bold text-rose-950 transition-colors hover:bg-rose-500"
		>
			Leave Room
		</button>
	</div>

	<div class="flex flex-col justify-between gap-2 md:flex-row">
		<div class="w-full md:w-1/4">
			<Collapsible title="Players">
				{#if $gameState?.players?.length > 0}
					<div class="w-full border-b-1 border-stone-700 bg-stone-950">
						<div class="flex border-b-1 border-stone-700">
							<div
								class="w-10 border-r-1 border-stone-700 bg-stone-900 p-2 text-center font-bold text-lime-400"
							>
								<span>{gameConfig.playerSymbols[0]}</span>
							</div>
							<div class="w-full bg-stone-950 p-2">
								<span class="text-stone-200">{$gameState.players[0] || 'Waiting...'}</span>
							</div>
						</div>
						<div class="flex">
							<div
								class="w-10 border-r-1 border-stone-700 bg-stone-900 p-2 text-center font-bold text-rose-400"
							>
								<span>{gameConfig.playerSymbols[1]}</span>
							</div>
							<div class="w-full bg-stone-950 p-2">
								<span class="text-stone-200">{$gameState.players[1] || 'Waiting...'}</span>
							</div>
						</div>
					</div>
					{#if $gameState.game.winner != ''}
						<div class="mt-4 w-full">
							<RematchButton
								rematchCount={$gameState.rematchCount}
								maxPlayers={$gameState.maxPlayers}
								onClick={handleRematch}
							/>
						</div>
					{/if}
				{:else}
					<div class="w-full p-4 text-center">
						<p class="text-stone-400">No players yet...</p>
					</div>
				{/if}
			</Collapsible>

			{#if $gameState?.spectatorCount > 0}
				<div class="mt-4">
					<Collapsible
						title={`${$gameState.spectatorCount} Spectator${$gameState.spectatorCount == 1 ? '' : 's'}`}
					>
						<div class="w-full border-b-1 border-stone-700 bg-stone-950 p-1">
							<ul class="text-sm text-stone-200">
								{#each $gameState.spectators as spectator}
									<li class="py-1">{spectator}</li>
								{/each}
							</ul>
						</div>
					</Collapsible>
				</div>
			{/if}
		</div>

		<div class="w-full md:w-2/5">
			{@render boardComponent({
				gameState: $gameState,
				myTurn,
				disabled: !myTurn || !!$gameState.game.winner
			})}
			<div class="mt-2">
				{@render gameStatusComponent({ gameState: $gameState, myTurn })}
			</div>
		</div>
		<div class="w-full md:w-1/4">
			<div class="h-[600px]">
				<GameChat />
			</div>
		</div>
	</div>
{:else}
	<div class="text-center">
		<h1 class="lora-700 text-warning-400 text-3xl">Loading Room...</h1>
		<p class="mb-4 text-stone-400">Connecting to the room, please wait.</p>
		<p class="text-sm text-stone-500">
			If this is taking too long, you can
			<a href={`/play/${gameConfig.path}`} class="text-lime-400 hover:underline">return to lobby</a
			>.
		</p>
	</div>
{/if}

<audio src={gameConfig.sounds.move1} bind:this={moveSoundPlayer1} preload="auto"></audio>
<audio src={gameConfig.sounds.move2} bind:this={moveSoundPlayer2} preload="auto"></audio>
<audio src={gameConfig.sounds.gameOver} bind:this={gameOverSound} preload="auto"></audio>
