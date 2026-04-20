<script lang="ts">
	import { page } from '$app/state';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { communitySections } from '$lib/features/community-nav';

	let { data, children } = $props();

	const currentPath = $derived(page.url.pathname.replace(/\/$/, ''));
</script>

<div class="min-h-screen bg-[linear-gradient(180deg,#f7f8fc_0%,#edf2f7_100%)]">
	<div class="mx-auto grid min-h-screen max-w-7xl grid-cols-1 lg:grid-cols-[280px_minmax(0,1fr)]">
		<aside
			class="border-b border-slate-200 bg-[linear-gradient(180deg,#0f172a_0%,#111827_100%)] px-6 py-8 text-slate-100 lg:border-r lg:border-b-0 lg:border-slate-800"
		>
			<div class="space-y-6">
				<div class="space-y-2">
					<p class="text-xs tracking-[0.28em] text-slate-400 uppercase">Comune</p>
					<h1 class="text-2xl font-semibold tracking-tight">
						{data.community?.name ?? 'Community workspace'}
					</h1>
					<p class="text-sm leading-6 text-slate-400">Community operations workspace</p>
				</div>

				<div class="flex flex-wrap gap-2">
					<Badge variant="secondary">{data.community?.status ?? 'Unavailable'}</Badge>
					{#if data.community?.timezone}
						<Badge variant="outline" class="border-slate-700 bg-slate-900 text-slate-300">
							{data.community.timezone}
						</Badge>
					{/if}
				</div>

				<nav class="space-y-2">
					{#each communitySections as section}
						{@const href = `/organizations/${data.community?.organization_id ?? ''}/communities/${data.community?.id ?? ''}${section.href ? `/${section.href}` : ''}`}
						<a
							class={`block rounded-2xl px-4 py-3 transition ${
								currentPath === href
									? 'bg-white text-slate-950 shadow-lg'
									: 'text-slate-300 hover:bg-slate-800 hover:text-white'
							}`}
							{href}
						>
							<p class="font-medium">{section.label}</p>
							<p class="mt-1 text-sm opacity-75">{section.description}</p>
						</a>
					{/each}
				</nav>

				<Button
					href="/"
					variant="ghost"
					class="w-full justify-start text-slate-300 hover:bg-slate-800 hover:text-white"
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
