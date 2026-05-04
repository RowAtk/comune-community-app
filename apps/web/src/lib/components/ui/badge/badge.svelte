<script lang="ts">
	import { cn } from '$lib/utils';

	type Variant = 'default' | 'secondary' | 'outline' | 'success' | 'warning';

	let {
		variant = 'default',
		class: className = '',
		children
	}: {
		variant?: Variant;
		class?: string;
		children?: import('svelte').Snippet;
	} = $props();

	const variants: Record<Variant, string> = {
		default: 'bg-[var(--color-brand-100)] text-[var(--color-brand-800)]',
		secondary: 'bg-[var(--color-secondary-100)] text-[var(--color-secondary-800)]',
		outline:
			'border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] text-[var(--color-app-muted)]',
		success: 'bg-[var(--color-success-100)] text-[var(--color-success-800)]',
		warning: 'bg-[var(--color-warning-100)] text-[var(--color-warning-800)]'
	};
</script>

<span
	class={cn(
		'inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium',
		variants[variant],
		className
	)}
>
	{@render children?.()}
</span>
