<script lang="ts">
	let {
		title,
		path,
	}: {
		title: string;
		path: string;
	} = $props();
</script>

<a href="/play/{path}">
	<div
		class="bg-stone-950 p-1 hover:px-7 hover:border-lime-400 border-b-1 border-stone-700 rounded-0 transition-all duration-150 ease-in-out"
	>
		<div class="m-6 text-center">
			<h2 class="h4 text-surface-100">{title}</h2>
		</div>
	</div>
</a>
