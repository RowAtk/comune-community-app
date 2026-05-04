<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { communitySections } from '$lib/features/community-nav';

	let { data } = $props();
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-3">
			<div
				class="theme-hero-chip inline-flex rounded-full px-4 py-2 text-xs font-semibold tracking-[0.24em] uppercase"
			>
				Community Dashboard
			</div>
			<h1 class="text-4xl font-semibold tracking-tight text-[var(--color-app-text)]">
				{data.community?.name ?? 'Community'}
			</h1>
			<p class="max-w-2xl text-base leading-7 text-[var(--color-app-muted)]">
				Use the community workspace to move from occupancy structure into resident operations and
				maintenance invoicing without losing the tenant boundary.
			</p>
		</div>
		<div class="flex gap-3">
			<Button
				href={`/organizations/${data.community?.organization_id}/communities/${data.community?.id}/units`}
			>
				Open Units
			</Button>
			<Button
				href={`/organizations/${data.community?.organization_id}/communities/${data.community?.id}/residents`}
				variant="outline"
			>
				View Residents
			</Button>
		</div>
	</div>

	<div class="grid gap-5 md:grid-cols-3">
		<Card class="theme-highlight-card border-0">
			<CardHeader>
				<CardDescription class="text-[color:rgba(253,250,244,0.78)]">Current status</CardDescription
				>
				<CardTitle class="text-[var(--color-app-inverse)]"
					>{data.community?.status ?? 'Unknown'}</CardTitle
				>
			</CardHeader>
		</Card>
		<Card>
			<CardHeader>
				<CardDescription>Timezone</CardDescription>
				<CardTitle>{data.community?.timezone ?? 'Not loaded'}</CardTitle>
			</CardHeader>
		</Card>
		<Card>
			<CardHeader>
				<CardDescription>Default next move</CardDescription>
				<CardTitle>Open your occupancy foundation</CardTitle>
			</CardHeader>
		</Card>
	</div>

	<div class="grid gap-6 xl:grid-cols-[minmax(0,1.15fr)_minmax(320px,0.85fr)]">
		<Card>
			<CardHeader>
				<CardTitle>Run The Community</CardTitle>
				<CardDescription>
					Start with the physical model, then move into people and maintenance-fee workflows.
				</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div
					class="rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] p-4"
				>
					<div class="flex items-start justify-between gap-3">
						<div class="space-y-1">
							<p class="text-lg font-semibold text-[var(--color-app-text)]">Occupancy Foundation</p>
							<p class="text-sm text-[var(--color-app-muted)]">
								Build the unit, household, and resident structure first so later workflows stay
								grounded.
							</p>
						</div>
						<Badge variant="secondary">Core</Badge>
					</div>
					<div class="mt-4 grid gap-3 md:grid-cols-3">
						{#each communitySections.filter( (section) => ['units', 'households', 'residents'].includes(section.href) ) as section (section.href)}
							<div class="theme-soft-card rounded-2xl p-4">
								<p class="font-medium text-[var(--color-app-text)]">{section.label}</p>
								<p class="mt-1 text-sm text-[var(--color-app-muted)]">{section.description}</p>
								<Button
									href={`/organizations/${data.community?.organization_id}/communities/${data.community?.id}/${section.href}`}
									variant="outline"
									size="sm"
									class="mt-4 w-full justify-between"
								>
									Open
									<span aria-hidden="true">→</span>
								</Button>
							</div>
						{/each}
					</div>
				</div>

				<div
					class="rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] p-4"
				>
					<div class="flex items-start justify-between gap-3">
						<div class="space-y-1">
							<p class="text-lg font-semibold text-[var(--color-app-text)]">Financial Operations</p>
							<p class="text-sm text-[var(--color-app-muted)]">
								Track maintenance dues, overdue households, and resident invoice visibility.
							</p>
						</div>
						<Badge>Invoicing</Badge>
					</div>
					<div
						class="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-[var(--color-app-border)] pt-4"
					>
						<p class="text-sm text-[var(--color-app-muted)]">
							Use invoicing when the physical and resident structure is ready.
						</p>
						<Button
							href={`/organizations/${data.community?.organization_id}/communities/${data.community?.id}/invoicing`}
							variant="outline"
						>
							Open Invoicing
						</Button>
					</div>
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<CardTitle>Workspace Signals</CardTitle>
				<CardDescription>
					Keep the community context obvious while you move through modules.
				</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="theme-soft-card rounded-2xl p-4">
					<p class="text-sm tracking-[0.18em] text-[var(--color-app-subtle)] uppercase">
						Primary workspace
					</p>
					<p class="mt-2 text-lg font-semibold text-[var(--color-app-text)]">
						Occupancy before operations
					</p>
					<p class="mt-2 text-sm text-[var(--color-app-muted)]">
						Keep the physical model and resident structure clean first, then layer invoicing on top
						of it.
					</p>
				</div>
				<div class="theme-soft-card rounded-2xl p-4">
					<p class="text-sm tracking-[0.18em] text-[var(--color-app-subtle)] uppercase">Slug</p>
					<p class="mt-2 text-lg font-semibold text-[var(--color-app-text)]">
						{data.community?.slug ?? 'Not loaded'}
					</p>
				</div>
				<div class="theme-soft-card rounded-2xl p-4">
					<p class="text-sm tracking-[0.18em] text-[var(--color-app-subtle)] uppercase">
						Suggested flow
					</p>
					<ol class="mt-3 space-y-2 text-sm text-[var(--color-app-muted)]">
						<li>1. Set up units.</li>
						<li>2. Group households.</li>
						<li>3. Add residents.</li>
						<li>4. Turn on invoicing.</li>
					</ol>
				</div>
			</CardContent>
		</Card>
	</div>
</div>
