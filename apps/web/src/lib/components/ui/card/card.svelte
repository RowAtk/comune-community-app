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

<section class={cn('theme-panel rounded-[calc(var(--radius-2xl)+2px)] border', className)}>
	{@render children?.()}
</section>
