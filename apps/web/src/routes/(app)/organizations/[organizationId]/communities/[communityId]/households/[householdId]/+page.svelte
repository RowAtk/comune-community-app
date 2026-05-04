<script lang="ts">
	import { page } from '$app/state';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';

	let { data } = $props();

	const communityBasePath = $derived(
		`/organizations/${page.params.organizationId}/communities/${page.params.communityId}`
	);
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-3">
			<div
				class="theme-hero-chip inline-flex rounded-full px-4 py-2 text-xs font-semibold tracking-[0.24em] uppercase"
			>
				Household
			</div>
			<h1 class="text-4xl font-semibold tracking-tight text-[var(--color-app-text)]">
				{data.household?.name || 'Unnamed household'}
			</h1>
			<p class="max-w-2xl text-base leading-7 text-[var(--color-app-muted)]">
				Use this page to confirm the household anchor, then move into residents or unit-level
				follow-up.
			</p>
		</div>
		<div class="flex gap-3">
			<Button href={`${communityBasePath}/households`} variant="outline">Back To Households</Button>
			<Badge variant="secondary">Household</Badge>
		</div>
	</div>

	{#if data.apiError || !data.household}
		<Card class="border-amber-200 bg-amber-50">
			<CardContent class="p-5 text-sm text-amber-900">
				{data.apiError ?? 'Household details are not available.'}
			</CardContent>
		</Card>
	{:else}
		<div class="grid gap-5 md:grid-cols-3">
			<Card class="theme-highlight-card border-0">
				<CardHeader>
					<CardDescription class="text-[color:rgba(253,250,244,0.78)]"
						>Household name</CardDescription
					>
					<CardTitle class="text-[var(--color-app-inverse)]"
						>{data.household.name || 'Unnamed household'}</CardTitle
					>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader>
					<CardDescription>Assigned unit</CardDescription>
					<CardTitle>{data.unit?.unit_number || 'Unknown unit'}</CardTitle>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader>
					<CardDescription>Unit details</CardDescription>
					<CardTitle
						>{data.unit?.block_floor || data.unit?.unit_type || 'No extra unit details'}</CardTitle
					>
				</CardHeader>
			</Card>
		</div>

		<div class="grid gap-5 xl:grid-cols-2">
			<Card>
				<CardHeader>
					<CardTitle>Record Metadata</CardTitle>
					<CardDescription>Helpful timing information for this household record.</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4 text-sm text-[var(--color-app-muted)]">
					<div>
						<p
							class="text-xs font-medium tracking-[0.2em] text-[var(--color-app-subtle)] uppercase"
						>
							Created
						</p>
						<p class="mt-1 text-[var(--color-app-text)]">
							{new Date(data.household.created_at).toLocaleString()}
						</p>
					</div>
					<div>
						<p
							class="text-xs font-medium tracking-[0.2em] text-[var(--color-app-subtle)] uppercase"
						>
							Updated
						</p>
						<p class="mt-1 text-[var(--color-app-text)]">
							{new Date(data.household.updated_at).toLocaleString()}
						</p>
					</div>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Linked Unit</CardTitle>
					<CardDescription>Jump directly to the unit this household is assigned to.</CardDescription
					>
				</CardHeader>
				<CardContent>
					{#if data.unit}
						<Button
							href={`${communityBasePath}/units/${data.unit.id}`}
							variant="outline"
							class="w-full justify-between"
						>
							Open {data.unit.unit_number}
							<span aria-hidden="true">→</span>
						</Button>
					{:else}
						<p class="text-sm text-[var(--color-app-muted)]">
							The linked unit could not be loaded for this household.
						</p>
					{/if}
				</CardContent>
			</Card>
		</div>
	{/if}
</div>
