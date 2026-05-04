<script lang="ts">
	import { resolve } from '$app/paths';
	import { Button } from '$lib/components/ui/button';
	import favicon from '$lib/assets/favicon.svg';
	import './layout.css';

	let { data, children } = $props();

	const profileName = $derived(data.auth?.user.first_name || data.auth?.user.email || '');
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div data-theme="estate-horizon" class="min-h-screen bg-[var(--color-app-bg)]">
	<header class="theme-topbar sticky top-0 z-40 border-b">
		<div class="mx-auto flex max-w-6xl items-center justify-between gap-4 px-6 py-4 lg:px-10">
			<a href={resolve('/')} class="flex items-center gap-3">
				<div
					class="theme-highlight-card flex h-10 w-10 items-center justify-center rounded-2xl text-sm font-semibold"
				>
					C
				</div>
				<div>
					<p class="theme-kicker text-xs tracking-[0.24em] uppercase">Comune</p>
					<p class="text-sm font-semibold text-[var(--color-app-text)]">Workspace</p>
				</div>
			</a>

			<nav class="flex items-center gap-3">
				{#if data.auth}
					<Button href="/dashboard" variant="ghost">Dashboard</Button>
					<div
						class="hidden rounded-full border border-[var(--color-app-border)] bg-[var(--color-app-panel)] px-3 py-2 text-sm text-[var(--color-app-muted)] sm:block"
					>
						{profileName}
					</div>
					<form method="POST" action="/">
						<Button type="submit" formAction="/?/logout" variant="outline">Log Out</Button>
					</form>
				{:else}
					<Button href="/" variant="outline">Log In</Button>
				{/if}
			</nav>
		</div>
	</header>

	{@render children()}
</div>
