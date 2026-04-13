<script lang="ts">
	import { page } from '$app/state';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { CardDescription, CardTitle } from '$lib/components/ui/card';
	import ResourceCollection from '$lib/features/community-resources/resource-collection/resource-collection.svelte';
	import ResourceListItemCard from '$lib/features/community-resources/resource-collection/resource-list-item-card.svelte';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { Unit } from '$lib/api/types';

	let { data, form } = $props();

	const unitsBasePath = $derived(
		`/organizations/${page.params.organizationId}/communities/${page.params.communityId}/units`
	);

	function statusVariant(status: string) {
		return status === 'ACTIVE' ? 'success' : 'outline';
	}

	function getUnitGroup(unit: Unit) {
		const [block] = unit.block_floor?.split('/') ?? [];
		return block?.trim() || 'Unassigned';
	}
</script>

<ResourceCollection
	resourceName="units"
	resourceHeading="Physical spaces"
	resourceDescription="Use units as the foundational map of the community. Households and residents can attach to these records next."
	creationDescription="Add a new building, apartment, lot, or physical space without leaving the listing view."
	resourceList={data.units}
	createForm={form}
	groupBy={getUnitGroup}
	resourceEmptyCreateDescription="Create the first unit to establish the physical map of this community."
>
	{#snippet createFormBody()}
		<div class="space-y-2">
			<Label for="unit_number">Unit Number</Label>
			<Input
				id="unit_number"
				name="unit_number"
				placeholder="A-12"
				required
				value={form?.values?.unit_number ?? ''}
			/>
		</div>
		<div class="space-y-2">
			<Label for="block_floor">Block / Floor</Label>
			<Input
				id="block_floor"
				name="block_floor"
				placeholder="Block B / Level 2"
				value={form?.values?.block_floor ?? ''}
			/>
		</div>
		<div class="space-y-2">
			<Label for="unit_type">Unit Type</Label>
			<Input
				id="unit_type"
				name="unit_type"
				placeholder="Apartment, Villa, Lot"
				value={form?.values?.unit_type ?? ''}
			/>
		</div>
		<div class="space-y-2">
			<Label for="status">Status</Label>
			<select
				id="status"
				name="status"
				class="flex h-11 w-full rounded-xl border border-(--color-app-border) bg-white px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm transition outline-none focus:border-[var(--color-secondary-400)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-secondary-300)_35%,transparent)]"
			>
				<option value="ACTIVE" selected={(form?.values?.status ?? 'ACTIVE') === 'ACTIVE'}
					>ACTIVE</option
				>
				<option value="INACTIVE" selected={form?.values?.status === 'INACTIVE'}>INACTIVE</option>
			</select>
		</div>
	{/snippet}

	{#snippet listItemCard(unit: Unit)}
		<ResourceListItemCard>
			<div class="flex items-start justify-between gap-4">
				<div class="space-y-1">
					<CardTitle>{unit.unit_number}</CardTitle>
					<CardDescription>{unit.block_floor || 'No block or floor set'}</CardDescription>
				</div>
				<div class="flex flex-col items-end gap-2">
					<Badge variant={statusVariant(unit.status)}>{unit.status}</Badge>
					<Badge variant="outline">{unit.unit_type || 'Unspecified type'}</Badge>
				</div>
			</div>
			<Button href={`${unitsBasePath}/${unit.id}`} variant="outline" size="sm" class="self-start">
				View Details
			</Button>
		</ResourceListItemCard>
	{/snippet}
</ResourceCollection>
