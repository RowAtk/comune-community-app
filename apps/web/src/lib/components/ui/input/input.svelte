<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		value = undefined,
		type = 'text',
		name = undefined,
		id = undefined,
		placeholder = undefined,
		required = false,
		disabled = false,
		min = undefined
	}: {
		class?: string;
		value?: string | number;
		type?: string;
		name?: string;
		id?: string;
		placeholder?: string;
		required?: boolean;
		disabled?: boolean;
		min?: string;
	} = $props();
</script>

<input
	class={cn(
		'flex h-11 w-full rounded-xl border border-[var(--color-app-border)] bg-white px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm outline-none transition placeholder:text-[var(--color-app-subtle)] focus:border-[var(--color-secondary-400)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-secondary-300)_35%,transparent)] disabled:cursor-not-allowed disabled:bg-slate-50',
		className
	)}
	{value}
	{type}
	{name}
	{id}
	{placeholder}
	{required}
	{disabled}
	{min}
/>
