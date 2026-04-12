import { ApiError } from '$lib/api/http';
import { apiServerRequest, apiServerResponse } from '$lib/server/api';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import type { OrganizationCreateValues } from '$lib/api/auth';
import type { DashboardPayload } from '$lib/api/auth';
import type { ApiEnvelope } from '$lib/api/types';

type CreateOrganizationPayload = OrganizationCreateValues & {
	status: string;
	settings: Record<string, never>;
};

const defaultValues = (): OrganizationCreateValues => ({
	name: '',
	slug: '',
	legal_name: '',
	billing_email: '',
	phone: '',
	country_code: 'JM',
	timezone: 'America/Jamaica'
});

export const load: PageServerLoad = async (event) => {
	const dashboard = await apiServerRequest<DashboardPayload>(event, '/v1/auth/dashboard');

	return {
		dashboard,
		defaultValues: defaultValues()
	};
};

export const actions: Actions = {
	default: async (event) => {
		const form = await event.request.formData();
		const values: OrganizationCreateValues = {
			name: String(form.get('name') ?? '').trim(),
			slug: String(form.get('slug') ?? '').trim(),
			legal_name: String(form.get('legal_name') ?? '').trim(),
			billing_email: String(form.get('billing_email') ?? '').trim(),
			phone: String(form.get('phone') ?? '').trim(),
			country_code: String(form.get('country_code') ?? 'JM').trim() || 'JM',
			timezone: String(form.get('timezone') ?? 'America/Jamaica').trim() || 'America/Jamaica'
		};

		if (!values.name || !values.slug) {
			return fail(400, {
				error: 'Organization name and slug are required.',
				values
			});
		}

		const payload: CreateOrganizationPayload = {
			...values,
			status: 'ACTIVE',
			settings: {}
		};

		const response = await apiServerResponse(event, '/v1/organizations', {
			method: 'POST',
			body: JSON.stringify(payload)
		});

		if (!response.ok) {
			const errorPayload = (await response.json().catch(() => null)) as
				| ApiEnvelope<unknown>
				| null;

			return fail(response.status, {
				error: errorPayload?.error?.message ?? 'Unable to create organization.',
				values
			});
		}

		throw redirect(303, '/dashboard');
	}
};
