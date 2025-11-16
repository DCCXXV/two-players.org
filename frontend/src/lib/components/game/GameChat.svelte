<script lang="ts">
	import { chatMessages, sendWebSocketMessage, displayName } from '$lib/socketStore';
	import { Send } from 'lucide-svelte';
	import { onMount, tick } from 'svelte';

	let messageInput = $state('');
	let chatContainer: HTMLDivElement;

	function sendMessage() {
		const trimmed = messageInput.trim();
		if (!trimmed) return;

		sendWebSocketMessage({
			type: 'chat_message',
			payload: {
				message: trimmed
			}
		});

		messageInput = '';
	}

	function handleKeyPress(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			sendMessage();
		}
	}

	$effect(() => {
		if ($chatMessages && chatContainer) {
			tick().then(() => {
				chatContainer.scrollTop = chatContainer.scrollHeight;
			});
		}
	});

	function formatTime(timestamp: string): string {
		const date = new Date(timestamp);
		return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
	}
</script>

<div class="flex h-full flex-col">
	<div
		class="flex items-center justify-between border-b-1 border-zinc-700 bg-zinc-800 px-3 py-2"
	>
		<span class="text-2xl text-zinc-400">
			<span class="text-blue-400">#</span> Chat <span class="text-red-400">#</span>
		</span>
	</div>

	<div
		bind:this={chatContainer}
		class="flex-1 overflow-y-auto border-zinc-700 bg-zinc-900 p-1"
	>
		{#if $chatMessages.length === 0}
			<div class="flex h-full items-center justify-center">
				<p class="text-sm text-zinc-500">No messages yet...</p>
			</div>
		{:else}
			<div class="flex flex-col gap-2">
				{#each $chatMessages as message}
					<div
						class="rounded-0 border-b-1 border-zinc-700 bg-zinc-900 p-2 {message.displayName ===
						$displayName
							? 'border-l-2 border-l-blue-400'
							: 'border-l-2 border-l-zinc-700'}"
					>
						<div class="mb-1 flex items-baseline justify-between gap-2">
							<span
								class="text-sm font-bold {message.displayName === $displayName
									? 'text-blue-400'
									: 'text-zinc-300'}"
							>
								{message.displayName === $displayName ? 'You' : message.displayName}
							</span>
							<span class="text-xs text-zinc-500">{formatTime(message.timestamp)}</span>
						</div>
						<p class="break-words text-sm text-zinc-200">{message.message}</p>
					</div>
				{/each}
			</div>
		{/if}
	</div>

	<div class="border-b-1 border-zinc-700 bg-zinc-800 p-2">
		<div class="flex gap-2">
			<input
				type="text"
				bind:value={messageInput}
				onkeypress={handleKeyPress}
				placeholder="Type a message..."
				class="flex-1 rounded-0 border-0 border-b-1 border-zinc-700 bg-zinc-900 px-3 py-2 text-sm text-zinc-200 placeholder-zinc-600 focus:border-blue-400 focus:ring-0"
			/>
			<button
				onclick={sendMessage}
				disabled={!messageInput.trim()}
				class="rounded-0 flex items-center gap-1 border-b-2 border-r-2 border-blue-900 bg-blue-400 px-3 py-2 text-sm font-bold text-zinc-950 transition-colors hover:bg-blue-500 disabled:border-zinc-800 disabled:bg-zinc-800 disabled:text-zinc-600"
			>
				<Send size={16} />
			</button>
		</div>
	</div>
</div>
