<script lang="ts">
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
	import type { Invoice, InvoicePlan } from '$lib/api/types';

	let { data, form } = $props();

	const currency = new Intl.NumberFormat('en-JM', {
		style: 'currency',
		currency: 'JMD',
		maximumFractionDigits: 2
	});

	const dateFormatter = new Intl.DateTimeFormat('en-JM', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});

	function formatMoney(value: number) {
		return currency.format(value ?? 0);
	}

	function formatDate(value?: string | null) {
		if (!value) {
			return 'Not scheduled';
		}

		return dateFormatter.format(new Date(value));
	}

	function invoiceStatusVariant(invoice: Invoice) {
		if (invoice.is_overdue) {
			return 'warning';
		}
		if (invoice.status === 'PAID') {
			return 'success';
		}
		return 'outline';
	}

	function planStatusVariant(plan: InvoicePlan) {
		if (plan.status === 'ACTIVE') {
			return 'success';
		}
		if (plan.status === 'PAUSED') {
			return 'warning';
		}
		return 'outline';
	}

	const overdueInvoices = $derived(data.invoices.filter((invoice: Invoice) => invoice.is_overdue));

	const currentInvoices = $derived(data.invoices.filter((invoice: Invoice) => !invoice.is_overdue));
</script>

<div class="space-y-8">
	<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
		<div class="space-y-3">
			<div
				class="theme-hero-chip inline-flex rounded-full px-4 py-2 text-xs font-semibold tracking-[0.24em] uppercase"
			>
				Invoicing
			</div>
			<h1 class="text-4xl font-semibold tracking-tight text-[var(--color-app-text)]">
				Maintenance Fees And Dues
			</h1>
			<p class="max-w-3xl text-base leading-7 text-[var(--color-app-muted)]">
				Track recurring maintenance charges, ad hoc invoices, overdue households, and resident due
				dates without depending on a payment gateway.
			</p>
		</div>
		<div class="flex flex-wrap gap-3">
			<Button href="../households" variant="outline">View Households</Button>
			<Button href="../residents" variant="outline">View Residents</Button>
		</div>
	</div>

	{#if data.residentOverview}
		<div class="grid gap-5 lg:grid-cols-3">
			<Card class="theme-highlight-card border-0">
				<CardHeader>
					<CardDescription class="text-[color:rgba(253,250,244,0.78)]"
						>Next maintenance due</CardDescription
					>
					<CardTitle class="text-[var(--color-app-inverse)]"
						>{formatDate(data.residentOverview.next_due_date)}</CardTitle
					>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader>
					<CardDescription>Outstanding balance</CardDescription>
					<CardTitle>{formatMoney(data.residentOverview.total_outstanding)}</CardTitle>
				</CardHeader>
			</Card>
			<Card>
				<CardHeader>
					<CardDescription>Overdue invoices</CardDescription>
					<CardTitle>{data.residentOverview.overdue_invoice_count}</CardTitle>
				</CardHeader>
			</Card>
		</div>

		<Card>
			<CardHeader>
				<div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
					<div class="space-y-1">
						<CardTitle>Resident Maintenance View</CardTitle>
						<CardDescription>
							{data.residentOverview.resident_name} in unit {data.residentOverview.unit_number}
							{#if data.residentOverview.household_name}
								· {data.residentOverview.household_name}
							{/if}
						</CardDescription>
					</div>
					{#if data.residentOverview.overdue_invoice_count > 0}
						<Badge variant="warning">Overdue invoices need attention</Badge>
					{:else}
						<Badge variant="success">No overdue maintenance invoices</Badge>
					{/if}
				</div>
			</CardHeader>
			<CardContent class="space-y-4">
				{#if data.residentOverview.invoices.length === 0}
					<p class="text-sm text-[var(--color-app-muted)]">
						No invoices have been issued to this unit yet.
					</p>
				{:else}
					<div class="space-y-3">
						<p class="text-sm tracking-[0.18em] text-[var(--color-app-subtle)] uppercase">
							Your maintenance ledger
						</p>
						<div
							class="overflow-x-auto rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)]"
						>
							<table class="min-w-full divide-y divide-[var(--color-app-border)] text-sm">
								<thead class="bg-[var(--color-app-panel)] text-left text-[var(--color-app-subtle)]">
									<tr>
										<th class="px-4 py-3 font-medium">Invoice</th>
										<th class="px-4 py-3 font-medium">Status</th>
										<th class="px-4 py-3 font-medium">Due</th>
										<th class="px-4 py-3 font-medium">Amount</th>
										<th class="px-4 py-3 font-medium">Outstanding</th>
									</tr>
								</thead>
								<tbody class="divide-y divide-[var(--color-app-border)]">
									{#each data.residentOverview.invoices as invoice (invoice.id)}
										<tr class={invoice.is_overdue ? 'bg-amber-50/85' : ''}>
											<td class="px-4 py-3 align-top">
												<p class="font-medium text-[var(--color-app-text)]">
													{invoice.description ||
														`Maintenance invoice for ${invoice.billing_period}`}
												</p>
												<p
													class="mt-1 text-xs tracking-[0.16em] text-[var(--color-app-subtle)] uppercase"
												>
													{invoice.billing_period}
												</p>
											</td>
											<td class="px-4 py-3 align-top">
												<div class="flex flex-wrap gap-2">
													<Badge variant={invoiceStatusVariant(invoice)}>{invoice.status}</Badge>
													{#if invoice.is_overdue}
														<Badge variant="warning">Overdue</Badge>
													{/if}
												</div>
											</td>
											<td class="px-4 py-3 align-top font-medium text-[var(--color-app-text)]">
												{formatDate(invoice.due_date)}
											</td>
											<td class="px-4 py-3 align-top font-medium text-[var(--color-app-text)]">
												{formatMoney(invoice.amount)}
											</td>
											<td
												class={`px-4 py-3 align-top font-medium ${invoice.is_overdue ? 'text-amber-900' : 'text-[var(--color-app-text)]'}`}
											>
												{formatMoney(invoice.outstanding_amount)}
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</div>
				{/if}
			</CardContent>
		</Card>
	{:else if data.residentError}
		<div class="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900">
			{data.residentError}
		</div>
	{/if}

	{#if data.adminAccess}
		{#if data.adminError}
			<div class="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900">
				{data.adminError}
			</div>
		{/if}
		<div class="grid gap-6 xl:grid-cols-[minmax(0,1.25fr)_minmax(340px,0.75fr)]">
			<div class="space-y-6">
				<div class="grid gap-5 md:grid-cols-3">
					<Card class="theme-highlight-card border-0">
						<CardHeader>
							<CardDescription class="text-[color:rgba(253,250,244,0.78)]"
								>Overdue households</CardDescription
							>
							<CardTitle class="text-[var(--color-app-inverse)]"
								>{data.overdueHouseholds.length}</CardTitle
							>
						</CardHeader>
					</Card>
					<Card>
						<CardHeader>
							<CardDescription>Recurring plans</CardDescription>
							<CardTitle>{data.plans.length}</CardTitle>
						</CardHeader>
					</Card>
					<Card>
						<CardHeader>
							<CardDescription>Issued invoices</CardDescription>
							<CardTitle>{data.invoices.length}</CardTitle>
						</CardHeader>
					</Card>
				</div>

				<Card>
					<CardHeader>
						<CardTitle>Overdue Households</CardTitle>
						<CardDescription>
							Community admins can quickly see which households are behind on maintenance.
						</CardDescription>
					</CardHeader>
					<CardContent class="space-y-3">
						{#if data.overdueHouseholds.length === 0}
							<p class="text-sm text-[var(--color-app-muted)]">
								No households are currently overdue.
							</p>
						{:else}
							<div class="overflow-x-auto rounded-2xl border border-amber-200 bg-amber-50">
								<table class="min-w-full divide-y divide-amber-200 text-sm">
									<thead
										class="bg-[color:rgba(255,248,239,0.92)] text-left text-[var(--color-app-subtle)]"
									>
										<tr>
											<th class="px-4 py-3 font-medium">Household</th>
											<th class="px-4 py-3 font-medium">Context</th>
											<th class="px-4 py-3 font-medium">Earliest due</th>
											<th class="px-4 py-3 font-medium">Overdue</th>
											<th class="px-4 py-3 font-medium">Outstanding</th>
										</tr>
									</thead>
									<tbody class="divide-y divide-amber-200">
										{#each data.overdueHouseholds as household (household.unit_id)}
											<tr>
												<td class="px-4 py-3 align-top">
													<p class="font-medium text-[var(--color-app-text)]">
														{household.household_name}
													</p>
												</td>
												<td class="px-4 py-3 align-top text-[var(--color-app-muted)]">
													<p>Unit {household.unit_number}</p>
													{#if household.primary_resident_name}
														<p class="mt-1">{household.primary_resident_name}</p>
													{/if}
												</td>
												<td class="px-4 py-3 align-top font-medium text-[var(--color-app-text)]">
													{formatDate(household.earliest_due_date)}
												</td>
												<td class="px-4 py-3 align-top">
													<Badge variant="warning">{household.overdue_invoice_count} overdue</Badge>
												</td>
												<td class="px-4 py-3 align-top font-medium text-amber-900">
													{formatMoney(household.total_outstanding)}
												</td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						{/if}
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Issued Invoices</CardTitle>
						<CardDescription>
							Use this dense view to compare due dates, balances, and overdue status quickly.
						</CardDescription>
					</CardHeader>
					<CardContent class="space-y-6">
						{#if data.invoices.length === 0}
							<p class="text-sm text-[var(--color-app-muted)]">No invoices have been issued yet.</p>
						{:else}
							<div class="space-y-4">
								<div class="space-y-3">
									<div class="flex items-center justify-between gap-3">
										<div>
											<p class="text-lg font-semibold text-[var(--color-app-text)]">
												Needs attention
											</p>
											<p class="text-sm text-[var(--color-app-muted)]">
												Overdue invoices should rise to the top first.
											</p>
										</div>
										<Badge variant="warning">{overdueInvoices.length} overdue</Badge>
									</div>
									{#if overdueInvoices.length === 0}
										<div
											class="rounded-2xl border border-dashed border-[var(--color-app-border-strong)] bg-[var(--color-app-panel)] px-4 py-4 text-sm text-[var(--color-app-muted)]"
										>
											No overdue invoices right now.
										</div>
									{:else}
										<div class="overflow-x-auto rounded-2xl border border-amber-200 bg-amber-50">
											<table class="min-w-full divide-y divide-amber-200 text-sm">
												<thead
													class="bg-[color:rgba(255,248,239,0.92)] text-left text-[var(--color-app-subtle)]"
												>
													<tr>
														<th class="px-4 py-3 font-medium">Unit</th>
														<th class="px-4 py-3 font-medium">Context</th>
														<th class="px-4 py-3 font-medium">Status</th>
														<th class="px-4 py-3 font-medium">Due</th>
														<th class="px-4 py-3 font-medium">Outstanding</th>
													</tr>
												</thead>
												<tbody class="divide-y divide-amber-200">
													{#each overdueInvoices as invoice (invoice.id)}
														<tr>
															<td class="px-4 py-3 align-top">
																<p class="font-medium text-[var(--color-app-text)]">
																	Unit {invoice.unit_number}
																</p>
																{#if invoice.household_name}
																	<p class="text-[var(--color-app-muted)]">
																		{invoice.household_name}
																	</p>
																{/if}
															</td>
															<td class="px-4 py-3 align-top text-[var(--color-app-muted)]">
																<p>
																	{invoice.description ||
																		`Maintenance invoice for ${invoice.billing_period}`}
																</p>
																<p
																	class="mt-1 text-xs tracking-[0.16em] text-[var(--color-app-subtle)] uppercase"
																>
																	{invoice.billing_period} · {invoice.source}
																</p>
															</td>
															<td class="px-4 py-3 align-top">
																<Badge variant={invoiceStatusVariant(invoice)}
																	>{invoice.status}</Badge
																>
															</td>
															<td
																class="px-4 py-3 align-top font-medium text-[var(--color-app-text)]"
															>
																{formatDate(invoice.due_date)}
															</td>
															<td class="px-4 py-3 align-top font-medium text-amber-900">
																{formatMoney(invoice.outstanding_amount)}
															</td>
														</tr>
													{/each}
												</tbody>
											</table>
										</div>
									{/if}
								</div>

								<div class="space-y-3">
									<div>
										<p class="text-lg font-semibold text-[var(--color-app-text)]">Current ledger</p>
										<p class="text-sm text-[var(--color-app-muted)]">
											Upcoming, partial, and settled invoices stay visible without stealing urgency
											from overdue work.
										</p>
									</div>
									<div
										class="overflow-x-auto rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)]"
									>
										<table class="min-w-full divide-y divide-[var(--color-app-border)] text-sm">
											<thead
												class="bg-[var(--color-app-panel)] text-left text-[var(--color-app-subtle)]"
											>
												<tr>
													<th class="px-4 py-3 font-medium">Unit</th>
													<th class="px-4 py-3 font-medium">Status</th>
													<th class="px-4 py-3 font-medium">Due</th>
													<th class="px-4 py-3 font-medium">Billing</th>
													<th class="px-4 py-3 font-medium">Amount</th>
													<th class="px-4 py-3 font-medium">Outstanding</th>
												</tr>
											</thead>
											<tbody class="divide-y divide-[var(--color-app-border)]">
												{#each currentInvoices as invoice (invoice.id)}
													<tr>
														<td class="px-4 py-3 align-top">
															<p class="font-medium text-[var(--color-app-text)]">
																Unit {invoice.unit_number}
															</p>
															<p class="text-[var(--color-app-muted)]">
																{invoice.household_name || 'No household name'}
															</p>
														</td>
														<td class="px-4 py-3 align-top">
															<div class="flex flex-wrap gap-2">
																<Badge variant={invoiceStatusVariant(invoice)}
																	>{invoice.status}</Badge
																>
																<Badge variant="outline">{invoice.source}</Badge>
															</div>
														</td>
														<td
															class="px-4 py-3 align-top font-medium text-[var(--color-app-text)]"
														>
															{formatDate(invoice.due_date)}
														</td>
														<td class="px-4 py-3 align-top text-[var(--color-app-muted)]">
															<p>{invoice.billing_period}</p>
															<p class="mt-1 text-xs">
																{invoice.description || 'Maintenance charge'}
															</p>
														</td>
														<td
															class="px-4 py-3 align-top font-medium text-[var(--color-app-text)]"
														>
															{formatMoney(invoice.amount)}
														</td>
														<td
															class="px-4 py-3 align-top font-medium text-[var(--color-app-text)]"
														>
															{formatMoney(invoice.outstanding_amount)}
														</td>
													</tr>
												{/each}
											</tbody>
										</table>
									</div>
								</div>
							</div>
						{/if}
					</CardContent>
				</Card>
			</div>

			<div class="space-y-6">
				<Card>
					<CardHeader>
						<CardTitle>Recurring Plan Setup</CardTitle>
						<CardDescription>
							Set the monthly maintenance fee schedule and keep configuration separate from the
							issued ledger.
						</CardDescription>
					</CardHeader>
					<CardContent>
						<form method="POST" class="grid gap-4 md:grid-cols-2">
							<div class="space-y-2 md:col-span-2">
								<Label for="name">Plan Name</Label>
								<Input
									id="name"
									name="name"
									value={form?.values?.name ?? ''}
									placeholder="Monthly maintenance"
								/>
							</div>
							<div class="space-y-2">
								<Label for="issue_day_of_month">Issue Day</Label>
								<Input
									id="issue_day_of_month"
									name="issue_day_of_month"
									type="number"
									min="1"
									max="28"
									value={form?.values?.issue_day_of_month ?? '1'}
								/>
							</div>
							<div class="space-y-2">
								<Label for="due_day_of_month">Due Day</Label>
								<Input
									id="due_day_of_month"
									name="due_day_of_month"
									type="number"
									min="1"
									max="28"
									value={form?.values?.due_day_of_month ?? '10'}
								/>
							</div>
							<div class="space-y-2">
								<Label for="default_amount">Default Amount</Label>
								<Input
									id="default_amount"
									name="default_amount"
									type="number"
									min="0"
									step="0.01"
									value={form?.values?.default_amount ?? '8500'}
								/>
							</div>
							<div class="space-y-2">
								<Label for="starts_on">Starts On</Label>
								<Input
									id="starts_on"
									name="starts_on"
									type="date"
									value={form?.values?.starts_on ?? ''}
								/>
							</div>
							<div class="space-y-2">
								<Label for="ends_on">Ends On</Label>
								<Input
									id="ends_on"
									name="ends_on"
									type="date"
									value={form?.values?.ends_on ?? ''}
								/>
							</div>
							<div class="space-y-2 md:col-span-2">
								<Label for="description">Description</Label>
								<Input
									id="description"
									name="description"
									value={form?.values?.description ?? ''}
									placeholder="Standard monthly community maintenance fee"
								/>
							</div>
							<div class="space-y-3 md:col-span-2">
								{#if form?.actionKind === 'createPlan' && form?.createPlanError}
									<p
										class="rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-900"
									>
										{form.createPlanError}
									</p>
								{/if}
								{#if form?.actionKind === 'createPlan' && form?.createPlanSuccess}
									<p
										class="rounded-xl border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm text-emerald-900"
									>
										{form.createPlanSuccess}
									</p>
								{/if}
								<Button type="submit" formAction="?/createPlan" class="w-full">
									Create Recurring Plan
								</Button>
							</div>
						</form>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Generation And Plan Snapshot</CardTitle>
						<CardDescription>
							Create invoice records for a billing period and keep the current recurring plans in
							view.
						</CardDescription>
					</CardHeader>
					<CardContent class="space-y-4">
						<form method="POST" class="space-y-4">
							<div class="space-y-2">
								<Label for="plan_id">Recurring Plan</Label>
								<select
									id="plan_id"
									name="plan_id"
									class="flex h-11 w-full rounded-xl border border-(--color-app-border) bg-[var(--color-app-surface-strong)] px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm transition outline-none focus:border-[var(--color-brand-300)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-brand-200)_45%,transparent)]"
								>
									<option value="">Select a plan</option>
									{#each data.plans as plan (plan.id)}
										<option value={plan.id}>{plan.name}</option>
									{/each}
								</select>
							</div>
							<div class="space-y-2">
								<Label for="billing_period">Billing Period</Label>
								<Input id="billing_period" name="billing_period" placeholder="2026-05" />
							</div>
							{#if form?.actionKind === 'generateInvoices' && form?.generateInvoicesError}
								<p
									class="rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-900"
								>
									{form.generateInvoicesError}
								</p>
							{/if}
							{#if form?.actionKind === 'generateInvoices' && form?.generateInvoicesSuccess}
								<p
									class="rounded-xl border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm text-emerald-900"
								>
									{form.generateInvoicesSuccess}
								</p>
							{/if}
							<Button type="submit" formAction="?/generateInvoices" class="w-full">
								Generate Invoices
							</Button>
						</form>

						<div class="space-y-3 border-t border-[var(--color-app-border)] pt-4">
							<p class="text-sm font-medium text-[var(--color-app-text)]">Current plans</p>
							{#if data.plans.length === 0}
								<p class="text-sm text-[var(--color-app-muted)]">
									Create a plan to start scheduled maintenance invoicing.
								</p>
							{:else}
								{#each data.plans as plan (plan.id)}
									<div
										class="rounded-2xl border border-[var(--color-app-border)] bg-[var(--color-app-surface-strong)] px-4 py-3"
									>
										<div class="flex items-center justify-between gap-3">
											<div>
												<p class="font-medium text-[var(--color-app-text)]">{plan.name}</p>
												<p class="text-sm text-[var(--color-app-muted)]">
													Issue day {plan.issue_day_of_month} · Due day {plan.due_day_of_month}
												</p>
											</div>
											<Badge variant={planStatusVariant(plan)}>{plan.status}</Badge>
										</div>
										<p class="mt-2 text-sm text-[var(--color-app-muted)]">
											{formatMoney(plan.default_amount)} default · {plan.override_count} unit override(s)
										</p>
									</div>
								{/each}
							{/if}
						</div>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Create Manual Invoice</CardTitle>
						<CardDescription>
							Add ad hoc charges without waiting on the recurring schedule.
						</CardDescription>
					</CardHeader>
					<CardContent>
						<form method="POST" class="grid gap-4">
							<div class="space-y-2">
								<Label for="unit_id">Unit</Label>
								<select
									id="unit_id"
									name="unit_id"
									class="flex h-11 w-full rounded-xl border border-(--color-app-border) bg-[var(--color-app-surface-strong)] px-3 py-2 text-sm text-[var(--color-app-text)] shadow-sm transition outline-none focus:border-[var(--color-brand-300)] focus:ring-2 focus:ring-[color-mix(in_oklab,var(--color-brand-200)_45%,transparent)]"
								>
									<option value="">Select a unit</option>
									{#each data.units as unit (unit.id)}
										<option value={unit.id} selected={form?.values?.unit_id === unit.id}>
											{unit.unit_number}
										</option>
									{/each}
								</select>
							</div>
							<div class="space-y-2">
								<Label for="amount">Amount</Label>
								<Input
									id="amount"
									name="amount"
									type="number"
									min="0"
									step="0.01"
									value={form?.values?.amount ?? ''}
								/>
							</div>
							<div class="space-y-2">
								<Label for="due_date">Due Date</Label>
								<Input
									id="due_date"
									name="due_date"
									type="date"
									value={form?.values?.due_date ?? ''}
								/>
							</div>
							<div class="space-y-2">
								<Label for="billing_period">Billing Period</Label>
								<Input
									id="billing_period"
									name="billing_period"
									placeholder="2026-05"
									value={form?.values?.billing_period ?? ''}
								/>
							</div>
							<div class="space-y-2">
								<Label for="manual_description">Description</Label>
								<Input
									id="manual_description"
									name="description"
									placeholder="Special levy or one-time maintenance charge"
									value={form?.values?.description ?? ''}
								/>
							</div>
							{#if form?.actionKind === 'createInvoice' && form?.createInvoiceError}
								<p
									class="rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-900"
								>
									{form.createInvoiceError}
								</p>
							{/if}
							{#if form?.actionKind === 'createInvoice' && form?.createInvoiceSuccess}
								<p
									class="rounded-xl border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm text-emerald-900"
								>
									{form.createInvoiceSuccess}
								</p>
							{/if}
							<Button type="submit" formAction="?/createInvoice" class="w-full">
								Create Manual Invoice
							</Button>
						</form>
					</CardContent>
				</Card>
			</div>
		</div>
	{:else if data.adminError}
		<div class="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900">
			{data.adminError}
		</div>
	{/if}
</div>
