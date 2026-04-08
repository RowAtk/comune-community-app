import { ApiError } from '$lib/api/http';
import type { Unit } from '$lib/api/types';
import { apiServerRequest } from '$lib/server/api';
import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		const units = await apiServerRequest<Unit[]>(
			event,
			`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/units`
		);
		return { units, apiError: null };
	} catch (error) {
		return {
			units: [],
			apiError: error instanceof ApiError ? error.message : 'Unable to load units.'
		};
	}
};

export const actions: Actions = {
	create: async (event) => {
		const form = await event.request.formData();
		const unit_number = String(form.get('unit_number') ?? '').trim();
		const block_floor = String(form.get('block_floor') ?? '').trim();
		const unit_type = String(form.get('unit_type') ?? '').trim();
		const status = String(form.get('status') ?? 'ACTIVE').trim().toUpperCase();

		if (!unit_number) {
			return fail(400, {
				createError: 'Unit number is required.',
				values: { unit_number, block_floor, unit_type, status }
			});
		}

		try {
			await apiServerRequest<Unit>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/units`,
				{
					method: 'POST',
					body: JSON.stringify({
						unit_number,
						block_floor,
						unit_type,
						status
					})
				}
			);

			return {
				createSuccess: `Unit ${unit_number} created successfully.`
			};
		} catch (error) {
			return fail(error instanceof ApiError ? error.status : 500, {
				createError: error instanceof ApiError ? error.message : 'Unable to create unit.',
				values: { unit_number, block_floor, unit_type, status }
			});
		}
	}
};
