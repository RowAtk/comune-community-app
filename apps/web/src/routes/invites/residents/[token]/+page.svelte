<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

let { data, form } = $props();

	const isInviteCreator = $derived(
		Boolean(data.auth && data.invite?.invited_by === data.auth.user.id)
	);
</script>

<svelte:head>
	<title>Resident Invite | Comune</title>
</svelte:head>

<div class="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(14,165,233,0.16),_transparent_28%),linear-gradient(180deg,_#f8fafc_0%,_#eef2ff_100%)]">
	<div class="mx-auto flex min-h-screen max-w-6xl flex-col justify-center gap-8 px-6 py-16 lg:px-10">
		<div class="max-w-3xl space-y-4">
			<p class="text-sm font-medium uppercase tracking-[0.24em] text-[var(--color-secondary-700)]">
				Resident Invitation
			</p>
			<h1 class="text-5xl font-semibold tracking-tight text-slate-950">
				Claim your resident access in Comune.
			</h1>
			<p class="max-w-2xl text-lg leading-8 text-slate-600">
				Sign in or create an account to link yourself to the resident record your community team prepared for you.
			</p>
		</div>

		{#if data.apiError || !data.invite}
			<Card class="max-w-3xl border-amber-200 bg-amber-50">
				<CardContent class="p-6 text-sm text-amber-900">
					{data.apiError ?? 'This resident invitation is not available.'}
				</CardContent>
			</Card>
		{:else}
			<div class="grid gap-6 lg:grid-cols-[minmax(0,1.1fr)_minmax(340px,420px)]">
				<Card class="border-slate-200/70 bg-white/90 backdrop-blur">
					<CardHeader>
						<CardTitle>{data.invite.resident_name}</CardTitle>
						<CardDescription>
							{data.invite.resident_type} for {data.invite.community}, managed by {data.invite.organization}.
						</CardDescription>
					</CardHeader>
					<CardContent class="space-y-4 text-sm text-slate-600">
						<div class="rounded-2xl border border-slate-200 bg-slate-50 p-4">
							<p class="text-xs font-medium uppercase tracking-[0.2em] text-slate-400">Unit</p>
							<p class="mt-2 text-base font-semibold text-slate-950">{data.invite.unit_number}</p>
						</div>
						<div class="rounded-2xl border border-slate-200 bg-slate-50 p-4">
							<p class="text-xs font-medium uppercase tracking-[0.2em] text-slate-400">Household access</p>
							<p class="mt-2 text-base font-semibold text-slate-950">
								{data.invite.household_role || 'Default resident access'}
							</p>
							{#if data.invite.household_name}
								<p class="mt-2 text-sm text-slate-600">{data.invite.household_name}</p>
							{/if}
						</div>
						<div class="rounded-2xl border border-slate-200 bg-slate-50 p-4">
							<p class="text-xs font-medium uppercase tracking-[0.2em] text-slate-400">Expires</p>
							<p class="mt-2 text-base font-semibold text-slate-950">
								{new Date(data.invite.expires_at).toLocaleString()}
							</p>
						</div>
					</CardContent>
				</Card>

				{#if data.auth}
					<Card class="border-slate-200/70 bg-white">
						<CardHeader>
							<CardTitle>Finish linking your account</CardTitle>
							<CardDescription>
								You are signed in as {data.auth.user.email}. Accepting this invite will link that account to this resident record.
							</CardDescription>
						</CardHeader>
						<CardContent class="space-y-4">
							<form method="POST" class="space-y-4">
								{#if isInviteCreator}
									<div class="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900">
										This invite was created by your account. Send this URL to the resident instead of opening it yourself.
									</div>
								{/if}
								{#if form?.acceptError}
									<p class="text-sm text-[var(--color-danger-700)]">{form.acceptError}</p>
								{/if}
								<Button type="submit" class="w-full" formAction="?/accept" disabled={!!isInviteCreator}>
									Accept Resident Invite
								</Button>
							</form>
						</CardContent>
					</Card>
				{:else}
					<div class="grid gap-6">
						<Card class="border-slate-200/70 bg-white">
							<CardHeader>
								<CardTitle>Log In</CardTitle>
								<CardDescription>Use an existing account and we will link it after login.</CardDescription>
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
									{#if form?.acceptError}
										<p class="text-sm text-[var(--color-danger-700)]">{form.acceptError}</p>
									{/if}
									<Button type="submit" class="w-full" formAction="?/login">Log In And Accept</Button>
								</form>
							</CardContent>
						</Card>

						<Card class="border-slate-200/70 bg-white">
							<CardHeader>
								<CardTitle>Create Account</CardTitle>
								<CardDescription>Create an account and link it to this resident in one flow.</CardDescription>
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
										<Button type="submit" class="w-full" formAction="?/signup">Create Account And Accept</Button>
									</div>
								</form>
							</CardContent>
						</Card>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>
