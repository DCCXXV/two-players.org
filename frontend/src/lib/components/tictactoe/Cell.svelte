<script lang="ts">
	let { value, disabled, onClick } = $props<{
		value: 'X' | 'O' | '';
		disabled: boolean;
		onClick: () => void;
	}>();
</script>

<button
	onclick={onClick}
	{disabled}
	class="bg-surface-600 hover:bg-surface-500 aspect-square w-full text-5xl font-bold transition-colors disabled:opacity-50
	{value === 'X' ? 'text-primary-400' : 'text-secondary-400'}"
>
	{value}
</button>
