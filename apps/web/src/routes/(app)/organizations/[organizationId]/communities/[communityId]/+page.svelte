<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { communitySections } from '$lib/features/community-nav';

	let { data } = $props();
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-3">
			<p class="text-sm font-medium uppercase tracking-[0.22em] text-[var(--color-secondary-700)]">
				Community Dashboard
			</p>
			<h1 class="text-4xl font-semibold tracking-tight text-slate-950">
				{data.community?.name ?? 'Community'}
			</h1>
			<p class="max-w-2xl text-base leading-7 text-slate-600">
				Start with the physical model of the community, then layer on households and residents.
			</p>
		</div>
		<div class="flex gap-3">
			<Button href={`/organizations/${data.community?.organization_id}/communities/${data.community?.id}/units`}>
				Open Units
			</Button>
			<Button
				href={`/organizations/${data.community?.organization_id}/communities/${data.community?.id}/residents`}
				variant="outline"
			>
				View Residents
			</Button>
		</div>
	</div>

	<div class="grid gap-5 md:grid-cols-3">
		<Card class="bg-slate-950 text-white">
			<CardHeader>
				<CardDescription class="text-slate-300">Current status</CardDescription>
				<CardTitle class="text-white">{data.community?.status ?? 'Unknown'}</CardTitle>
			</CardHeader>
		</Card>
		<Card>
			<CardHeader>
				<CardDescription>Timezone</CardDescription>
				<CardTitle>{data.community?.timezone ?? 'Not loaded'}</CardTitle>
			</CardHeader>
		</Card>
		<Card>
			<CardHeader>
				<CardDescription>Community slug</CardDescription>
				<CardTitle>{data.community?.slug ?? 'Not loaded'}</CardTitle>
			</CardHeader>
		</Card>
	</div>

	<div class="grid gap-5 xl:grid-cols-3">
		{#each communitySections.filter((section) => section.href) as section}
			<Card>
				<CardHeader>
					<div class="flex items-center justify-between">
						<CardTitle>{section.label}</CardTitle>
						<Badge>{section.label}</Badge>
					</div>
					<CardDescription>{section.description}</CardDescription>
				</CardHeader>
				<CardContent>
					<Button
						href={`/organizations/${data.community?.organization_id}/communities/${data.community?.id}/${section.href}`}
						variant="outline"
						class="w-full justify-between"
					>
						Open {section.label}
						<span aria-hidden="true">→</span>
					</Button>
				</CardContent>
			</Card>
		{/each}
	</div>
</div>
