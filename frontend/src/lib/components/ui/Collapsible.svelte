<script lang="ts">
	import type { Snippet } from 'svelte';
	import { ChevronUp } from 'lucide-svelte';
	import { slide } from 'svelte/transition';

	let { title, children }: { title: string; children?: Snippet } = $props();
	let isOpen = $state(true);

	function toggle() {
		isOpen = !isOpen;
	}
</script>

<section class="mb-8">
	<button
		onclick={toggle}
		class="flex w-full cursor-pointer items-center justify-between border-b-1 border-zinc-700 bg-zinc-800 px-3 py-2 text-left text-2xl text-zinc-300 transition-colors hover:bg-zinc-700"
	>
		<span
			><span class="text-blue-400">#</span> {title} <span class="text-red-400">#</span></span
		>
		<span
			class="text-lg transition-transform duration-75"
			style:transform={isOpen ? 'rotate(-180deg)' : 'rotate(0deg)'}
		>
			<ChevronUp />
		</span>
	</button>
	{#if isOpen}
		<div
			transition:slide={{ duration: 75 }}
			class="flex flex-wrap overflow-hidden bg-zinc-900 text-center"
		>
			{#if children}
				{@render children()}
			{/if}
		</div>
	{/if}
</section>
