import { ApiError } from '$lib/api/http';
import type { Household, Unit } from '$lib/api/types';
import { apiServerRequest } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		const household = await apiServerRequest<Household>(
			event,
			`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/households/${event.params.householdId}`
		);

		let unit: Unit | null = null;

		try {
			unit = await apiServerRequest<Unit>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/units/${household.unit_id}`
			);
		} catch {
			unit = null;
		}

		return { household, unit, apiError: null };
	} catch (error) {
		return {
			household: null,
			unit: null,
			apiError: error instanceof ApiError ? error.message : 'Unable to load household details.'
		};
	}
};
