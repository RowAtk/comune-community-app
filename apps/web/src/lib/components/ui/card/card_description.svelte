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

<p class={cn('text-sm leading-6 text-[var(--color-app-muted)]', className)}>
	{@render children?.()}
</p>
