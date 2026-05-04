<script lang="ts">
	import AppShell from '$lib/components/app-shell.svelte';
	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	let { data, form } = $props();

	const dashboard = $derived(form?.dashboard ?? data.dashboard);
	const values = $derived(form?.values ?? data.defaultValues);
</script>

<svelte:head>
	<title>New Community | Comune</title>
</svelte:head>

<AppShell
	eyebrow="Community Setup"
	title="Create a new community"
	description="Choose the organization it belongs to, then create the community workspace."
>
	{#snippet actions()}
		<Button href="/dashboard" variant="ghost">Back to Dashboard</Button>
	{/snippet}

	<Card>
		<CardHeader>
			<CardTitle>Community Details</CardTitle>
			<CardDescription>
				Communities live inside an organization and become their own operational workspace.
			</CardDescription>
		</CardHeader>
		<CardContent>
			<form method="POST" class="grid gap-4 md:grid-cols-2">
				<div class="space-y-2 md:col-span-2">
					<Label for="organization_id">Organization</Label>
					<select
						id="organization_id"
						name="organization_id"
						class="flex h-11 w-full rounded-xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm transition outline-none focus:border-[var(--color-brand-300)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-brand-200)_45%,transparent)]"
						required
					>
						<option value="" disabled selected={!values.organization_id}
							>Select an organization</option
						>
						{#each dashboard.organizations as organization (organization.organization_id)}
							<option
								value={organization.organization_id}
								selected={organization.organization_id === values.organization_id}
							>
								{organization.name}
							</option>
						{/each}
					</select>
				</div>
				<div class="space-y-2 md:col-span-2">
					<Label for="name">Community Name</Label>
					<Input
						id="name"
						name="name"
						required
						value={values.name}
						placeholder="Harbour View Residences"
					/>
				</div>
				<div class="space-y-2">
					<Label for="slug">Slug</Label>
					<Input id="slug" name="slug" required value={values.slug} placeholder="harbour-view" />
				</div>
				<div class="space-y-2">
					<Label for="timezone">Timezone</Label>
					<Input
						id="timezone"
						name="timezone"
						value={values.timezone}
						placeholder="America/Jamaica"
					/>
				</div>
				<div class="space-y-2 md:col-span-2">
					<Label for="address">Address</Label>
					<Input
						id="address"
						name="address"
						value={values.address}
						placeholder="24 Harbour Road, Kingston"
					/>
				</div>
				{#if form?.error}
					<p class="text-sm text-[var(--color-danger-700)] md:col-span-2">{form.error}</p>
				{/if}
				<div class="flex flex-wrap justify-end gap-3 pt-2 md:col-span-2">
					<Button href="/dashboard" variant="outline">Cancel</Button>
					<Button type="submit" disabled={dashboard.organizations.length === 0}
						>Create Community</Button
					>
				</div>
				{#if dashboard.organizations.length === 0}
					<p class="text-sm text-[var(--color-app-muted)] md:col-span-2">
						Create an organization first before adding a community.
					</p>
				{/if}
			</form>
		</CardContent>
	</Card>
</AppShell>
