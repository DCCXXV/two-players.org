<script lang="ts">
	import '../app.css';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import {
		connectWebSocket,
		disconnectWebSocket,
		isConnected,
		displayName,
		gameState
	} from '$lib/socketStore';

	let menuOpen = $state(false);

	onMount(() => {
		console.log('Layout onMount: Attempting to call connectWebSocket...');
		connectWebSocket();

		const handleBeforeUnload = (event: BeforeUnloadEvent) => {
			// Only show confirmation if the user is in a game room.
			if ($gameState && page.url.pathname.includes('/play/tic-tac-toe/')) {
				event.preventDefault();
				event.returnValue = 'Are you sure you want to leave the room?';
			}
		};

		window.addEventListener('beforeunload', handleBeforeUnload);

		return () => {
			console.log('Layout onDestroy: Calling disconnectWebSocket');
			window.removeEventListener('beforeunload', handleBeforeUnload);
			disconnectWebSocket();
		};
	});

	let { children } = $props();
</script>

<div class="grid h-screen w-full grid-rows-[auto_1fr_auto]">
	<header
		class="border-surface-500 xl-landscape:px-16 flex min-h-14 flex-wrap justify-between gap-4 border-b-2 px-4 py-2 sm:justify-between sm:py-0"
	>
		<div class="flex w-full justify-between sm:w-auto">
			<a class="text-surface-500 lora-700 text-2xl font-bold hover:brightness-80" href="/">
				<span class="text-primary-400">two</span><span class="text-surface-500">-</span><span
					class="text-secondary-400">players</span
				>.org
			</a>
			<div class="flex items-center sm:hidden">
				{#if $isConnected && $displayName}
					<input
						type="text"
						readonly
						class="input text-primary-400 bg-surface-900 lora-400 mr-4 max-w-35 rounded-none text-center text-lg"
						value={$displayName}
					/>
				{:else if $isConnected}
					<span class="text-primary-400 mr-4 text-xl">Loading...</span>
				{/if}
				<button class="text-surface-200 text-3xl" onclick={() => (menuOpen = !menuOpen)}>
					&#9776;
				</button>
			</div>
		</div>
		<nav
			class="text-surface-200 bg-surface-950 z-10 ms-4 me-auto w-full flex-col gap-2 sm:flex sm:w-auto sm:flex-row"
			class:hidden={!menuOpen}
		>
			<section
				class="bg-surface-950 flex h-full flex-col justify-center px-2"
				class:border-b-3={page.url.pathname.startsWith('/play')}
				class:border-primary-400={page.url.pathname.startsWith('/play')}
			>
				<a
					class="hover:text-primary-100 transition-colors duration-200"
					class:text-primary-400={page.url.pathname.startsWith('/play')}
					href="/play"
				>
					PLAY
				</a>
			</section>
			<section
				class="bg-surface-950 flex h-full flex-col justify-center px-2"
				class:border-b-3={page.url.pathname === '/explore'}
				class:border-primary-400={page.url.pathname === '/explore'}
			>
				<a
					class="hover:text-primary-100 transition-colors duration-200"
					class:text-primary-400={page.url.pathname === '/explore'}
					href="/play"
				>
					EXPLORE
				</a>
			</section>
			<section
				class="bg-surface-950 flex h-full flex-col justify-center px-2"
				class:border-b-3={page.url.pathname === '/learn'}
				class:border-primary-400={page.url.pathname === '/learn'}
			>
				<a
					class="hover:text-primary-100 transition-colors duration-200"
					class:text-primary-400={page.url.pathname === '/learn'}
					href="/play"
				>
					LEARN
				</a>
			</section>
		</nav>
		<div class="hidden gap-6 sm:flex">
			{#if $isConnected && $displayName}
				<input
					type="text"
					readonly
					class="input text-primary-400 bg-surface-900 lora-400 my-auto h-10 max-w-35 rounded-none text-center text-lg"
					value={$displayName}
				/>
			{:else if $isConnected}
				<span class="text-primary-400 text-xl">Loading...</span>
			{/if}
		</div>
	</header>
	<div class="xl-landscape:px-16 mt-4 w-full px-4">
		<main class="w-full">
			{@render children()}
		</main>
	</div>
	<footer class="text-surface-400 bg-surface-900 mt-4 p-8 text-center">
		<nav class="flex justify-center space-x-4">
			<a href="/privacy-policy" class="hover:text-primary-400 transition-colors duration-200"
				>Privacy Policy</a
			>
			<a href="/terms-of-service" class="hover:text-primary-400 transition-colors duration-200"
				>Terms of Service</a
			>
			<a
				href="https://github.com/DCCXXV/two-players.org"
				target="_blank"
				rel="noopener noreferrer"
				class="hover:text-primary-400 transition-colors duration-200">GitHub</a
			>
			<a href="/contact" class="hover:text-primary-400 transition-colors duration-200">Contact</a>
		</nav>
		<p class="mt-4 text-sm">
			&copy; {new Date().getFullYear()} two-players.org. All rights reserved.
		</p>
	</footer>
</div>

<style>
</style>
