	<script lang="ts">
	import { slide } from 'svelte/transition';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { cn } from '$lib/utils';

	let { data, form } = $props();

	let createOpen = $state(false);

	$effect(() => {
		if (form?.createError || form?.createSuccess) {
			createOpen = true;
		}
	});

	function toggleCreateCard(nextState?: boolean) {
		createOpen = nextState ?? !createOpen;
	}

	function statusVariant(status: string) {
		return status === 'ACTIVE' ? 'success' : 'outline';
	}
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-2">
			<p class="text-sm font-medium uppercase tracking-[0.22em] text-[var(--color-secondary-700)]">
				Units
			</p>
			<h1 class="text-4xl font-semibold tracking-tight text-slate-950">Physical spaces</h1>
			<p class="max-w-2xl text-base leading-7 text-slate-600">
				Use units as the foundational map of the community. Households and residents can attach to
				these records next.
			</p>
		</div>
		<Badge variant="secondary">{data.units.length} units</Badge>
	</div>

	<Card
		class={cn(
			'overflow-hidden transition-all duration-300',
			createOpen && 'border-[var(--color-secondary-200)] shadow-lg shadow-[color-mix(in_oklab,var(--color-secondary-200)_35%,transparent)]'
		)}
	>
		<CardHeader class="gap-4">
			<div class="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
				<div class="space-y-2">
					<div class="space-y-1">
						<CardTitle>Create Unit</CardTitle>
						<CardDescription>
							Add a new building, apartment, lot, or physical space without leaving the listing view.
						</CardDescription>
					</div>
				</div>
				<div class="flex shrink-0 gap-2">
					{#if createOpen}
						<Button variant="ghost" size="sm" onclick={() => toggleCreateCard(false)}>Close</Button>
					{:else}
						<Button size="sm" onclick={() => toggleCreateCard(true)}>Add New Unit</Button>
					{/if}
				</div>
			</div>
		</CardHeader>

		{#if createOpen}
			<div transition:slide={{ duration: 220 }}>
				<CardContent class="border-t border-slate-100 pt-6">
					<form method="POST" class="grid gap-4 md:grid-cols-2">
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
								class="flex h-11 w-full rounded-xl border border-(--color-app-border) bg-white px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm outline-none transition focus:border-[var(--color-secondary-400)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-secondary-300)_35%,transparent)]"
							>
								<option value="ACTIVE" selected={(form?.values?.status ?? 'ACTIVE') === 'ACTIVE'}>ACTIVE</option>
								<option value="INACTIVE" selected={form?.values?.status === 'INACTIVE'}>INACTIVE</option>
							</select>
						</div>
						<div class="space-y-3 md:col-span-2">
							{#if form?.createError}
								<p class="text-sm text-[var(--color-danger-700)]">{form.createError}</p>
							{/if}
							{#if form?.createSuccess}
								<p class="text-sm text-[var(--color-success-700)]">{form.createSuccess}</p>
							{/if}
							<div class="flex flex-col gap-3 sm:flex-row sm:justify-end">
								<Button variant="outline" type="button" onclick={() => toggleCreateCard(false)}>
									Cancel
								</Button>
								<Button type="submit" formAction="?/create">Create Unit</Button>
							</div>
						</div>
					</form>
				</CardContent>
			</div>
		{/if}
	</Card>

	{#if data.apiError}
		<Card class="border-amber-200 bg-amber-50">
			<CardContent class="p-5 text-sm text-amber-900">{data.apiError}</CardContent>
		</Card>
	{/if}

	{#if data.units.length === 0}
		<Card>
			<CardHeader>
				<CardTitle>No units yet</CardTitle>
				<CardDescription>
					Create the first unit to establish the physical map of this community.
				</CardDescription>
			</CardHeader>
		</Card>
	{:else}
		<div class={cn('grid gap-4 md:grid-cols-2 transition-all duration-300', !createOpen && 'xl:grid-cols-3')}>
			{#each data.units as unit}
				<Card class="transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md">
					<CardHeader class="gap-3">
						<div class="flex items-start justify-between gap-4">
							<div class="space-y-1">
								<CardTitle>{unit.unit_number}</CardTitle>
								<CardDescription>{unit.block_floor || 'No block or floor set'}</CardDescription>
							</div>
							<div class="flex flex-col items-end gap-2">
								<Badge variant={statusVariant(unit.status)}>
									{unit.status}
								</Badge>
								<Badge variant="outline">{unit.unit_type || 'Unspecified type'}</Badge>
							</div>
						</div>
					</CardHeader>
				</Card>
			{/each}
		</div>
	{/if}
</div>
