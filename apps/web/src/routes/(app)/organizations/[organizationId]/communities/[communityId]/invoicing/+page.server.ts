import { ApiError } from '$lib/api/http';
import { apiServerRequest } from '$lib/server/api';
import type {
	Invoice,
	InvoicePlan,
	OverdueHouseholdSummary,
	ResidentInvoiceOverview,
	Unit
} from '$lib/api/types';
import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

type OptionalResult<T> = {
	data: T | null;
	error: string | null;
};

async function optionalRequest<T>(
	event: Parameters<PageServerLoad>[0],
	path: string,
	ignoredStatuses: number[]
): Promise<OptionalResult<T>> {
	try {
		return {
			data: await apiServerRequest<T>(event, path),
			error: null
		};
	} catch (error) {
		if (error instanceof ApiError && ignoredStatuses.includes(error.status)) {
			return { data: null, error: null };
		}

		return {
			data: null,
			error: error instanceof ApiError ? error.message : 'Unable to load this invoicing view.'
		};
	}
}

export const load: PageServerLoad = async (event) => {
	const communityBase = `/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}`;

	const [plansResult, invoicesResult, unitsResult, overdueResult, residentOverviewResult] =
		await Promise.all([
			optionalRequest<InvoicePlan[]>(event, `${communityBase}/invoicing/plans`, [401, 403]),
			optionalRequest<Invoice[]>(event, `${communityBase}/invoicing/invoices`, [401, 403]),
			optionalRequest<Unit[]>(event, `${communityBase}/units`, [401, 403]),
			optionalRequest<OverdueHouseholdSummary[]>(
				event,
				`${communityBase}/invoicing/households/overdue`,
				[401, 403]
			),
			optionalRequest<ResidentInvoiceOverview>(
				event,
				`${communityBase}/invoicing/me`,
				[401, 403, 404]
			)
		]);

	const adminAccess = Boolean(plansResult.data || invoicesResult.data || overdueResult.data);

	return {
		plans: plansResult.data ?? [],
		invoices: invoicesResult.data ?? [],
		units: unitsResult.data ?? [],
		overdueHouseholds: overdueResult.data ?? [],
		residentOverview: residentOverviewResult.data,
		adminAccess,
		adminError:
			plansResult.error || invoicesResult.error || unitsResult.error || overdueResult.error || null,
		residentError: residentOverviewResult.error
	};
};

export const actions: Actions = {
	createPlan: async (event) => {
		const form = await event.request.formData();
		const name = String(form.get('name') ?? '').trim();
		const issue_day_of_month = Number(form.get('issue_day_of_month') ?? 0);
		const due_day_of_month = Number(form.get('due_day_of_month') ?? 0);
		const default_amount = Number(form.get('default_amount') ?? 0);
		const starts_on = String(form.get('starts_on') ?? '').trim();
		const ends_on = String(form.get('ends_on') ?? '').trim();
		const description = String(form.get('description') ?? '').trim();

		const values = {
			name,
			issue_day_of_month: String(form.get('issue_day_of_month') ?? ''),
			due_day_of_month: String(form.get('due_day_of_month') ?? ''),
			default_amount: String(form.get('default_amount') ?? ''),
			starts_on,
			ends_on,
			description
		};

		if (!name || !starts_on || issue_day_of_month < 1 || due_day_of_month < 1) {
			return fail(400, {
				actionKind: 'createPlan',
				createPlanError: 'Name, issue day, due day, and start date are required.',
				values
			});
		}

		try {
			await apiServerRequest<InvoicePlan>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/invoicing/plans`,
				{
					method: 'POST',
					body: JSON.stringify({
						name,
						issue_day_of_month,
						due_day_of_month,
						default_amount,
						starts_on,
						ends_on,
						description
					})
				}
			);

			return {
				actionKind: 'createPlan',
				createPlanSuccess: `Created invoice plan ${name}.`
			};
		} catch (error) {
			return fail(error instanceof ApiError ? error.status : 500, {
				actionKind: 'createPlan',
				createPlanError:
					error instanceof ApiError ? error.message : 'Unable to create invoice plan.',
				values
			});
		}
	},
	generateInvoices: async (event) => {
		const form = await event.request.formData();
		const plan_id = String(form.get('plan_id') ?? '').trim();
		const billing_period = String(form.get('billing_period') ?? '').trim();

		if (!plan_id || !billing_period) {
			return fail(400, {
				actionKind: 'generateInvoices',
				generateInvoicesError: 'Select a plan and billing period.'
			});
		}

		try {
			await apiServerRequest<{ plan_id: string; billing_period: string; created_count: number }>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/invoicing/plans/${plan_id}/generate`,
				{
					method: 'POST',
					body: JSON.stringify({ billing_period })
				}
			);

			return {
				actionKind: 'generateInvoices',
				generateInvoicesSuccess: `Generated invoices for ${billing_period}.`
			};
		} catch (error) {
			return fail(error instanceof ApiError ? error.status : 500, {
				actionKind: 'generateInvoices',
				generateInvoicesError:
					error instanceof ApiError ? error.message : 'Unable to generate invoices.'
			});
		}
	},
	createInvoice: async (event) => {
		const form = await event.request.formData();
		const unit_id = String(form.get('unit_id') ?? '').trim();
		const amount = Number(form.get('amount') ?? 0);
		const due_date = String(form.get('due_date') ?? '').trim();
		const billing_period = String(form.get('billing_period') ?? '').trim();
		const description = String(form.get('description') ?? '').trim();

		const values = {
			unit_id,
			amount: String(form.get('amount') ?? ''),
			due_date,
			billing_period,
			description
		};

		if (!unit_id || !due_date || !billing_period || amount <= 0) {
			return fail(400, {
				actionKind: 'createInvoice',
				createInvoiceError: 'Unit, amount, due date, and billing period are required.',
				values
			});
		}

		try {
			await apiServerRequest<Invoice>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/invoicing/invoices`,
				{
					method: 'POST',
					body: JSON.stringify({
						unit_id,
						amount,
						due_date,
						billing_period,
						description
					})
				}
			);

			return {
				actionKind: 'createInvoice',
				createInvoiceSuccess: `Created a manual invoice for ${billing_period}.`
			};
		} catch (error) {
			return fail(error instanceof ApiError ? error.status : 500, {
				actionKind: 'createInvoice',
				createInvoiceError:
					error instanceof ApiError ? error.message : 'Unable to create manual invoice.',
				values
			});
		}
	}
};
