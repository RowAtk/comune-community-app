<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { communitySections } from '$lib/features/community-nav';

	let { data, children } = $props();

	const currentPath = $derived(page.url.pathname.replace(/\/$/, ''));
</script>

<div class="min-h-screen">
	<div class="mx-auto grid min-h-screen max-w-7xl grid-cols-1 lg:grid-cols-[280px_minmax(0,1fr)]">
		<aside
			class="theme-sidebar border-b border-[color:rgba(255,255,255,0.08)] px-6 py-8 lg:border-r lg:border-b-0"
		>
			<div class="space-y-6">
				<div class="space-y-2">
					<p class="text-xs tracking-[0.28em] text-[var(--color-app-sidebar-muted)] uppercase">
						Comune
					</p>
					<h1 class="text-2xl font-semibold tracking-tight">
						{data.community?.name ?? 'Community workspace'}
					</h1>
					<p class="text-sm leading-6 text-[var(--color-app-sidebar-muted)]">
						Community operations workspace
					</p>
				</div>

				<div class="flex flex-wrap gap-2">
					<Badge variant="secondary">{data.community?.status ?? 'Unavailable'}</Badge>
					{#if data.community?.timezone}
						<Badge
							variant="outline"
							class="border-white/15 bg-white/6 text-[var(--color-app-sidebar-text)]"
						>
							{data.community.timezone}
						</Badge>
					{/if}
				</div>

				<nav class="space-y-2">
					{#each communitySections as section (section.href)}
						{@const href = `/organizations/${data.community?.organization_id ?? ''}/communities/${data.community?.id ?? ''}${section.href ? `/${section.href}` : ''}`}
						<a
							class={`block rounded-2xl px-4 py-3 transition ${
								currentPath === href
									? 'bg-[var(--color-app-surface-strong)] text-[var(--color-app-text)] shadow-[var(--shadow-soft)]'
									: 'text-[var(--color-app-sidebar-text)] hover:bg-white/8 hover:text-[var(--color-app-inverse)]'
							}`}
							href={resolve(href)}
						>
							<p class="font-medium">{section.label}</p>
							<p class="mt-1 text-sm opacity-75">{section.description}</p>
						</a>
					{/each}
				</nav>

				<Button
					href="/"
					variant="ghost"
					class="w-full justify-start text-[var(--color-app-sidebar-text)] hover:bg-white/8 hover:text-[var(--color-app-inverse)]"
				>
					Switch Community
				</Button>
			</div>
		</aside>

		<main class="px-6 py-8 lg:px-10">
			{#if data.apiError}
				<div
					class="mb-6 rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900"
				>
					{data.apiError}
				</div>
			{/if}
			{@render children()}
		</main>
	</div>
</div>
