<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Building2, Home, LogIn, ShieldCheck, UserPlus, Users } from 'lucide-svelte';

	let { data, form } = $props();
</script>

<svelte:head>
	<title>Comune</title>
</svelte:head>

<div class="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(14,165,233,0.16),_transparent_28%),linear-gradient(180deg,_#f8fafc_0%,_#eef2ff_100%)]">
	<div class="mx-auto flex min-h-screen max-w-6xl flex-col justify-center gap-10 px-6 py-16 lg:px-10">
		<div class="max-w-3xl space-y-5">
			<p class="text-sm font-medium uppercase tracking-[0.24em] text-[var(--color-secondary-700)]">
				Comune Workspace
			</p>
			<h1 class="max-w-2xl text-5xl font-semibold tracking-tight text-slate-950">
				Run gated-community operations from one connected workspace.
			</h1>
			<p class="max-w-2xl text-lg leading-8 text-slate-600">
				Authentication is now wired through the web app, so protected community screens only load when
				there is a valid backend session.
			</p>
		</div>

		<div class="grid gap-6 lg:grid-cols-[minmax(0,1.1fr)_minmax(320px,420px)]">
			<Card class="border-slate-200/70 bg-white/80 backdrop-blur">
				<CardHeader>
					<CardTitle>App Foundation</CardTitle>
					<CardDescription>
						SvelteKit server loads and actions are handling auth, mutations, and route protection without
						introducing a separate state library.
					</CardDescription>
				</CardHeader>
				<CardContent>
					<div class="grid gap-4 md:grid-cols-3">
						<div class="rounded-2xl bg-slate-950 p-5 text-white">
							<Building2 class="h-5 w-5 text-slate-300" />
							<p class="mt-4 text-sm text-slate-300">Vertical slice</p>
							<p class="mt-3 text-2xl font-semibold">Units</p>
						</div>
						<div class="rounded-2xl border border-slate-200 bg-slate-50 p-5">
							<Home class="h-5 w-5 text-slate-500" />
							<p class="mt-4 text-sm text-slate-500">Scaffolded next</p>
							<p class="mt-3 text-2xl font-semibold text-slate-900">Households</p>
						</div>
						<div class="rounded-2xl border border-slate-200 bg-slate-50 p-5">
							<Users class="h-5 w-5 text-slate-500" />
							<p class="mt-4 text-sm text-slate-500">Scaffolded next</p>
							<p class="mt-3 text-2xl font-semibold text-slate-900">Residents</p>
						</div>
					</div>
				</CardContent>
			</Card>

			{#if data.auth}
				<Card class="border-slate-200/70 bg-white">
					<CardHeader>
						<div class="flex items-center justify-between gap-4">
							<div>
								<CardTitle>Authenticated</CardTitle>
								<CardDescription>
									Signed in as {data.auth.user.first_name || data.auth.user.email}
								</CardDescription>
							</div>
							<Badge variant="success">
								<ShieldCheck class="mr-1 h-3.5 w-3.5" />
								Session active
							</Badge>
						</div>
					</CardHeader>
					<CardContent class="space-y-5">
						<div class="space-y-2 rounded-2xl border border-slate-200 bg-slate-50 p-4">
							<p class="text-sm font-medium text-slate-900">{data.auth.user.email}</p>
							<p class="text-sm text-slate-600">
								Session expires {new Date(data.auth.session.expires_at).toLocaleString()}
							</p>
						</div>

						<form method="POST" class="space-y-4">
							<div class="space-y-2">
								<Label for="organizationId">Organization ID</Label>
								<Input id="organizationId" name="organizationId" placeholder="8c49... or demo-org" required />
							</div>
							<div class="space-y-2">
								<Label for="communityId">Community ID</Label>
								<Input id="communityId" name="communityId" placeholder="f39d... or main-estate" required />
							</div>
							{#if form?.error}
								<p class="text-sm text-[var(--color-danger-700)]">{form.error}</p>
							{/if}
							<Button type="submit" class="w-full" formAction="?/openWorkspace">Open Workspace</Button>
						</form>

						<form method="POST">
							<Button type="submit" variant="ghost" class="w-full" formAction="?/logout">Log Out</Button>
						</form>
					</CardContent>
				</Card>
			{:else}
				<div class="grid gap-6">
					<Card class="border-slate-200/70 bg-white">
						<CardHeader>
							<div class="flex items-center gap-3">
								<LogIn class="h-5 w-5 text-[var(--color-secondary-700)]" />
								<div>
									<CardTitle>Log In</CardTitle>
									<CardDescription>Use your existing account to unlock the workspace.</CardDescription>
								</div>
							</div>
						</CardHeader>
						<CardContent>
							<form method="POST" class="space-y-4">
								<div class="space-y-2">
									<Label for="email">Email</Label>
									<Input id="email" name="email" type="email" required value={form?.loginValues?.email ?? ''} />
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

					<Card class="border-slate-200/70 bg-white">
						<CardHeader>
							<div class="flex items-center gap-3">
								<UserPlus class="h-5 w-5 text-[var(--color-brand-700)]" />
								<div>
									<CardTitle>Create Account</CardTitle>
									<CardDescription>Spin up a new session against the backend auth API.</CardDescription>
								</div>
							</div>
						</CardHeader>
						<CardContent>
							<form method="POST" class="grid gap-4 md:grid-cols-2">
								<div class="space-y-2 md:col-span-2">
									<Label for="signup_email">Email</Label>
									<Input id="signup_email" name="signup_email" type="email" required value={form?.signupValues?.email ?? ''} />
								</div>
								<div class="space-y-2">
									<Label for="first_name">First Name</Label>
									<Input id="first_name" name="first_name" value={form?.signupValues?.first_name ?? ''} />
								</div>
								<div class="space-y-2">
									<Label for="last_name">Last Name</Label>
									<Input id="last_name" name="last_name" value={form?.signupValues?.last_name ?? ''} />
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
				</div>
			{/if}
		</div>
	</div>
</div>
