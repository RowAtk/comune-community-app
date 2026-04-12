<script lang="ts">
	import AppShell from '$lib/components/app-shell.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Building2, MapPinned, Plus } from 'lucide-svelte';

	let { data } = $props();

	const formatRole = (role: string) =>
		role
			.toLowerCase()
			.split('_')
			.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
			.join(' ');
</script>

<svelte:head>
	<title>Dashboard | Comune</title>
</svelte:head>

<AppShell
	eyebrow="User Dashboard"
	title={`Welcome back, ${data.auth.user.first_name || data.auth.user.email}`}
	description="Start by choosing one of your organizations or jump directly into a community workspace."
>
	{#snippet actions()}
		<form method="POST" action="/">
			<Button type="submit" variant="ghost" formAction="/?/logout">Log Out</Button>
		</form>
	{/snippet}

	<div class="grid gap-6 xl:grid-cols-2">
		<Card class="border-slate-200/80 bg-white/95">
			<CardHeader>
				<div class="flex items-start justify-between gap-4">
					<div class="flex items-center gap-3">
						<div class="rounded-2xl bg-slate-950 p-3 text-white">
							<Building2 class="h-5 w-5" />
						</div>
						<div>
							<CardTitle>Organizations</CardTitle>
							<CardDescription>Your current organization memberships.</CardDescription>
						</div>
					</div>
					<Button href="/organizations/new" variant="outline" size="sm">
						<Plus class="h-4 w-4" />
						New Organization
					</Button>
				</div>
			</CardHeader>
			<CardContent class="space-y-4">
				{#if data.dashboard.organizations.length > 0}
					{#each data.dashboard.organizations as organization}
						<div class="rounded-2xl border border-slate-200 bg-slate-50 p-4">
							<div class="flex flex-wrap items-start justify-between gap-3">
								<div class="space-y-1">
									<p class="text-lg font-semibold text-slate-950">{organization.name}</p>
									<p class="text-sm text-slate-600">{organization.slug}</p>
								</div>
								<Badge variant="secondary" class="capitalize">{formatRole(organization.role)}</Badge>
							</div>
						</div>
					{/each}
				{:else}
					<div class="rounded-2xl border border-dashed border-slate-300 bg-slate-50 p-6 text-sm text-slate-600">
						No organization memberships yet.
					</div>
				{/if}
			</CardContent>
		</Card>

		<Card class="border-slate-200/80 bg-white/95">
			<CardHeader>
				<div class="flex items-start justify-between gap-4">
					<div class="flex items-center gap-3">
						<div class="rounded-2xl bg-[var(--color-secondary-700)] p-3 text-white">
							<MapPinned class="h-5 w-5" />
						</div>
						<div>
							<CardTitle>Communities</CardTitle>
							<CardDescription>Your current community memberships.</CardDescription>
						</div>
					</div>
					<Button href="/communities/new" variant="outline" size="sm">
						<Plus class="h-4 w-4" />
						New Community
					</Button>
				</div>
			</CardHeader>
			<CardContent class="space-y-4">
				{#if data.dashboard.communities.length > 0}
					{#each data.dashboard.communities as community}
						<div class="rounded-2xl border border-slate-200 bg-slate-50 p-4">
							<div class="flex flex-wrap items-start justify-between gap-3">
								<div class="space-y-1">
									<p class="text-lg font-semibold text-slate-950">{community.community_name}</p>
									<p class="text-sm text-slate-600">{community.organization_name}</p>
								</div>
								<Badge variant="default" class="capitalize">{formatRole(community.role)}</Badge>
							</div>
							<div class="mt-4 flex flex-wrap items-center justify-between gap-3">
								<p class="text-sm text-slate-500">{community.community_slug}</p>
								<Button
									href={`/organizations/${community.organization_id}/communities/${community.community_id}`}
									variant="outline"
									size="sm"
								>
									Open Community
								</Button>
							</div>
						</div>
					{/each}
				{:else}
					<div class="rounded-2xl border border-dashed border-slate-300 bg-slate-50 p-6 text-sm text-slate-600">
						No community memberships yet.
					</div>
				{/if}
			</CardContent>
		</Card>
	</div>
</AppShell>
