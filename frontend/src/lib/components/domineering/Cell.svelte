<script lang="ts">
	let {
		value = '',
		disabled,
		onClick
	} = $props<{
		value: 'H' | 'V' | '';
		disabled: boolean;
		onClick: () => void;
	}>();
</script>

<button
	onclick={onClick}
	{disabled}
	class="aspect-square w-full text-4xl font-bold disabled:opacity-100
	{disabled === true ? 'hover:opacity-100' : 'hover:opacity-50'}"
>
	{value}
</button>
