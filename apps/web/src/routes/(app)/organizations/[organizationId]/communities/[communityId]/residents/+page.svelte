<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';

	let { data } = $props();
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-2">
			<p class="text-sm font-medium uppercase tracking-[0.22em] text-[var(--color-secondary-700)]">
				Residents
			</p>
			<h1 class="text-4xl font-semibold tracking-tight text-slate-950">People and contacts</h1>
			<p class="max-w-2xl text-base leading-7 text-slate-600">
				The resident slice is now wired into the app shell. Next we can add the richer resident form
				with unit and household selection.
			</p>
		</div>
		<div class="flex gap-3">
			<Badge variant="secondary">{data.residents.length} residents</Badge>
			<Button href="../households" variant="outline">View Households</Button>
		</div>
	</div>

	{#if data.apiError}
		<Card class="border-amber-200 bg-amber-50">
			<CardContent class="p-5 text-sm text-amber-900">{data.apiError}</CardContent>
		</Card>
	{/if}

	<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
		{#each data.residents as resident}
			<Card>
				<CardHeader>
					<div class="flex items-start justify-between gap-4">
						<div>
							<CardTitle>{resident.first_name} {resident.last_name}</CardTitle>
							<CardDescription>{resident.resident_type}</CardDescription>
						</div>
						<Badge variant={resident.status === 'ACTIVE' ? 'success' : 'outline'}>
							{resident.status}
						</Badge>
					</div>
				</CardHeader>
				<CardContent class="space-y-2 text-sm text-slate-600">
					<p>Unit: {resident.unit_id}</p>
					<p>Household: {resident.household_id || 'Unassigned'}</p>
					<p>Email: {resident.email || 'Not set'}</p>
				</CardContent>
			</Card>
		{/each}
	</div>
</div>
