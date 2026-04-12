<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { data, children } = $props();

	const profileName = $derived(data.auth?.user.first_name || data.auth?.user.email || '');
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="min-h-screen bg-[var(--color-app-bg)]">
	<header class="sticky top-0 z-40 border-b border-slate-200/80 bg-white/85 backdrop-blur">
		<div class="mx-auto flex max-w-6xl items-center justify-between gap-4 px-6 py-4 lg:px-10">
			<a href="/" class="flex items-center gap-3">
				<div class="flex h-10 w-10 items-center justify-center rounded-2xl bg-slate-950 text-sm font-semibold text-white">
					C
				</div>
				<div>
					<p class="text-xs uppercase tracking-[0.24em] text-slate-500">Comune</p>
					<p class="text-sm font-semibold text-slate-950">Workspace</p>
				</div>
			</a>

			<nav class="flex items-center gap-3">
				{#if data.auth}
					<Button href="/dashboard" variant="ghost" class="text-slate-700">
						Dashboard
					</Button>
					<div class="hidden rounded-full border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-700 sm:block">
						{profileName}
					</div>
					<form method="POST" action="/">
						<Button type="submit" formAction="/?/logout" variant="outline">
							Log Out
						</Button>
					</form>
				{:else}
					<Button href="/" variant="outline">Log In</Button>
				{/if}
			</nav>
		</div>
	</header>

	{@render children()}
</div>
