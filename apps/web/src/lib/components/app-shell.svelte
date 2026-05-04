<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		eyebrow,
		title,
		description = '',
		class: className = '',
		actions,
		children
	}: {
		eyebrow: string;
		title: string;
		description?: string;
		class?: string;
		actions?: import('svelte').Snippet;
		children?: import('svelte').Snippet;
	} = $props();
</script>

<div class={cn('min-h-screen', className)}>
	<div
		class="mx-auto flex min-h-[calc(100vh-73px)] max-w-6xl flex-col gap-8 px-6 py-12 lg:px-10 lg:py-14"
	>
		<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
			<div class="space-y-3">
				<p class="theme-kicker text-sm font-medium tracking-[0.22em] uppercase">
					{eyebrow}
				</p>
				<h1 class="text-4xl font-semibold tracking-tight text-[var(--color-app-text)]">{title}</h1>
				{#if description}
					<p class="max-w-2xl text-base leading-7 text-[var(--color-app-muted)]">{description}</p>
				{/if}
			</div>
			{#if actions}
				<div class="shrink-0">
					{@render actions()}
				</div>
			{/if}
		</div>

		{@render children?.()}
	</div>
</div>
