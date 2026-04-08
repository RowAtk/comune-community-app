import { ApiError } from '$lib/api/http';
import type { Household } from '$lib/api/types';
import { apiServerRequest } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		const households = await apiServerRequest<Household[]>(
			event,
			`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/households`
		);
		return { households, apiError: null };
	} catch (error) {
		return {
			households: [],
			apiError: error instanceof ApiError ? error.message : 'Unable to load households.'
		};
	}
};
