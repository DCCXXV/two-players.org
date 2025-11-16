<script lang="ts">
	let {
		value = '',
		disabled,
		onClick
	} = $props<{
		value: 'X' | 'O' | '';
		disabled: boolean;
		onClick: () => void;
	}>();
</script>

<button
	onclick={onClick}
	{disabled}
	class="bg-zinc-950 aspect-square w-full border-2 border-zinc-700 text-8xl disabled:opacity-100
	{value === 'X' ? 'text-blue-400' : 'text-red-400'}
	{disabled === true ? 'cursor-not-allowed' : 'hover:border-zinc-500'}"
>
	{value}
</button>
