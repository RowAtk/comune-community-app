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
				Households
			</p>
			<h1 class="text-4xl font-semibold tracking-tight text-slate-950">Occupancy groups</h1>
			<p class="max-w-2xl text-base leading-7 text-slate-600">
				This page is scaffolded and already connected to the API. The next step is the same create and
				edit workflow we established for units.
			</p>
		</div>
		<div class="flex gap-3">
			<Badge>{data.households.length} households</Badge>
			<Button href="../units" variant="outline">Back To Units</Button>
		</div>
	</div>

	{#if data.apiError}
		<Card class="border-amber-200 bg-amber-50">
			<CardContent class="p-5 text-sm text-amber-900">{data.apiError}</CardContent>
		</Card>
	{/if}

	<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
		{#each data.households as household}
			<Card>
				<CardHeader>
					<CardTitle>{household.name || 'Unnamed household'}</CardTitle>
					<CardDescription>Unit {household.unit_id}</CardDescription>
				</CardHeader>
				<CardContent>
					<p class="text-sm text-slate-600">
						Created {new Date(household.created_at).toLocaleDateString()}
					</p>
				</CardContent>
			</Card>
		{/each}
	</div>
</div>
