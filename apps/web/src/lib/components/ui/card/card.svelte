<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		class: className = '',
		children
	}: {
		class?: string;
		children?: import('svelte').Snippet;
	} = $props();
</script>

<section
	class={cn(
		'rounded-[calc(var(--radius-2xl)+2px)] border border-[var(--color-app-border)] bg-white shadow-[var(--shadow-card)]',
		className
	)}
>
	{@render children?.()}
</section>
