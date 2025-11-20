<script lang="ts">
	import { getGameConfig } from '$lib/config/games';
	import GameLobby from '$lib/components/game/GameLobby.svelte';

	const gameConfig = getGameConfig('dots-and-boxes');
</script>

<GameLobby {gameConfig} />
