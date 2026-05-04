<script lang="ts">
	import { page } from '$app/state';
	import type { Resident } from '$lib/api/types';
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

	type HouseholdResidentGroup = {
		id: string;
		label: string;
		href: string | null;
		residents: Resident[];
	};

	const residentsByHousehold = $derived.by(() => {
		const householdMap = new Map<string, HouseholdResidentGroup>(
			data.households.map((household) => [
				household.id,
				{
					id: household.id,
					label: household.name || 'Unnamed household',
					href: `${communityBasePath}/households/${household.id}`,
					residents: []
				}
			])
		);

		const unassignedResidents: Resident[] = [];

		for (const resident of data.residents) {
			if (resident.household_id && householdMap.has(resident.household_id)) {
				householdMap.get(resident.household_id)?.residents.push(resident);
				continue;
			}

			unassignedResidents.push(resident);
		}

		const grouped = Array.from(householdMap.values()).sort((leftGroup, rightGroup) =>
			leftGroup.label.localeCompare(rightGroup.label, undefined, { sensitivity: 'base' })
		);

		if (unassignedResidents.length > 0) {
			grouped.push({
				id: 'unassigned',
				label: 'Unassigned',
				href: null,
				residents: unassignedResidents
			});
		}

		return grouped;
	});

	function statusVariant(status: string) {
		return status === 'ACTIVE' ? 'success' : 'outline';
	}
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-3">
			<div
				class="theme-hero-chip inline-flex rounded-full px-4 py-2 text-xs font-semibold tracking-[0.24em] uppercase"
			>
				Unit
			</div>
			<h1 class="text-4xl font-semibold tracking-tight text-[var(--color-app-text)]">
				{data.unit?.unit_number ?? 'Unit details'}
			</h1>
			<p class="max-w-2xl text-base leading-7 text-[var(--color-app-muted)]">
				Use this view to understand the physical record first, then see the households and residents
				anchored to it.
			</p>
		</div>
		<div class="flex gap-3">
			<Button href={`${communityBasePath}/units`} variant="outline">Back To Units</Button>
			{#if data.unit}
				<Badge variant={statusVariant(data.unit.status)}>{data.unit.status}</Badge>
			{/if}
		</div>
	</div>

	{#if data.apiError || !data.unit}
		<Card class="border-amber-200 bg-amber-50">
			<CardContent class="p-5 text-sm text-amber-900">
				{data.apiError ?? 'Unit details are not available.'}
			</CardContent>
		</Card>
	{:else}
		<div class="grid gap-5 md:grid-cols-3">
			<Card class="theme-highlight-card border-0">
				<CardHeader>
					<CardDescription class="text-[color:rgba(253,250,244,0.78)]">Unit number</CardDescription>
					<CardTitle class="text-[var(--color-app-inverse)]">{data.unit.unit_number}</CardTitle>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader>
					<CardDescription>Block / Floor</CardDescription>
					<CardTitle>{data.unit.block_floor || 'Not set'}</CardTitle>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader>
					<CardDescription>Type</CardDescription>
					<CardTitle>{data.unit.unit_type || 'Unspecified'}</CardTitle>
				</CardHeader>
			</Card>
		</div>

		<div class="grid gap-5 xl:grid-cols-2">
			<Card>
				<CardHeader>
					<CardTitle>Record Metadata</CardTitle>
					<CardDescription>Useful timing information for this unit record.</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4 text-sm text-[var(--color-app-muted)]">
					<div>
						<p
							class="text-xs font-medium tracking-[0.2em] text-[var(--color-app-subtle)] uppercase"
						>
							Created
						</p>
						<p class="mt-1 text-[var(--color-app-text)]">
							{new Date(data.unit.created_at).toLocaleString()}
						</p>
					</div>
					<div>
						<p
							class="text-xs font-medium tracking-[0.2em] text-[var(--color-app-subtle)] uppercase"
						>
							Updated
						</p>
						<p class="mt-1 text-[var(--color-app-text)]">
							{new Date(data.unit.updated_at).toLocaleString()}
						</p>
					</div>
				</CardContent>
			</Card>

			<Card>
				<CardHeader>
					<CardTitle>Next Step</CardTitle>
					<CardDescription
						>Use this unit as the anchor for household and resident records.</CardDescription
					>
				</CardHeader>
				<CardContent class="space-y-3">
					<Button
						href={`${communityBasePath}/households`}
						variant="outline"
						class="w-full justify-between"
					>
						View Households
						<span aria-hidden="true">→</span>
					</Button>
					<Button
						href={`${communityBasePath}/residents`}
						variant="outline"
						class="w-full justify-between"
					>
						View Residents
						<span aria-hidden="true">→</span>
					</Button>
				</CardContent>
			</Card>
		</div>

		<Card>
			<CardHeader>
				<CardTitle>Residents In This Unit</CardTitle>
				<CardDescription>
					Residents are grouped by household so the occupancy structure is easy to scan.
				</CardDescription>
			</CardHeader>
			<CardContent class="space-y-6">
				{#if (data.residents ?? []).length === 0}
					<p class="text-sm text-[var(--color-app-muted)]">
						No residents are assigned to this unit yet.
					</p>
				{:else}
					<div class="space-y-6">
						{#each residentsByHousehold as householdGroup}
							<section class="space-y-3">
								<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
									<div class="space-y-1">
										<h2 class="text-lg font-semibold tracking-tight text-[var(--color-app-text)]">
											{householdGroup.label}
										</h2>
										<p class="text-sm text-[var(--color-app-muted)]">
											{householdGroup.residents.length}
											{householdGroup.residents.length === 1 ? ' resident' : ' residents'}
										</p>
									</div>
									{#if householdGroup.href}
										<Button href={householdGroup.href} variant="outline" size="sm">
											Open Household
										</Button>
									{/if}
								</div>

								<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
									{#each householdGroup.residents as resident}
										<Card>
											<CardHeader class="gap-3">
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
												<div class="space-y-1 text-sm text-[var(--color-app-muted)]">
													<p>{resident.email || 'No email on file'}</p>
													<p>{resident.phone || 'No phone on file'}</p>
												</div>
											</CardHeader>
										</Card>
									{/each}
								</div>
							</section>
						{/each}
					</div>
				{/if}
			</CardContent>
		</Card>
	{/if}
</div>
