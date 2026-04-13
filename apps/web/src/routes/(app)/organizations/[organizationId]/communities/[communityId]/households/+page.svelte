<script lang="ts">
	import { page } from '$app/state';
	import type { Household } from '$lib/api/types';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { CardDescription, CardTitle } from '$lib/components/ui/card';
	import ResourceCollection from '$lib/features/community-resources/resource-collection/resource-collection.svelte';
	import ResourceListItemCard from '$lib/features/community-resources/resource-collection/resource-list-item-card.svelte';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	let { data, form } = $props();

	const householdsBasePath = $derived(
		`/organizations/${page.params.organizationId}/communities/${page.params.communityId}/households`
	);
</script>

<ResourceCollection
	resourceName="households"
	resourceHeading="Occupancy groups"
	resourceDescription="Group together family members, tenants, or roommates while keeping each household mapped to a unit."
	creationDescription="Create a household and assign it to an existing unit."
	resourceList={data.households}
	createForm={form}
	apiError={data.apiError}
	resourceEmptyCreateDescription="Create the first household to start mapping occupancy groups to units."
>
	{#snippet headerActions()}
		<Button href="./units" variant="outline">Back To Units</Button>
	{/snippet}

	{#snippet createFormBody()}
		<div class="space-y-2">
			<Label for="unit_id">Unit</Label>
			<select
				id="unit_id"
				name="unit_id"
				required
				class="flex h-11 w-full rounded-xl border border-(--color-app-border) bg-white px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm outline-none transition focus:border-[var(--color-secondary-400)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-secondary-300)_35%,transparent)]"
			>
				<option value="" disabled selected={!form?.values?.unit_id}>Select a unit</option>
				{#each data.units as unit}
					<option value={unit.id} selected={form?.values?.unit_id === unit.id}>{unit.unit_number}</option>
				{/each}
			</select>
			{#if data.units.length === 0}
				<p class="text-sm text-[var(--color-app-muted)]">Create a unit before adding households.</p>
			{/if}
		</div>
		<div class="space-y-2">
			<Label for="name">Household Name</Label>
			<Input
				id="name"
				name="name"
				placeholder="Brown Family"
				value={form?.values?.name ?? ''}
			/>
		</div>
	{/snippet}

	{#snippet listItemCard(household: Household)}
		<ResourceListItemCard>
			<div class="space-y-1">
				<CardTitle>{household.name || 'Unnamed household'}</CardTitle>
				<CardDescription>Occupancy group for unit assignment.</CardDescription>
			</div>
			<div class="flex flex-wrap gap-2">
				<Badge variant="outline">Unit {household.unit_id}</Badge>
				<Badge variant="secondary">Household</Badge>
			</div>
			<Button href={`${householdsBasePath}/${household.id}`} variant="outline" size="sm" class="self-start">
				View Details
			</Button>
		</ResourceListItemCard>
	{/snippet}
</ResourceCollection>
