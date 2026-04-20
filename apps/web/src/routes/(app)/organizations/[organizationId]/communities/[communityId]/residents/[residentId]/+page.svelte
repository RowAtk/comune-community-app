<script lang="ts">
	import { page } from '$app/state';
	import { Badge } from '$lib/components/ui/badge';
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

	const communityBasePath = $derived(
		`/organizations/${page.params.organizationId}/communities/${page.params.communityId}`
	);

	const inviteBaseUrl = $derived(`${page.url.origin}/invites/residents`);
	const hasActiveInvite = $derived(
		data.invitations.some(
			(invitation) =>
				!invitation.accepted_at && new Date(invitation.expires_at).getTime() > Date.now()
		)
	);
	const displayEmail = $derived(
		data.resident
			? data.resident.user_id
				? data.resident.linked_user_email || data.resident.email
				: data.resident.email
			: ''
	);
	const displayPhone = $derived(
		data.resident
			? data.resident.user_id
				? data.resident.linked_user_phone || data.resident.phone
				: data.resident.phone
			: ''
	);

	function statusVariant(status: string) {
		return status === 'ACTIVE' ? 'success' : 'outline';
	}
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-2">
			<p class="text-sm font-medium tracking-[0.22em] text-[var(--color-secondary-700)] uppercase">
				Resident
			</p>
			<h1 class="text-4xl font-semibold tracking-tight text-slate-950">
				{data.resident
					? `${data.resident.first_name} ${data.resident.last_name}`
					: 'Resident details'}
			</h1>
			<p class="max-w-2xl text-base leading-7 text-slate-600">
				{#if data.resident?.user_id}
					Review the linked account and resident record for this community member.
				{:else}
					Generate a token or invite link so this resident can claim their account and join the
					community workspace.
				{/if}
			</p>
		</div>
		<div class="flex flex-wrap gap-3">
			<Button href={`${communityBasePath}/residents`} variant="outline">Back To Residents</Button>
			{#if data.resident}
				<Badge variant={statusVariant(data.resident.status)}>{data.resident.status}</Badge>
			{/if}
		</div>
	</div>

	{#if data.apiError || !data.resident}
		<Card class="border-amber-200 bg-amber-50">
			<CardContent class="p-5 text-sm text-amber-900">
				{data.apiError ?? 'Resident details are not available.'}
			</CardContent>
		</Card>
	{:else}
		<div class="grid gap-5 md:grid-cols-3">
			<Card class="bg-slate-950 text-white">
				<CardHeader>
					<CardDescription class="text-slate-300">Resident type</CardDescription>
					<CardTitle class="text-white">{data.resident.resident_type}</CardTitle>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader>
					<CardDescription>Unit</CardDescription>
					<CardTitle>{data.unit?.unit_number || 'Unknown unit'}</CardTitle>
					{#if data.household?.name}
						<CardDescription class="pt-2">{data.household.name}</CardDescription>
					{/if}
				</CardHeader>
			</Card>
			<Card>
				<CardHeader>
					<CardDescription>Linked access</CardDescription>
					<CardTitle>
						{#if data.resident.user_id}
							Linked to user "{data.resident.linked_user_name ||
								data.resident.linked_user_email ||
								'linked account'}"
						{:else if hasActiveInvite}
							Invite pending
						{:else}
							Account not linked
						{/if}
					</CardTitle>
					{#if data.resident.user_id && data.resident.linked_user_name && data.resident.linked_user_email}
						<CardDescription class="pt-2">{data.resident.linked_user_email}</CardDescription>
					{/if}
				</CardHeader>
				{#if data.resident.user_id}
					<CardContent class="space-y-3 pt-0">
						<p class="text-sm text-slate-600">
							Remove the linked account if this resident was connected to the wrong user or needs to
							be reassigned.
						</p>
						<form
							method="POST"
							class="space-y-3"
							onsubmit={(event) => {
								if (!confirm('Are you sure you want to unlink this resident from their account?')) {
									event.preventDefault();
								}
							}}
						>
							{#if form?.unlinkError}
								<p class="text-sm text-[var(--color-danger-700)]">{form.unlinkError}</p>
							{/if}
							<Button type="submit" variant="outline" class="w-full" formAction="?/unlinkUser">
								Unlink Resident Account
							</Button>
						</form>
					</CardContent>
				{:else if !hasActiveInvite}
					<CardContent class="pt-0 text-sm text-slate-600">
						Generate a new invite to link this resident to an account.
					</CardContent>
				{/if}
			</Card>
		</div>

		<div class="grid gap-5 xl:grid-cols-[minmax(0,1fr)_minmax(360px,420px)]">
			<Card>
				<CardHeader>
					<CardTitle>Invitation History</CardTitle>
					<CardDescription>
						Share the token directly or send the generated invite link to the resident.
					</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4">
					{#if form?.createdInvitation}
						<div
							class="rounded-2xl border border-[var(--color-success-200)] bg-[var(--color-success-50)] p-4"
						>
							<p class="text-sm font-medium text-[var(--color-success-700)]">Newest invitation</p>
							<p class="mt-2 font-mono text-sm break-all text-slate-900">
								{form.createdInvitation.token}
							</p>
							<p class="mt-3 text-sm break-all text-slate-700">
								{inviteBaseUrl}/{form.createdInvitation.token}
							</p>
						</div>
					{/if}

					{#if data.invitations.length === 0}
						<p class="text-sm text-slate-600">No resident invitations have been generated yet.</p>
					{:else}
						<div class="space-y-4">
							{#each data.invitations as invitation}
								<div class="rounded-2xl border border-slate-200 bg-slate-50 p-4">
									<div class="flex flex-wrap items-start justify-between gap-3">
										<div class="space-y-2">
											<p class="text-xs font-medium tracking-[0.2em] text-slate-400 uppercase">
												Token
											</p>
											<p class="font-mono text-[13px] break-all text-slate-900">
												{invitation.token}
											</p>
										</div>
										<Badge variant={invitation.accepted_at ? 'success' : 'outline'}>
											{invitation.accepted_at ? 'Accepted' : 'Pending'}
										</Badge>
									</div>
									<div class="mt-4 space-y-2 text-sm text-slate-600">
										<p>
											Link: <span class="break-all text-slate-900"
												>{inviteBaseUrl}/{invitation.token}</span
											>
										</p>
										<p>Expires: {new Date(invitation.expires_at).toLocaleString()}</p>
										{#if invitation.accepted_at}
											<p>Accepted: {new Date(invitation.accepted_at).toLocaleString()}</p>
										{/if}
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</CardContent>
			</Card>

			{#if !data.resident.user_id}
				<Card>
					<CardHeader>
						<CardTitle>Generate Invite</CardTitle>
						<CardDescription>
							Create a fresh resident claim token. The default expiry is seven days if you leave this
							blank.
						</CardDescription>
					</CardHeader>
					<CardContent>
						<form method="POST" class="space-y-4">
							<div class="space-y-2">
								<Label for="expires_at">Expiry Date</Label>
								<Input
									id="expires_at"
									name="expires_at"
									type="date"
									value={form?.values?.expires_at ?? ''}
								/>
							</div>
							<div class="space-y-3">
								{#if form?.createError}
									<p class="text-sm text-[var(--color-danger-700)]">{form.createError}</p>
								{/if}
								{#if form?.createSuccess}
									<p class="text-sm text-[var(--color-success-700)]">{form.createSuccess}</p>
								{/if}
								<Button type="submit" class="w-full" formAction="?/createInvite">
									Generate Resident Invite
								</Button>
							</div>
						</form>
					</CardContent>
				</Card>
			{/if}
		</div>

		<Card>
			<CardHeader>
				<CardTitle>Contact Info On File</CardTitle>
				<CardDescription>
					{#if data.resident.user_id}
						This contact information is currently sourced from the linked user account.
					{:else}
						This contact information is currently stored on the resident record.
					{/if}
				</CardDescription>
			</CardHeader>
			<CardContent class="grid gap-5 md:grid-cols-2">
				<div>
					<p class="text-xs font-medium uppercase tracking-[0.2em] text-slate-400">Email</p>
					<p class="mt-2 text-slate-900">{displayEmail || 'No email on file'}</p>
				</div>
				<div>
					<p class="text-xs font-medium uppercase tracking-[0.2em] text-slate-400">Phone</p>
					<p class="mt-2 text-slate-900">{displayPhone || 'No phone on file'}</p>
				</div>
			</CardContent>
		</Card>
	{/if}
</div>
