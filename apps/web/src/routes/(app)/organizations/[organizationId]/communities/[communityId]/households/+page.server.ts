import { ApiError } from '$lib/api/http';
import type { Household, Unit } from '$lib/api/types';
import { apiServerRequest } from '$lib/server/api';
import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		const [households, units] = await Promise.all([
			apiServerRequest<Household[]>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/households`
			),
			apiServerRequest<Unit[]>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/units`
			)
		]);

		return { households, units, apiError: null };
	} catch (error) {
		return {
			households: [],
			units: [],
			apiError: error instanceof ApiError ? error.message : 'Unable to load households.'
		};
	}
};

export const actions: Actions = {
	create: async (event) => {
		const form = await event.request.formData();
		const unit_id = String(form.get('unit_id') ?? '').trim();
		const name = String(form.get('name') ?? '').trim();

		if (!unit_id) {
			return fail(400, {
				createError: 'A unit is required.',
				values: { unit_id, name }
			});
		}

		try {
			await apiServerRequest<Household>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/households`,
				{
					method: 'POST',
					body: JSON.stringify({
						unit_id,
						name
					})
				}
			);

			return {
				createSuccess: name
					? `Household ${name} created successfully.`
					: 'Household created successfully.'
			};
		} catch (error) {
			return fail(error instanceof ApiError ? error.status : 500, {
				createError: error instanceof ApiError ? error.message : 'Unable to create household.',
				values: { unit_id, name }
			});
		}
	}
};
