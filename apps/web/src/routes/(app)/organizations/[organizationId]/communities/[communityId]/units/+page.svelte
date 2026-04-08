<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	let { data, form } = $props();
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

	<div class="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
		<Card class="h-fit">
			<CardHeader>
				<CardTitle>Create Unit</CardTitle>
				<CardDescription>Add a new building, apartment, lot, or physical space.</CardDescription>
			</CardHeader>
			<CardContent>
				<form method="POST" class="space-y-4">
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
							class="flex h-11 w-full rounded-xl border border-[var(--color-app-border)] bg-white px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm outline-none transition focus:border-[var(--color-secondary-400)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-secondary-300)_35%,transparent)]"
						>
							<option value="ACTIVE" selected={(form?.values?.status ?? 'ACTIVE') === 'ACTIVE'}>ACTIVE</option>
							<option value="INACTIVE" selected={form?.values?.status === 'INACTIVE'}>INACTIVE</option>
						</select>
					</div>
					{#if form?.createError}
						<p class="text-sm text-[var(--color-danger-700)]">{form.createError}</p>
					{/if}
					{#if form?.createSuccess}
						<p class="text-sm text-[var(--color-success-700)]">{form.createSuccess}</p>
					{/if}
					<Button type="submit" class="w-full" formAction="?/create">Create Unit</Button>
				</form>
			</CardContent>
		</Card>

		<div class="space-y-4">
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
				<div class="grid gap-4 md:grid-cols-2">
					{#each data.units as unit}
						<Card class="overflow-hidden">
							<CardHeader class="gap-4 border-b border-slate-100 bg-slate-50">
								<div class="flex items-start justify-between gap-4">
									<div>
										<CardTitle>{unit.unit_number}</CardTitle>
										<CardDescription>{unit.block_floor || 'No block/floor set'}</CardDescription>
									</div>
									<Badge variant={unit.status === 'ACTIVE' ? 'success' : 'outline'}>
										{unit.status}
									</Badge>
								</div>
							</CardHeader>
							<CardContent class="space-y-3 pt-6">
								<div>
									<p class="text-xs uppercase tracking-[0.2em] text-slate-400">Type</p>
									<p class="mt-1 text-sm font-medium text-slate-900">{unit.unit_type || 'Unspecified'}</p>
								</div>
								<div>
									<p class="text-xs uppercase tracking-[0.2em] text-slate-400">Created</p>
									<p class="mt-1 text-sm text-slate-600">
										{new Date(unit.created_at).toLocaleDateString()}
									</p>
								</div>
							</CardContent>
						</Card>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
