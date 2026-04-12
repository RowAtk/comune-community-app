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

<div
	class={cn(
		'min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(14,165,233,0.16),_transparent_28%),linear-gradient(180deg,_#f8fafc_0%,_#eef2ff_100%)]',
		className
	)}
>
	<div class="mx-auto flex min-h-[calc(100vh-73px)] max-w-6xl flex-col gap-8 px-6 py-12 lg:px-10 lg:py-14">
		<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
			<div class="space-y-3">
				<p class="text-sm font-medium uppercase tracking-[0.22em] text-[var(--color-secondary-700)]">
					{eyebrow}
				</p>
				<h1 class="text-4xl font-semibold tracking-tight text-slate-950">{title}</h1>
				{#if description}
					<p class="max-w-2xl text-base leading-7 text-slate-600">{description}</p>
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
