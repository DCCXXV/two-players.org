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
	class="blackboard-darker erasdust-400 aspect-square w-full text-7xl font-bold disabled:opacity-100
	{value === 'X' ? 'text-primary-300' : 'text-secondary-300'}
	{disabled === true ? 'hover:opacity-100' : 'hover:opacity-50'}"
>
	{value}
</button>
