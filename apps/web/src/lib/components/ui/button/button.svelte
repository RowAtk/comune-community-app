<script lang="ts">
	import { resolve } from '$app/paths';
	import { cn } from '$lib/utils';

	type Variant = 'default' | 'secondary' | 'outline' | 'ghost' | 'destructive';
	type Size = 'default' | 'sm' | 'lg';

	let {
		type = 'button',
		variant = 'default',
		size = 'default',
		href = undefined,
		onclick = undefined,
		formAction = undefined,
		class: className = '',
		disabled = false,
		children
	}: {
		type?: 'button' | 'submit' | 'reset';
		variant?: Variant;
		size?: Size;
		href?: string;
		onclick?: ((event: MouseEvent) => void) | undefined;
		formAction?: string;
		class?: string;
		disabled?: boolean;
		children?: import('svelte').Snippet;
	} = $props();

	const base =
		'inline-flex shrink-0 items-center justify-center gap-2 whitespace-nowrap rounded-xl text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50';

	const variants: Record<Variant, string> = {
		default:
			'bg-[var(--color-app-accent)] text-[var(--color-app-inverse)] shadow-sm hover:bg-[color-mix(in_oklab,var(--color-app-accent)_86%,black)] focus-visible:ring-[var(--color-brand-300)]',
		secondary:
			'bg-[var(--color-secondary-600)] text-[var(--color-app-inverse)] shadow-sm hover:bg-[var(--color-secondary-700)] focus-visible:ring-[var(--color-secondary-400)]',
		outline:
			'border border-[var(--color-app-border-strong)] bg-[var(--color-app-surface-strong)] text-[var(--color-app-text)] shadow-sm hover:border-[var(--color-brand-300)] hover:bg-[var(--color-app-highlight)] focus-visible:ring-[var(--color-brand-300)]',
		ghost:
			'text-[var(--color-app-muted)] hover:bg-[var(--color-app-accent-soft)] hover:text-[var(--color-app-text)] focus-visible:ring-[var(--color-brand-300)]',
		destructive:
			'bg-[var(--color-danger-600)] text-[var(--color-app-inverse)] shadow-sm hover:bg-[var(--color-danger-700)] focus-visible:ring-[var(--color-danger-300)]'
	};

	const sizes: Record<Size, string> = {
		default: 'h-10 px-4 py-2',
		sm: 'h-9 px-4',
		lg: 'h-11 px-5 text-base'
	};

	const normalizedHref = $derived(href && href.startsWith('/') ? resolve(href) : href);
</script>

{#if href}
	<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
	<a
		class={cn(base, variants[variant], sizes[size], className)}
		href={normalizedHref}
		aria-disabled={disabled}
		{onclick}
	>
		{@render children?.()}
	</a>
{:else}
	<button
		class={cn(base, variants[variant], sizes[size], className)}
		{type}
		{disabled}
		formaction={formAction}
		{onclick}
	>
		{@render children?.()}
	</button>
{/if}
