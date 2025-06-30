<script lang="ts">
	// Temp
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	const targetUrl = '/play/tic-tac-toe';

	onMount(() => {
		goto(targetUrl, { replaceState: true });
	});
</script>

<section class="w-full">
	<ol class="flex items-center gap-4">
		<li><a class="text-primary-400 hover:underline" href="/play">Play</a></li>
	</ol>
</section>
