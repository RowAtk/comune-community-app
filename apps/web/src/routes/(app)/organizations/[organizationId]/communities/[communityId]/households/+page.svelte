<script lang="ts">
	import { CounterBadge, Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';

	let { data } = $props();
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-2">
			<p class="text-sm font-medium tracking-[0.22em] text-[var(--color-secondary-700)] uppercase">
				Households
			</p>
			<h1 class="text-4xl font-semibold tracking-tight text-slate-950">Occupancy groups</h1>
			<p class="max-w-2xl text-base leading-7 text-slate-600">
				This page is scaffolded and already connected to the API. The next step is reusing the new
				compact card-plus-expand pattern here for household creation and edits.
			</p>
		</div>
		<div class="flex gap-3">
			<CounterBadge counter={data.households?.length || 0} pl="households" />
			<Button href="./units" variant="outline">Back To Units</Button>
		</div>
	</div>

	{#if data.apiError}
		<Card class="border-amber-200 bg-amber-50">
			<CardContent class="p-5 text-sm text-amber-900">{data.apiError}</CardContent>
		</Card>
	{/if}

	<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
		{#each data.households as household}
			<Card class="transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md">
				<CardHeader class="gap-4">
					<div class="space-y-1">
						<CardTitle>{household.name || 'Unnamed household'}</CardTitle>
						<CardDescription>Occupancy group for unit assignment.</CardDescription>
					</div>
					<div class="flex flex-wrap gap-2">
						<Badge variant="outline">Unit {household.unit_id}</Badge>
						<Badge variant="secondary">Household</Badge>
					</div>
				</CardHeader>
			</Card>
		{/each}
	</div>
</div>
