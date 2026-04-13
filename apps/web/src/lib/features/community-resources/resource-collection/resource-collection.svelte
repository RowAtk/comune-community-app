<script lang="ts" generics="T extends CommunityResource, TValues extends Record<string, unknown> = Record<string, unknown>">
	import { untrack } from 'svelte';
	import type { CommunityResource } from '$lib/api/types';
	import { CounterBadge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { capitalize, cn, dePluralize } from '$lib/utils';
	import { slide } from 'svelte/transition';
	import type { ResourceCollectionProps } from './types';

	let {
		resourceName = '',
		resourceSingleName = dePluralize(resourceName),
		resourceHeading = '',
		resourceDescription = '',
		creationDescription = '',
		resourceList = [],
		createForm = undefined,
		createFormBody = undefined,
		createAction = '?/create',
		apiError = null,
		resourceEmptyCreateDescription = `Create the first ${dePluralize(resourceName)} to get started.`,
		headerActions = undefined,
		groupBy = undefined,
		listItemCard
	}: ResourceCollectionProps<T, TValues> = $props();

	let createOpen = $state(false);
	let openGroups = $state<Record<string, boolean>>({});

	$effect(() => {
		if (createForm?.createError || createForm?.createSuccess) {
			createOpen = true;
		}
	});

	function toggleCreateCard(nextState?: boolean) {
		createOpen = nextState ?? !createOpen;
	}

	function toggleGroup(groupLabel: string) {
		openGroups = {
			...openGroups,
			[groupLabel]: !openGroups[groupLabel]
		};
	}

	const groupedResources = $derived.by(() => {
		if (!groupBy) {
			return [];
		}

		const groups = new Map<string, T[]>();

		for (const resourceItem of resourceList) {
			const groupLabel = groupBy(resourceItem) ?? 'Other';
			const existingGroup = groups.get(groupLabel);

			if (existingGroup) {
				existingGroup.push(resourceItem);
				continue;
			}

			groups.set(groupLabel, [resourceItem]);
		}

		return Array.from(groups.entries())
			.sort(([leftLabel], [rightLabel]) =>
				leftLabel.localeCompare(rightLabel, undefined, { numeric: true, sensitivity: 'base' })
			)
			.map(([label, items]) => ({ label, items }));
	});

	$effect(() => {
		if (!groupBy) {
			if (Object.keys(untrack(() => openGroups)).length > 0) {
				openGroups = {};
			}
			return;
		}

		const currentOpenGroups = untrack(() => openGroups);
		const nextOpenGroups = Object.fromEntries(
			groupedResources.map((group, index) => [group.label, currentOpenGroups[group.label] ?? index === 0])
		);

		const hasChanged =
			Object.keys(currentOpenGroups).length !== Object.keys(nextOpenGroups).length ||
			Object.entries(nextOpenGroups).some(
				([label, isOpen]) => currentOpenGroups[label] !== isOpen
			);

		if (hasChanged) {
			openGroups = nextOpenGroups;
		}
	});
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-2">
			<p class="text-sm font-medium tracking-[0.22em] text-[var(--color-secondary-700)] uppercase">
				{resourceName}
			</p>
			<h1 class="text-4xl font-semibold tracking-tight text-slate-950">{resourceHeading}</h1>
			<p class="max-w-2xl text-base leading-7 text-slate-600">
				{resourceDescription}
			</p>
		</div>
		<div class="flex flex-wrap items-center gap-3">
			<CounterBadge counter={resourceList.length} pl={resourceName} sl={resourceSingleName} />
			{@render headerActions?.()}
		</div>
	</div>

	<Card
		class={cn(
			'overflow-hidden transition-all duration-300',
			createOpen &&
				'border-[var(--color-secondary-200)] shadow-lg shadow-[color-mix(in_oklab,var(--color-secondary-200)_35%,transparent)]'
		)}
	>
		<CardHeader class="gap-4">
			<div class="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
				<div class="space-y-2">
					<div class="space-y-1">
						<CardTitle>Create {resourceSingleName}</CardTitle>
						<CardDescription>{creationDescription}</CardDescription>
					</div>
				</div>
				<div class="flex shrink-0 gap-2">
					{#if createOpen}
						<Button variant="ghost" size="sm" onclick={() => toggleCreateCard(false)}>Close</Button>
					{:else}
						<Button size="sm" onclick={() => toggleCreateCard(true)}>
							Add New {capitalize(resourceSingleName)}
						</Button>
					{/if}
				</div>
			</div>
		</CardHeader>

		{#if createOpen}
			<div transition:slide={{ duration: 220 }}>
				<CardContent class="border-t border-slate-100 pt-6">
					<form method="POST" class="grid gap-4 md:grid-cols-2">
						{@render createFormBody?.()}
						<div class="space-y-3 md:col-span-2">
							{#if createForm?.createError}
								<p class="text-sm text-[var(--color-danger-700)]">{createForm.createError}</p>
							{/if}
							{#if createForm?.createSuccess}
								<p class="text-sm text-[var(--color-success-700)]">{createForm.createSuccess}</p>
							{/if}
							<div class="flex flex-col gap-3 sm:flex-row sm:justify-end">
								<Button variant="outline" type="button" onclick={() => toggleCreateCard(false)}>
									Cancel
								</Button>
								<Button type="submit" formAction={createAction}>
									Create {capitalize(resourceSingleName)}
								</Button>
							</div>
						</div>
					</form>
				</CardContent>
			</div>
		{/if}
	</Card>

	{#if apiError}
		<Card class="border-amber-200 bg-amber-50">
			<CardContent class="p-5 text-sm text-amber-900">{apiError}</CardContent>
		</Card>
	{/if}

	{#if resourceList.length === 0}
		<Card>
			<CardHeader>
				<CardTitle>No {resourceName} yet</CardTitle>
				<CardDescription>{resourceEmptyCreateDescription}</CardDescription>
			</CardHeader>
		</Card>
	{:else}
		{#if groupBy}
			<div class="space-y-8">
				{#each groupedResources as group}
					<section class="space-y-4">
						<button
							type="button"
							class="flex w-full items-center justify-between gap-4 rounded-2xl border border-slate-200 bg-white px-4 py-3 text-left transition hover:border-slate-300 hover:bg-slate-50"
							aria-expanded={openGroups[group.label]}
							onclick={() => toggleGroup(group.label)}
						>
							<div class="space-y-1">
								<h2 class="text-lg font-semibold tracking-tight text-slate-950">{group.label}</h2>
								<p class="text-sm text-slate-500">
									{group.items.length} {group.items.length === 1 ? resourceSingleName : resourceName}
								</p>
							</div>
							<span class="text-sm font-medium text-slate-500">
								{openGroups[group.label] ? 'Hide' : 'Show'}
							</span>
						</button>
						{#if openGroups[group.label]}
							<div transition:slide={{ duration: 220 }}>
								<div
									class={cn(
										'grid gap-4 transition-all duration-300 md:grid-cols-2',
										!createOpen && 'xl:grid-cols-3'
									)}
								>
									{#each group.items as resourceItem (resourceItem.id)}
										{@render listItemCard(resourceItem)}
									{/each}
								</div>
							</div>
						{/if}
					</section>
				{/each}
			</div>
		{:else}
			<div
				class={cn(
					'grid gap-4 transition-all duration-300 md:grid-cols-2',
					!createOpen && 'xl:grid-cols-3'
				)}
			>
				{#each resourceList as resourceItem (resourceItem.id)}
					{@render listItemCard(resourceItem)}
				{/each}
			</div>
		{/if}
	{/if}
</div>
