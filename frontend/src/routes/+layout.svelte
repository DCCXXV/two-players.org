<script lang="ts">
	import '../app.css';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import {
		connectWebSocket,
		disconnectWebSocket,
		isConnected,
		displayName
	} from '$lib/socketStore';

onMount(() => {
		console.log('Layout onMount: Attempting to call connectWebSocket...');
		connectWebSocket();

		return () => {
			console.log('Layout onDestroy: Calling disconnectWebSocket');
			disconnectWebSocket();
		};
	});

	let { children } = $props();
</script>

<div class="grid h-screen w-full grid-rows-[auto_1fr_auto]">
	<header
		class="border-surface-500 xl-landscape:px-16 flex min-h-14 items-center justify-between border-b-2 px-4"
	>
		<div class="flex h-full items-center">
			<a
				class="text-surface-500 hover:brightness-80 lora-700 me-6 pb-2 text-2xl font-bold"
				href="/"
			>
				<span class="text-primary-400">two</span><span class="text-surface-500">-</span><span
					class="text-secondary-400">players</span
				>.org
			</a>
			<nav class="text-surface-200 flex h-full justify-center gap-2">
				<section
					class="flex h-full flex-col justify-center px-2"
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
					class="flex h-full flex-col justify-center px-2"
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
					class="flex h-full flex-col justify-center px-2"
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
		</div>
		<div class="flex gap-6">
			<!--<button class="btn preset-filled-surface-500 text-surface-100"> Change Display Name </button>-->
			{#if $isConnected && $displayName}
				<input type="text" readonly class="input text-lg text-primary-400 bg-surface-900 max-w-35 text-center lora-400 rounded-none" value={$displayName}>
			{:else if $isConnected}
				<span class="text-xl text-primary-400">Loading...</span>
			{/if}
			<!--
			<a aria-label="github repository" href="https://github.com/DCCXXV/two-players.org">
				<svg
					class="h-7 w-7 transform transition duration-200 ease-in-out hover:-translate-y-1 pt-1"
					viewBox="0 0 1024 1024"
					fill="none"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path
						fill-rule="evenodd"
						clip-rule="evenodd"
						d="M8 0C3.58 0 0 3.58 0 8C0 11.54 2.29 14.53 5.47 15.59C5.87 15.66 6.02 15.42 6.02 15.21C6.02 15.02 6.01 14.39 6.01 13.72C4 14.09 3.48 13.23 3.32 12.78C3.23 12.55 2.84 11.84 2.5 11.65C2.22 11.5 1.82 11.13 2.49 11.12C3.12 11.11 3.57 11.7 3.72 11.94C4.44 13.15 5.59 12.81 6.05 12.6C6.12 12.08 6.33 11.73 6.56 11.53C4.78 11.33 2.92 10.64 2.92 7.58C2.92 6.71 3.23 5.99 3.74 5.43C3.66 5.23 3.38 4.41 3.82 3.31C3.82 3.31 4.49 3.10 6.02 4.13C6.66 3.95 7.34 3.86 8.02 3.86C8.70 3.86 9.38 3.95 10.02 4.13C11.55 3.09 12.22 3.31 12.22 3.31C12.66 4.41 12.38 5.23 12.30 5.43C12.81 5.99 13.12 6.70 13.12 7.58C13.12 10.65 11.25 11.33 9.47 11.53C9.76 11.78 10.01 12.26 10.01 13.01C10.01 14.08 10 14.94 10 15.21C10 15.42 10.15 15.67 10.55 15.59C13.71 14.53 16 11.53 16 8C16 3.58 12.42 0 8 0Z"
						transform="scale(64)"
						fill="#FFFFFF99"
					/>
				</svg>
			</a>
			-->
		</div>
	</header>
	<div class="xl-landscape:px-16 mt-4 w-full px-4">
		<main class="w-full">
			{@render children()}
		</main>
	</div>
	<footer class="p-8"></footer>
</div>

<style>
</style>
