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
	class="aspect-square w-full border-1 border-r-2 border-b-2 border-stone-700 text-8xl disabled:opacity-100
	{value === 'X' ? 'text-lime-400' : 'text-rose-400'}
	{disabled === true ? 'cursor-not-allowed' : 'hover:border-stone-500'}"
>
	{value}
</button>
