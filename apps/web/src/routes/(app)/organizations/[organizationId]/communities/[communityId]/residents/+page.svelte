<script lang="ts">
	import { page } from '$app/state';
	import type { Household, Resident } from '$lib/api/types';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { CardDescription, CardTitle } from '$lib/components/ui/card';
	import ResourceCollection from '$lib/features/community-resources/resource-collection/resource-collection.svelte';
	import ResourceListItemCard from '$lib/features/community-resources/resource-collection/resource-list-item-card.svelte';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	let { data, form } = $props();

	const communityBasePath = $derived(
		`/organizations/${page.params.organizationId}/communities/${page.params.communityId}`
	);

	const householdsByUnit = $derived.by(() => {
		const groups = new Map<string, Household[]>();

		for (const household of data.households) {
			const existing = groups.get(household.unit_id);

			if (existing) {
				existing.push(household);
				continue;
			}

			groups.set(household.unit_id, [household]);
		}

		return groups;
	});

	let selectedUnitId = $state('');
	let selectedHouseholdId = $state('');

	const selectedUnitHouseholds = $derived(
		selectedUnitId ? householdsByUnit.get(selectedUnitId) ?? [] : []
	);

	$effect(() => {
		const incomingUnitId = form?.values?.unit_id ?? '';
		if (incomingUnitId && incomingUnitId !== selectedUnitId) {
			selectedUnitId = incomingUnitId;
		}
	});

	$effect(() => {
		const incomingHouseholdId = form?.values?.household_id ?? '';
		if (incomingHouseholdId && incomingHouseholdId !== selectedHouseholdId) {
			selectedHouseholdId = incomingHouseholdId;
		}
	});

	$effect(() => {
		if (!selectedUnitId) {
			if (selectedHouseholdId) {
				selectedHouseholdId = '';
			}
			return;
		}

		const currentUnitHouseholds = householdsByUnit.get(selectedUnitId) ?? [];

		if (currentUnitHouseholds.length === 1) {
			const [defaultHousehold] = currentUnitHouseholds;
			if (selectedHouseholdId !== defaultHousehold.id) {
				selectedHouseholdId = defaultHousehold.id;
			}
			return;
		}

		if (!currentUnitHouseholds.some((household) => household.id === selectedHouseholdId)) {
			selectedHouseholdId = '';
		}
	});

	function statusVariant(status: string) {
		return status === 'ACTIVE' ? 'success' : 'outline';
	}
</script>

<ResourceCollection
	resourceName="residents"
	resourceHeading="People and contacts"
	resourceDescription="Create and manage resident records now, then link a user profile later through invitations when needed."
	creationDescription="Add a resident with their role and unit assignment, then attach a profile later through the invite flow."
	resourceList={data.residents}
	createForm={form}
	apiError={data.apiError}
	resourceEmptyCreateDescription="Create the first resident to start linking households and contacts."
	groupBy={(resident) =>
		data.households.find((household) => household.id === resident.household_id)?.name ||
		(resident.household_id ? 'Unnamed household' : 'Unassigned')
	}
>
	{#snippet headerActions()}
		<Button href={`${communityBasePath}/households`} variant="outline">View Households</Button>
	{/snippet}

	{#snippet createFormBody()}
		<div class="space-y-2">
			<Label for="first_name">First Name</Label>
			<Input id="first_name" name="first_name" required value={form?.values?.first_name ?? ''} />
		</div>
		<div class="space-y-2">
			<Label for="last_name">Last Name</Label>
			<Input id="last_name" name="last_name" required value={form?.values?.last_name ?? ''} />
		</div>
		<div class="space-y-2">
			<Label for="unit_id">Unit</Label>
			<select
				id="unit_id"
				name="unit_id"
				required
				bind:value={selectedUnitId}
				class="flex h-11 w-full rounded-xl border border-(--color-app-border) bg-white px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm outline-none transition focus:border-[var(--color-secondary-400)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-secondary-300)_35%,transparent)]"
			>
				<option value="" disabled>Select a unit</option>
				{#each data.units as unit}
					<option value={unit.id}>{unit.unit_number}</option>
				{/each}
			</select>
		</div>
		<div class="space-y-2">
			<Label for="household_id">Household</Label>
			<select
				id="household_id"
				name="household_id"
				bind:value={selectedHouseholdId}
				disabled={!selectedUnitId || selectedUnitHouseholds.length === 0}
				class="flex h-11 w-full rounded-xl border border-(--color-app-border) bg-white px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm outline-none transition focus:border-[var(--color-secondary-400)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-secondary-300)_35%,transparent)] disabled:bg-slate-50 disabled:text-slate-400"
			>
				<option value="">
					{#if !selectedUnitId}
						Select a unit first
					{:else if selectedUnitHouseholds.length === 0}
						No households available
					{:else}
						Select a household
					{/if}
				</option>
				{#each selectedUnitHouseholds as household}
					<option value={household.id}>{household.name || 'Unnamed household'}</option>
				{/each}
			</select>
			{#if selectedUnitHouseholds.length === 1}
				<p class="text-sm text-[var(--color-app-muted)]">
					The only household for this unit has been selected automatically.
				</p>
			{/if}
		</div>
		<div class="space-y-2">
			<Label for="resident_type">Resident Type</Label>
			<select
				id="resident_type"
				name="resident_type"
				required
				class="flex h-11 w-full rounded-xl border border-(--color-app-border) bg-white px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm outline-none transition focus:border-[var(--color-secondary-400)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-secondary-300)_35%,transparent)]"
			>
				<option value="" disabled selected={!form?.values?.resident_type}>Select a resident type</option>
				<option value="OWNER" selected={form?.values?.resident_type === 'OWNER'}>Owner</option>
				<option value="TENANT" selected={form?.values?.resident_type === 'TENANT'}>Tenant</option>
				<option value="DEPENDENT" selected={form?.values?.resident_type === 'DEPENDENT'}>
					Dependent
				</option>
				<option value="OCCUPANT" selected={form?.values?.resident_type === 'OCCUPANT'}>Occupant</option>
			</select>
		</div>
		<div class="space-y-2">
			<Label for="household_role">Household Access</Label>
			<select
				id="household_role"
				name="household_role"
				class="flex h-11 w-full rounded-xl border border-(--color-app-border) bg-white px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm outline-none transition focus:border-[var(--color-secondary-400)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-secondary-300)_35%,transparent)]"
			>
				<option value="" selected={!form?.values?.household_role}>Default access</option>
				<option value="HOUSEHOLD_ADMIN" selected={form?.values?.household_role === 'HOUSEHOLD_ADMIN'}>
					Household Admin
				</option>
				<option value="HOUSEHOLD_MEMBER" selected={form?.values?.household_role === 'HOUSEHOLD_MEMBER'}>
					Household Member
				</option>
				<option value="HOUSEHOLD_VIEWER" selected={form?.values?.household_role === 'HOUSEHOLD_VIEWER'}>
					Household Viewer
				</option>
			</select>
			<p class="text-sm text-[var(--color-app-muted)]">
				Leave this blank to keep the resident at the lowest household privilege level.
			</p>
		</div>
		<div class="space-y-2">
			<Label for="email">Email</Label>
			<Input id="email" name="email" type="email" value={form?.values?.email ?? ''} />
		</div>
		<div class="space-y-2">
			<Label for="phone">Phone</Label>
			<Input id="phone" name="phone" value={form?.values?.phone ?? ''} />
		</div>
		<div class="space-y-2 md:col-span-2">
			<Label for="move_in_date">Move In Date</Label>
			<Input id="move_in_date" name="move_in_date" type="date" value={form?.values?.move_in_date ?? ''} />
		</div>
	{/snippet}

	{#snippet listItemCard(resident: Resident)}
		<ResourceListItemCard>
			<div class="flex items-start justify-between gap-4">
				<div class="space-y-1">
					<CardTitle>{resident.first_name} {resident.last_name}</CardTitle>
					<CardDescription>{resident.resident_type}</CardDescription>
				</div>
				<Badge variant={statusVariant(resident.status)}>{resident.status}</Badge>
			</div>
			{#if resident.household_role}
				<div class="flex flex-wrap gap-2">
					<Badge variant="outline">{resident.household_role}</Badge>
				</div>
			{/if}
			<div class="space-y-1 text-sm text-slate-600">
				<p>{resident.linked_user_email || resident.email || 'No email on file'}</p>
				<p>{resident.linked_user_phone || resident.phone || 'No phone on file'}</p>
			</div>
			<Button
				href={`${communityBasePath}/residents/${resident.id}`}
				variant="outline"
				size="sm"
				class="self-start"
			>
				View Details
			</Button>
		</ResourceListItemCard>
	{/snippet}
</ResourceCollection>
