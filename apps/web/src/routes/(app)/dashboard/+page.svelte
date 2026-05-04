<script lang="ts">
	import AppShell from '$lib/components/app-shell.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Building2, MapPinned, Plus } from 'lucide-svelte';

	let { data } = $props();

	const formatRole = (role: string) =>
		role
			.toLowerCase()
			.split('_')
			.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
			.join(' ');

	const organizationCount = data.dashboard.organizations.length;
	const communityCount = data.dashboard.communities.length;
	const hasCommunityWork = communityCount > 0;
	const recommendedAction = hasCommunityWork
		? 'Open a community workspace'
		: 'Create your first organization';
</script>

<svelte:head>
	<title>Dashboard | Comune</title>
</svelte:head>

<AppShell
	eyebrow="User Dashboard"
	title={`Welcome back, ${data.auth.user.first_name || data.auth.user.email}`}
	description="Choose your next operational workspace quickly: communities first, then organizations when you need setup or governance changes."
>
	{#snippet actions()}
		<form method="POST" action="/">
			<Button type="submit" variant="ghost" formAction="/?/logout">Log Out</Button>
		</form>
	{/snippet}

	<div class="grid gap-5 md:grid-cols-3">
		<Card class="theme-highlight-card border-0">
			<CardHeader>
				<CardDescription class="text-[color:rgba(253,250,244,0.78)]"
					>Community workspaces</CardDescription
				>
				<CardTitle class="text-[var(--color-app-inverse)]">{communityCount}</CardTitle>
			</CardHeader>
		</Card>
		<Card>
			<CardHeader>
				<CardDescription>Organizations</CardDescription>
				<CardTitle>{organizationCount}</CardTitle>
			</CardHeader>
		</Card>
		<Card>
			<CardHeader>
				<CardDescription>Fastest next action</CardDescription>
				<CardTitle>{recommendedAction}</CardTitle>
			</CardHeader>
		</Card>
	</div>

	<div class="grid gap-6 xl:grid-cols-[minmax(0,1.25fr)_minmax(320px,0.75fr)]">
		<Card>
			<CardHeader>
				<div class="flex items-start justify-between gap-4">
					<div class="flex items-center gap-3">
						<div
							class="rounded-2xl bg-[var(--color-secondary-600)] p-3 text-[var(--color-app-inverse)]"
						>
							<MapPinned class="h-5 w-5" />
						</div>
						<div>
							<CardTitle>Community Workspaces</CardTitle>
							<CardDescription
								>Jump into the communities where day-to-day work actually happens.</CardDescription
							>
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
					{#each data.dashboard.communities as community (community.community_id)}
						<div
							class="rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] p-4"
						>
							<div class="flex flex-wrap items-start justify-between gap-3">
								<div class="space-y-1">
									<p class="text-lg font-semibold text-[var(--color-app-text)]">
										{community.community_name}
									</p>
									<p class="text-sm text-[var(--color-app-muted)]">{community.organization_name}</p>
									<p class="text-xs tracking-[0.18em] text-[var(--color-app-subtle)] uppercase">
										{community.community_slug}
									</p>
								</div>
								<Badge variant="default" class="capitalize">{formatRole(community.role)}</Badge>
							</div>
							<div
								class="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-[var(--color-app-border)] pt-4"
							>
								<p class="text-sm text-[var(--color-app-muted)]">
									Open this workspace to manage residents, households, and invoicing.
								</p>
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
					<div
						class="rounded-2xl border border-dashed border-[var(--color-app-border-strong)] bg-[var(--color-app-panel)] p-6 text-sm text-[var(--color-app-muted)]"
					>
						No organization memberships yet.
					</div>
				{/if}
			</CardContent>
		</Card>

		<div class="space-y-6">
			<Card class="theme-soft-card border-0">
				<CardHeader>
					<CardTitle>Where To Start Today</CardTitle>
					<CardDescription>
						Keep the next move obvious instead of treating every workspace as equally urgent.
					</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4">
					<div
						class="rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] p-4"
					>
						<p class="text-sm tracking-[0.18em] text-[var(--color-app-subtle)] uppercase">
							Primary move
						</p>
						<p class="mt-2 text-lg font-semibold text-[var(--color-app-text)]">
							{recommendedAction}
						</p>
						<p class="mt-2 text-sm text-[var(--color-app-muted)]">
							{#if hasCommunityWork}
								Community workspaces are where residents, households, and invoicing activity now
								live.
							{:else}
								An organization unlocks community setup, memberships, and the rest of the
								operational structure.
							{/if}
						</p>
					</div>
					<div
						class="rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] p-4"
					>
						<p class="text-sm tracking-[0.18em] text-[var(--color-app-subtle)] uppercase">
							Suggested sequence
						</p>
						<ol class="mt-3 space-y-2 text-sm text-[var(--color-app-muted)]">
							<li>1. Enter a community workspace.</li>
							<li>2. Review units, households, and residents.</li>
							<li>3. Use invoicing for maintenance dues and overdue visibility.</li>
						</ol>
					</div>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<div class="flex items-start justify-between gap-4">
						<div class="flex items-center gap-3">
							<div class="theme-highlight-card rounded-2xl p-3">
								<Building2 class="h-5 w-5" />
							</div>
							<div>
								<CardTitle>Organizations</CardTitle>
								<CardDescription
									>Use these for top-level setup, governance, and future portfolio actions.</CardDescription
								>
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
						{#each data.dashboard.organizations as organization (organization.organization_id)}
							<div class="theme-soft-card rounded-2xl p-4">
								<div class="flex flex-wrap items-start justify-between gap-3">
									<div class="space-y-1">
										<p class="text-lg font-semibold text-[var(--color-app-text)]">
											{organization.name}
										</p>
										<p class="text-sm text-[var(--color-app-muted)]">{organization.slug}</p>
									</div>
									<Badge variant="secondary" class="capitalize"
										>{formatRole(organization.role)}</Badge
									>
								</div>
							</div>
						{/each}
					{:else}
						<div
							class="rounded-2xl border border-dashed border-[var(--color-app-border-strong)] bg-[var(--color-app-panel)] p-6 text-sm text-[var(--color-app-muted)]"
						>
							No organization memberships yet.
						</div>
					{/if}
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Join A Resident Invite</CardTitle>
					<CardDescription>
						If a community admin shared a token with you directly, paste it here to open the
						resident claim flow.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form method="GET" action="/invites/residents" class="flex flex-col gap-3 sm:flex-row">
						<input
							name="token"
							placeholder="Paste resident invite token"
							required
							class="flex-1 rounded-xl border border-(--color-app-border) bg-[var(--color-app-surface-strong)] px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm transition outline-none focus:border-[var(--color-brand-300)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-brand-200)_45%,transparent)]"
						/>
						<Button type="submit" variant="outline">Open Invite</Button>
					</form>
				</CardContent>
			</Card>
		</div>
	</div>
</AppShell>
