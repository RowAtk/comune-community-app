import { ApiError } from '$lib/api/http';
import type { Household, Resident, Unit } from '$lib/api/types';
import { apiServerRequest } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		const [unit, households, residents] = await Promise.all([
			apiServerRequest<Unit>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/units/${event.params.unitId}`
			),
			apiServerRequest<Household[]>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/households`
			),
			apiServerRequest<Resident[]>(
				event,
				`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/residents`
			)
		]);

		return {
			unit,
			households: households.filter((household) => household.unit_id === event.params.unitId),
			residents: residents.filter((resident) => resident.unit_id === event.params.unitId),
			apiError: null
		};
	} catch (error) {
		return {
			unit: null,
			households: [],
			residents: [],
			apiError: error instanceof ApiError ? error.message : 'Unable to load unit details.'
		};
	}
};
