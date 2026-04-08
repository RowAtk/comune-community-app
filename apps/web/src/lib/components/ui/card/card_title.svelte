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

<h2 class={cn('text-lg font-semibold tracking-tight text-[var(--color-app-text)]', className)}>
	{@render children?.()}
</h2>
