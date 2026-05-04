<script lang="ts">
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
	import { Building2, Home, LogIn, UserPlus, Users } from 'lucide-svelte';

	let { form } = $props();

	const operationalAreas = [
		{
			label: 'Foundation ready',
			title: 'Units',
			copy: 'Model the physical layout first so every later workflow has a real place to land.',
			icon: Building2,
			highlight: true
		},
		{
			label: 'Foundation ready',
			title: 'Households',
			copy: 'Group people the way communities actually manage occupancy and responsibility.',
			icon: Home,
			highlight: false
		},
		{
			label: 'Current growth lane',
			title: 'Invoicing',
			copy: 'Move next into maintenance dues, overdue visibility, and resident billing clarity.',
			icon: Users,
			highlight: false
		}
	];
</script>

<svelte:head>
	<title>Comune</title>
</svelte:head>

<div class="min-h-screen">
	<div
		class="mx-auto flex min-h-screen max-w-6xl flex-col justify-center gap-10 px-6 py-16 lg:px-10"
	>
		<div class="grid gap-6 xl:grid-cols-[minmax(0,1.15fr)_minmax(300px,0.85fr)] xl:items-end">
			<div class="max-w-3xl space-y-5">
				<div
					class="theme-hero-chip inline-flex rounded-full px-4 py-2 text-xs font-semibold tracking-[0.26em] uppercase"
				>
					Comune Workspace
				</div>
				<h1 class="max-w-2xl text-5xl font-semibold tracking-tight text-[var(--color-app-text)]">
					Run gated-community operations from one connected workspace.
				</h1>
				<p class="max-w-2xl text-lg leading-8 text-[var(--color-app-muted)]">
					Move from organization setup to day-to-day community operations in one calm, tenant-safe
					workspace built for real staff, residents, and managers.
				</p>
			</div>

			<Card class="theme-soft-card border-0">
				<CardHeader>
					<CardTitle>Choose your starting lane</CardTitle>
					<CardDescription>
						Sign in if you already work in a community. Create an account if you are setting up a
						new organization or portfolio.
					</CardDescription>
				</CardHeader>
				<CardContent class="grid gap-3 sm:grid-cols-2 xl:grid-cols-1">
					<div
						class="rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] p-4"
					>
						<p class="text-sm font-medium text-[var(--color-app-text)]">Fastest route back in</p>
						<p class="mt-1 text-sm text-[var(--color-app-muted)]">
							Log in to continue where your organization or community team left off.
						</p>
					</div>
					<div
						class="rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] p-4"
					>
						<p class="text-sm font-medium text-[var(--color-app-text)]">New setup</p>
						<p class="mt-1 text-sm text-[var(--color-app-muted)]">
							Create an account first, then move into organization and community setup.
						</p>
					</div>
				</CardContent>
			</Card>
		</div>

		<div class="grid gap-6 lg:grid-cols-[minmax(0,1.1fr)_minmax(320px,420px)]">
			<Card>
				<CardHeader>
					<CardTitle>What The Product Handles Today</CardTitle>
					<CardDescription>
						The current product is organized around a few real operational lanes instead of a broad,
						generic dashboard surface.
					</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4">
					<div class="space-y-3">
						<p class="text-sm tracking-[0.18em] text-[var(--color-app-subtle)] uppercase">
							Operational lanes
						</p>
						<div class="grid gap-4 md:grid-cols-3">
							{#each operationalAreas as area (area.title)}
								<div
									class={area.highlight
										? 'theme-highlight-card rounded-2xl p-5'
										: 'theme-soft-card rounded-2xl p-5'}
								>
									<svelte:component
										this={area.icon}
										class={area.highlight
											? 'h-5 w-5 text-[color:rgba(253,250,244,0.82)]'
											: 'h-5 w-5 text-[var(--color-brand-600)]'}
									/>
									<p
										class={area.highlight
											? 'mt-4 text-sm text-[color:rgba(253,250,244,0.76)]'
											: 'mt-4 text-sm text-[var(--color-app-muted)]'}
									>
										{area.label}
									</p>
									<p
										class={area.highlight
											? 'mt-3 text-2xl font-semibold'
											: 'mt-3 text-2xl font-semibold text-[var(--color-app-text)]'}
									>
										{area.title}
									</p>
									<p
										class={area.highlight
											? 'mt-2 text-sm leading-6 text-[color:rgba(253,250,244,0.78)]'
											: 'mt-2 text-sm leading-6 text-[var(--color-app-muted)]'}
									>
										{area.copy}
									</p>
								</div>
							{/each}
						</div>
					</div>

					<div
						class="rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] p-5"
					>
						<p class="text-sm tracking-[0.18em] text-[var(--color-app-subtle)] uppercase">
							How to enter the app
						</p>
						<ol class="mt-4 space-y-2 text-sm text-[var(--color-app-muted)]">
							<li>1. Log in if you already belong to an organization or community.</li>
							<li>2. Create an account if you are starting a new setup.</li>
							<li>3. Use a resident invite token only when an admin sent one directly to you.</li>
						</ol>
					</div>
				</CardContent>
			</Card>

			<div class="grid gap-6">
				<Card>
					<CardHeader>
						<div class="flex items-center gap-3">
							<LogIn class="h-5 w-5 text-[var(--color-brand-700)]" />
							<div>
								<CardTitle>Log In</CardTitle>
								<CardDescription
									>Continue into the operational workspace with your existing account.</CardDescription
								>
							</div>
						</div>
					</CardHeader>
					<CardContent>
						<form method="POST" class="space-y-4">
							<div class="space-y-2">
								<Label for="email">Email</Label>
								<Input
									id="email"
									name="email"
									type="email"
									required
									value={form?.loginValues?.email ?? ''}
								/>
							</div>
							<div class="space-y-2">
								<Label for="password">Password</Label>
								<Input id="password" name="password" type="password" required />
							</div>
							{#if form?.loginError}
								<p class="text-sm text-[var(--color-danger-700)]">{form.loginError}</p>
							{/if}
							<Button type="submit" class="w-full" formAction="?/login">Log In</Button>
						</form>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<div class="flex items-center gap-3">
							<UserPlus class="h-5 w-5 text-[var(--color-brand-700)]" />
							<div>
								<CardTitle>Create Account</CardTitle>
								<CardDescription
									>Create a new account for organization and community setup.</CardDescription
								>
							</div>
						</div>
					</CardHeader>
					<CardContent>
						<form method="POST" class="grid gap-4 md:grid-cols-2">
							<div class="space-y-2 md:col-span-2">
								<Label for="signup_email">Email</Label>
								<Input
									id="signup_email"
									name="signup_email"
									type="email"
									required
									value={form?.signupValues?.email ?? ''}
								/>
							</div>
							<div class="space-y-2">
								<Label for="first_name">First Name</Label>
								<Input
									id="first_name"
									name="first_name"
									value={form?.signupValues?.first_name ?? ''}
								/>
							</div>
							<div class="space-y-2">
								<Label for="last_name">Last Name</Label>
								<Input
									id="last_name"
									name="last_name"
									value={form?.signupValues?.last_name ?? ''}
								/>
							</div>
							<div class="space-y-2">
								<Label for="phone">Phone</Label>
								<Input id="phone" name="phone" value={form?.signupValues?.phone ?? ''} />
							</div>
							<div class="space-y-2">
								<Label for="signup_password">Password</Label>
								<Input id="signup_password" name="signup_password" type="password" required />
							</div>
							<div class="md:col-span-2">
								{#if form?.signupError}
									<p class="text-sm text-[var(--color-danger-700)]">{form.signupError}</p>
								{/if}
							</div>
							<div class="md:col-span-2">
								<Button type="submit" class="w-full" formAction="?/signup">Create Account</Button>
							</div>
						</form>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<div>
							<CardTitle>Have An Invite Token?</CardTitle>
							<CardDescription>
								Paste a resident invite token to jump directly into the claim flow with fewer
								decisions.
							</CardDescription>
						</div>
					</CardHeader>
					<CardContent>
						<form method="GET" action="/invites/residents" class="space-y-4">
							<div class="space-y-2">
								<Label for="resident_invite_token">Resident Invite Token</Label>
								<Input
									id="resident_invite_token"
									name="token"
									placeholder="Paste token here"
									required
								/>
							</div>
							<Button type="submit" class="w-full" variant="outline">Open Invite</Button>
						</form>
					</CardContent>
				</Card>
			</div>
		</div>
	</div>
</div>
