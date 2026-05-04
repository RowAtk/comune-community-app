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

	const values = $derived(form?.values ?? data.defaultValues);
</script>

<svelte:head>
	<title>New Organization | Comune</title>
</svelte:head>

<AppShell
	eyebrow="Organization Setup"
	title="Create a new organization"
	description="Set up the top-level organization first, then you can attach one or more communities to it."
>
	{#snippet actions()}
		<Button href="/dashboard" variant="ghost">Back to Dashboard</Button>
	{/snippet}

	<Card>
		<CardHeader>
			<CardTitle>Organization Details</CardTitle>
			<CardDescription>
				Keep the first step focused on the essential operational details and defer the rest.
			</CardDescription>
		</CardHeader>
		<CardContent>
			<form method="POST" class="grid gap-4 md:grid-cols-2">
				<div class="space-y-2 md:col-span-2">
					<Label for="name">Organization Name</Label>
					<Input
						id="name"
						name="name"
						required
						value={values.name}
						placeholder="Kingston Estates Management"
					/>
				</div>
				<div class="space-y-2">
					<Label for="slug">Slug</Label>
					<Input
						id="slug"
						name="slug"
						required
						value={values.slug}
						placeholder="kingston-estates"
					/>
				</div>
				<div class="space-y-2">
					<Label for="legal_name">Legal Name</Label>
					<Input
						id="legal_name"
						name="legal_name"
						value={values.legal_name}
						placeholder="Kingston Estates Management Ltd."
					/>
				</div>
				<div class="space-y-2">
					<Label for="billing_email">Billing Email</Label>
					<Input
						id="billing_email"
						name="billing_email"
						type="email"
						value={values.billing_email}
						placeholder="billing@example.com"
					/>
				</div>
				<div class="space-y-2">
					<Label for="phone">Phone</Label>
					<Input id="phone" name="phone" value={values.phone} placeholder="+1 876 555 0100" />
				</div>
				<div class="space-y-2">
					<Label for="country_code">Country Code</Label>
					<Input
						id="country_code"
						name="country_code"
						value={values.country_code}
						placeholder="JM"
					/>
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
				{#if form?.error}
					<p class="text-sm text-[var(--color-danger-700)] md:col-span-2">{form.error}</p>
				{/if}
				<div class="flex flex-wrap justify-end gap-3 pt-2 md:col-span-2">
					<Button href="/dashboard" variant="outline">Cancel</Button>
					<Button type="submit">Create Organization</Button>
				</div>
			</form>
		</CardContent>
	</Card>
</AppShell>
