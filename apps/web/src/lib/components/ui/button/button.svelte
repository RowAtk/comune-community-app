<script lang="ts">
	import { cn } from '$lib/utils';

	type Variant = 'default' | 'secondary' | 'outline' | 'ghost' | 'destructive';
	type Size = 'default' | 'sm' | 'lg';

	let {
		type = 'button',
		variant = 'default',
		size = 'default',
		href = undefined,
		formAction = undefined,
		class: className = '',
		disabled = false,
		children
	}: {
		type?: 'button' | 'submit' | 'reset';
		variant?: Variant;
		size?: Size;
		href?: string;
		formAction?: string;
		class?: string;
		disabled?: boolean;
		children?: import('svelte').Snippet;
	} = $props();

	const base =
		'inline-flex items-center justify-center gap-2 rounded-xl text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50';

	const variants: Record<Variant, string> = {
		default:
			'bg-[var(--color-app-text)] text-white shadow-sm hover:bg-[color-mix(in_oklab,var(--color-app-text)_88%,white)] focus-visible:ring-[var(--color-secondary-500)]',
		secondary:
			'bg-[var(--color-secondary-600)] text-white shadow-sm hover:bg-[var(--color-secondary-700)] focus-visible:ring-[var(--color-secondary-500)]',
		outline:
			'border border-[var(--color-app-border)] bg-white text-[var(--color-app-text)] hover:bg-[var(--color-brand-50)] focus-visible:ring-[var(--color-brand-400)]',
		ghost:
			'text-[var(--color-app-muted)] hover:bg-white hover:text-[var(--color-app-text)] focus-visible:ring-[var(--color-secondary-500)]',
		destructive:
			'bg-[var(--color-danger-600)] text-white shadow-sm hover:bg-[var(--color-danger-700)] focus-visible:ring-[var(--color-danger-500)]'
	};

	const sizes: Record<Size, string> = {
		default: 'h-10 px-4 py-2',
		sm: 'h-9 px-3',
		lg: 'h-11 px-5 text-base'
	};
</script>

{#if href}
	<a class={cn(base, variants[variant], sizes[size], className)} {href} aria-disabled={disabled}>
		{@render children?.()}
	</a>
{:else}
	<button class={cn(base, variants[variant], sizes[size], className)} {type} {disabled} formaction={formAction}>
		{@render children?.()}
	</button>
{/if}
