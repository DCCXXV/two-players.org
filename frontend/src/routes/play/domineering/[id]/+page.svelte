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
	import RematchButton from '$lib/components/ui/RematchButton.svelte';
	import Board from '$lib/components/domineering/Board.svelte';
	import GameStatus from '$lib/components/domineering/GameStatus.svelte';
	import Collapsible from '$lib/components/ui/Collapsible.svelte';

	import type { PageData } from './$types';

	export let data: PageData;

	let hasJoined = false;

	let moveSoundH: HTMLAudioElement;
	let moveSoundV: HTMLAudioElement;
	let gameOverSound: HTMLAudioElement;

	let prevTurn: number | undefined = undefined;
	let gameHasEnded = false;

	onMount(() => {
		if (moveSoundH) {
			moveSoundH.volume = 0.5;
		}
		if (moveSoundV) {
			moveSoundV.volume = 0.5;
		}
		if (gameOverSound) {
			gameOverSound.volume = 0.4;
		}
	});

	$: if ($gameState && prevTurn === undefined) {
		prevTurn = $gameState.game.currentTurn;
	}

	$: if ($gameState && prevTurn !== undefined && $gameState.game.currentTurn !== prevTurn) {
		if (prevTurn === 0 && moveSoundH) {
			moveSoundH.currentTime = 0;
			moveSoundH.play().catch((e) => console.error('Error playing sound H', e));
		} else if (prevTurn === 1 && moveSoundV) {
			moveSoundV.currentTime = 0;
			moveSoundV.play().catch((e) => console.error('Error playing sound V', e));
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
		alert($roomClosedMessage);
		$roomClosedMessage = null;
		goto('/play/domineering');
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
	$: mySymbol =
		$gameState?.players[0] === $displayName
			? 'H'
			: $gameState?.players[1] === $displayName
				? 'V'
				: undefined;

	function onMove(row: number, col: number) {
		sendWebSocketMessage({
			type: 'make_move',
			payload: {
				row: row,
				col: col
			}
		});
	}

	function handleRematch() {
		sendWebSocketMessage({ type: 'rematch_request' });
	}
</script>

{#if data.error}
	<div class="text-center">
		<h1 class="text-error-400 mb-8 text-3xl">Room Not Found</h1>
		<a href="/play/domineering" class="bg-lime-400">Go to Lobby</a>
	</div>
{:else if $gameState}
	<h1 class="mb-4 text-3xl text-lime-400">{data.room.name}</h1>

	<div class="flex flex-col justify-between gap-4 md:flex-row">
		<div class="w-full md:w-1/5">
			<Collapsible title="Players">
				{#if $gameState?.players?.length > 0}
					<div class="w-full border-r-2 border-b-2 border-stone-700 bg-stone-950">
						<div class="flex border-b-1 border-stone-700">
							<div
								class="w-10 border-r-1 border-stone-700 bg-stone-900 p-2 text-center font-bold text-lime-400"
							>
								<span>H</span>
							</div>
							<div class="w-full bg-stone-950 p-2">
								<span class="text-stone-200">{$gameState.players[0] || 'Waiting...'}</span>
							</div>
						</div>
						<div class="flex">
							<div
								class="w-10 border-r-1 border-stone-700 bg-stone-900 p-2 text-center font-bold text-rose-400"
							>
								<span>V</span>
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

		<div class="w-full md:w-4/5">
			{#if $gameState.players.length == 2}
				<Board
					board={$gameState.game.board}
					disabled={!myTurn || !!$gameState.game.winner}
					{onMove}
					{mySymbol}
				/>
			{:else}
				<Board board={$gameState.game.board} disabled={true} {onMove} />
			{/if}
			<div class="mt-2">
				<GameStatus gameState={$gameState} {myTurn} />
			</div>
		</div>
		<div class="w-full md:w-1/5"></div>
	</div>
{:else}
	<div class="text-center">
		<h1 class="lora-700 h3 text-stone-400">Loading Room...</h1>
		<p class="mb-4 text-stone-400">Connecting to the room, please wait.</p>
		<p class="text-sm text-stone-500">
			If this is taking too long, you can
			<a href="/play/domineering" class="text-lime-400 hover:underline">return to lobby</a>.
		</p>
	</div>
{/if}

<audio src="/sounds/moveX.wav" bind:this={moveSoundH} preload="auto"></audio>
<audio src="/sounds/moveO.wav" bind:this={moveSoundV} preload="auto"></audio>
<audio src="/sounds/gameOver.wav" bind:this={gameOverSound} preload="auto"></audio>
