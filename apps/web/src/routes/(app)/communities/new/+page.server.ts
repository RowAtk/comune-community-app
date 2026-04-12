import type { Community, ApiEnvelope } from '$lib/api/types';
import type { CommunityCreateValues, DashboardPayload } from '$lib/api/auth';
import { apiServerRequest, apiServerResponse } from '$lib/server/api';
import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

type CreateCommunityPayload = {
	name: string;
	slug: string;
	address: string;
	timezone: string;
	status: string;
	settings: Record<string, never>;
};

const defaultValues = (): CommunityCreateValues => ({
	organization_id: '',
	name: '',
	slug: '',
	address: '',
	timezone: 'America/Jamaica'
});

export const load: PageServerLoad = async (event) => {
	const dashboard = await apiServerRequest<DashboardPayload>(event, '/v1/auth/dashboard');

	return {
		dashboard,
		defaultValues: {
			...defaultValues(),
			organization_id: dashboard.organizations[0]?.organization_id ?? ''
		}
	};
};

export const actions: Actions = {
	default: async (event) => {
		const form = await event.request.formData();
		const values: CommunityCreateValues = {
			organization_id: String(form.get('organization_id') ?? '').trim(),
			name: String(form.get('name') ?? '').trim(),
			slug: String(form.get('slug') ?? '').trim(),
			address: String(form.get('address') ?? '').trim(),
			timezone: String(form.get('timezone') ?? 'America/Jamaica').trim() || 'America/Jamaica'
		};

		if (!values.organization_id || !values.name || !values.slug) {
			const dashboard = await apiServerRequest<DashboardPayload>(event, '/v1/auth/dashboard');
			return fail(400, {
				error: 'Organization, community name, and slug are required.',
				values,
				dashboard
			});
		}

		const payload: CreateCommunityPayload = {
			name: values.name,
			slug: values.slug,
			address: values.address,
			timezone: values.timezone,
			status: 'ACTIVE',
			settings: {}
		};

		const response = await apiServerResponse(
			event,
			`/v1/organizations/${values.organization_id}/communities`,
			{
				method: 'POST',
				body: JSON.stringify(payload)
			}
		);

		if (!response.ok) {
			const errorPayload = (await response.json().catch(() => null)) as
				| ApiEnvelope<unknown>
				| null;
			const dashboard = await apiServerRequest<DashboardPayload>(event, '/v1/auth/dashboard');

			return fail(response.status, {
				error: errorPayload?.error?.message ?? 'Unable to create community.',
				values,
				dashboard
			});
		}

		const payloadEnvelope = (await response.json()) as ApiEnvelope<Community>;
		const community = payloadEnvelope.data;

		throw redirect(303, `/organizations/${community.organization_id}/communities/${community.id}`);
	}
};
