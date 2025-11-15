<script lang="ts">
	import { getGameConfig } from '$lib/config/games';
	import GameLobby from '$lib/components/game/GameLobby.svelte';

	const gameConfig = getGameConfig('tic-tac-toe');
</script>

<GameLobby {gameConfig} />
