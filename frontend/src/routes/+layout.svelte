<script lang="ts">
	import '../app.css';
	import { page as pageStore } from '$app/state';
	import { onMount } from 'svelte';
	import {
		connectWebSocket,
		disconnectWebSocket,
		isConnected,
		displayName,
		gameState,
		updateDisplayName
	} from '$lib/socketStore';

	let menuOpen = $state(false);
	let showNameModal = $state(false);
	let tempName = $state('');

	function openNameModal() {
		tempName = $displayName;
		showNameModal = true;
	}

	function closeNameModal() {
		showNameModal = false;
	}

	function saveName() {
		if (tempName.trim() && tempName !== $displayName) {
			updateDisplayName(tempName.trim());
		}
		closeNameModal();
	}

	function handleModalKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			closeNameModal();
		} else if (e.key === 'Enter') {
			saveName();
		}
	}

	onMount(() => {
		console.log('Layout onMount: Attempting to call connectWebSocket...');
		connectWebSocket();

		const handleBeforeUnload = (event: BeforeUnloadEvent) => {
			if ($gameState && pageStore.url.pathname.includes('/play/tic-tac-toe/')) {
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

	let { children, data } = $props();

	let meta = $derived(pageStore.data.meta || data.meta);
</script>

<svelte:head>
	<title>{meta.title}</title>
	<meta name="description" content={meta.description} />

	<meta property="og:type" content="website" />
	<meta property="og:url" content={meta.url} />
	<meta property="og:title" content={meta.title} />
	<meta property="og:description" content={meta.description} />
	<meta property="og:image" content={meta.imageUrl} />

	<meta property="twitter:card" content="summary_large_image" />
	<meta property="twitter:url" content={meta.url} />
	<meta property="twitter:title" content={meta.title} />
	<meta property="twitter:description" content={meta.description} />
	<meta property="twitter:image" content={meta.imageUrl} />
</svelte:head>

<div class="titillium-web-light grid h-screen w-full grid-rows-[auto_1fr_auto]">
	<header
		class="xl-landscape:px-16 relative flex min-h-12 flex-wrap justify-between gap-4 bg-linear-to-t from-zinc-900 to-zinc-800 px-4 py-2 sm:justify-between sm:py-0"
	>
		<div class="flex w-full justify-between sm:w-auto">
			<a
				class="lora-700 text-md me-6 pt-2 font-bold text-zinc-500 hover:brightness-80 md:text-2xl"
				href="/"
			>
				<span class="text-blue-400">two</span><span class="text-zinc-500">-</span><span
					class="text-red-400">players</span
				>.org
			</a>
			<div class="flex items-center sm:hidden">
				{#if $isConnected && $displayName}
					<button
						onclick={openNameModal}
						class="lora-400 my-auto h-10 max-w-35 rounded-none border-0 bg-transparent text-center text-lg text-zinc-300 hover:text-blue-400 focus:ring-0"
					>
						{$displayName}
					</button>
				{:else if $isConnected}
					<span class="mr-4 text-xl text-blue-400">Loading...</span>
				{/if}
				<button class="text-2xl text-zinc-200" onclick={() => (menuOpen = !menuOpen)}>
					&#9776;
				</button>
			</div>
		</div>
		<nav
			class="text-md absolute top-full left-0 z-10 w-full flex-col gap-2 border-2 border-zinc-500 text-zinc-200 opacity-90 sm:static sm:ms-4 sm:me-auto sm:flex sm:w-auto sm:flex-row sm:border-none"
			class:hidden={!menuOpen}
		>
			<section
				class="py-1.5-2 flex h-full flex-col justify-center px-8"
				class:border-b-1={pageStore.url.pathname === '/'}
				class:border-blue-400={pageStore.url.pathname === '/'}
			>
				<a
					class="transition-colors duration-200 hover:text-blue-400"
					class:text-blue-400={pageStore.url.pathname === '/'}
					href="/"
				>
					ABOUT
				</a>
			</section>
			<section
				class="py-1.5-2 flex h-full flex-col justify-center px-8"
				class:border-b-1={pageStore.url.pathname.startsWith('/play')}
				class:border-blue-400={pageStore.url.pathname.startsWith('/play')}
			>
				<a
					class="transition-colors duration-200 hover:text-blue-400"
					class:text-blue-400={pageStore.url.pathname.startsWith('/play')}
					href="/play"
				>
					PLAY
				</a>
			</section>
		</nav>
		<div class="hidden gap-6 sm:flex">
			{#if $isConnected && $displayName}
				<button
					onclick={openNameModal}
					class="my-auto h-10 max-w-35 cursor-pointer rounded-none border-0 bg-transparent text-center text-lg text-zinc-300 hover:text-blue-400 focus:ring-0"
				>
					{$displayName}
				</button>
			{:else if $isConnected}
				<span class="text-xl text-blue-400">Loading...</span>
			{/if}
		</div>
	</header>
	<div class="xl-landscape:px-16 mt-4 w-full bg-zinc-900 px-4">
		<main class="w-full">
			{@render children()}
		</main>
	</div>
	<footer class="mt-4 bg-zinc-900 p-4 text-center text-zinc-400 sm:p-8">
		<nav class="flex flex-wrap justify-center gap-2 sm:gap-4">
			<a href="/privacy-policy" class="transition-colors duration-200 hover:text-blue-400"
				>Privacy Policy</a
			>
			<a href="/terms-of-service" class="transition-colors duration-200 hover:text-blue-400"
				>Terms of Service</a
			>
			<a href="/contact" class="transition-colors duration-200 hover:text-blue-400">Contact</a>
			<a
				href="https://discord.gg/8s9NneBy"
				target="_blank"
				rel="noopener noreferrer"
				class="transition-colors duration-200 hover:text-blue-400">Discord</a
			>
			<a
				href="https://github.com/DCCXXV/two-players.org"
				target="_blank"
				rel="noopener noreferrer"
				class="transition-colors duration-200 hover:text-blue-400">GitHub</a
			>
		</nav>
		<p class="mt-4 text-sm">
			&copy; {new Date().getFullYear()} two-players.org. All rights reserved.
		</p>
	</footer>
</div>

{#if showNameModal}
	<div
		class="titillium-web-light fixed inset-0 z-50 flex items-center justify-center bg-zinc-950/80"
		onclick={closeNameModal}
		onkeydown={handleModalKeydown}
		role="button"
		tabindex="-1"
	>
		<div
			class="rounded-0 w-full max-w-md border-b-1 border-zinc-700 bg-zinc-900 p-6"
			onclick={(e) => e.stopPropagation()}
			role="dialog"
			tabindex="-1"
		>
			<h2 class="mb-4 text-center text-2xl text-blue-400">Change Display Name</h2>
			<input
				type="text"
				bind:value={tempName}
				class="mb-6 w-full rounded-none border-0 border-b-1 border-zinc-700 bg-zinc-900 px-4 py-2 text-center text-xl text-zinc-300 focus:border-blue-400 focus:ring-0"
				placeholder="Enter new name"
				autofocus
			/>
			<div class="flex justify-center gap-4">
				<button
					onclick={closeNameModal}
					class="cursor-pointer rounded-none border-b-1 border-red-400 bg-transparent px-6 py-2 text-lg text-red-400 hover:bg-red-400 hover:text-zinc-950"
				>
					Cancel
				</button>
				<button
					onclick={saveName}
					class="cursor-pointer rounded-none border-b-1 border-blue-400 bg-transparent px-6 py-2 text-lg text-blue-400 hover:bg-blue-400 hover:text-zinc-950"
				>
					Accept
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
</style>
