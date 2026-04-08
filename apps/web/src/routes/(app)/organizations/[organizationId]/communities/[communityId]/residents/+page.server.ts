import { ApiError } from '$lib/api/http';
import type { Resident } from '$lib/api/types';
import { apiServerRequest } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		const residents = await apiServerRequest<Resident[]>(
			event,
			`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}/residents`
		);
		return { residents, apiError: null };
	} catch (error) {
		return {
			residents: [],
			apiError: error instanceof ApiError ? error.message : 'Unable to load residents.'
		};
	}
};
