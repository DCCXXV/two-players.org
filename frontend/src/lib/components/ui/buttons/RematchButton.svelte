<script lang="ts">
	export let rematchCount = 0;
	export let maxPlayers = 2;
	export let onClick: () => void = () => {};

	$: buttonText = `Rematch? ${rematchCount}/${maxPlayers}`;
</script>

<button
	id="rematch-button"
	class="border-surface-400 bg-surface-900 text-surface-200 lora-700 w-full border-2 p-2 text-xl"
	on:click={onClick}
>
	{buttonText}
	<div class="ribbon-left"></div>
	<div class="ribbon-right"></div>
</button>

<style>
	#rematch-button {
		cursor: pointer;
		position: relative;
		display: inline-block;
		outline: none;
	}
	#rematch-button::before {
		content: '';
		display: block;
		position: absolute;
		background-color: var(--color-secondary-400);
		top: 0;
		left: 1rem;
		width: 1rem;
		height: 100%;
	}
	#rematch-button::after {
		content: '';
		display: block;
		position: absolute;
		background-color: var(--color-secondary-400);
		top: 0;
		right: 1rem;
		width: 1rem;
		height: 100%;
	}
	.ribbon-left,
	.ribbon-right {
		position: absolute;
		width: 1rem;
		height: 4rem;
		background-color: var(--color-secondary-400);
		top: 100%;
		filter: brightness(50%);
		clip-path: polygon(0 0, 100% 0, 100% 100%, 50% 80%, 0 100%);
		transform-origin: center top;
	}
	.ribbon-left {
		left: calc(1rem + 4px);
	}
	.ribbon-right {
		right: calc(1rem + 4px);
	}
	#rematch-button:active {
		transform: translateY(2px);
	}

	#rematch-button:active .ribbon-left {
		animation: ribbon-sway-left 0.6s ease-in-out infinite alternate;
	}

	#rematch-button:active .ribbon-right {
		animation: ribbon-sway-right 0.6s ease-in-out infinite alternate;
	}

	@keyframes ribbon-sway-left {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(-8deg);
		}
	}

	@keyframes ribbon-sway-right {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(8deg);
		}
	}
</style>
