<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		for: htmlFor = undefined,
		class: className = '',
		children
	}: {
		for?: string;
		class?: string;
		children?: import('svelte').Snippet;
	} = $props();
</script>

<label class={cn('text-sm font-medium text-[var(--color-app-text)]', className)} for={htmlFor}>
	{@render children?.()}
</label>
