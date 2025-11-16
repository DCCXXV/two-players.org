<script lang="ts">
	export let rematchCount = 0;
	export let maxPlayers = 2;
	export let onClick: () => void = () => {};

	$: buttonText = `Rematch? ${rematchCount}/${maxPlayers}`;

	function playClickSound() {
		const audio = new Audio('/sounds/click.wav');
		audio.volume = 0.5;
		audio.play();
	}
</script>

<button
	id="rematch-button"
	class="w-full cursor-pointer border-r-2 border-b-2 border-blue-900 bg-blue-500 p-2 text-xl font-bold text-blue-950"
	on:click={() => {
		onClick();
		playClickSound();
	}}
>
	{buttonText}
	<div class="ribbon-left"></div>
	<div class="ribbon-right"></div>
</button>
