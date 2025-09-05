<script lang="ts">
	import InfoBadge from '$lib/components/ui/badges/InfoBadge.svelte';
	let {
		title,
		path,
		src,
		desc,
		tags
	}: {
		title: string;
		path: string;
		src: string;
		desc: string;
		tags: string[];
	} = $props();
</script>

<a href="/play/{path}">
	<div
		class="card group preset-filled-surface-100-900 hover:border-primary-400 relative flex aspect-[2/1] rounded-none border-2 border-transparent shadow-sm transition-all duration-200 ease-in-out hover:border-2"
	>
		<img class="ms-4 aspect-square w-1/2 object-contain" {src} alt="Picture of a {title} match" />
		<div class="m-6 flex w-1/2 flex-col">
			<h4 class="h4 text-surface-100 lora-700 mb-6">{title}</h4>
			<div class="flex flex-grow flex-col overflow-auto">
				<p class="mb-6 opacity-60">
					{desc}
				</p>
				<div class="mt-auto flex flex-wrap gap-2">
					<InfoBadge label="Board" value={tags[0]} />
					<InfoBadge label="Class" value={tags[1]} />
					<InfoBadge label="Solved" value={tags[2]} />
					<InfoBadge label="Outcome" value={tags[3]} />
				</div>
			</div>
		</div>
	</div>
</a>
