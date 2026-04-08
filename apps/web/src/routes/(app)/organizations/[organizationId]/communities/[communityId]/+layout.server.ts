import type { Community } from '$lib/api/types';
import { ApiError } from '$lib/api/http';
import { apiServerRequest } from '$lib/server/api';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async (event) => {
	try {
		const community = await apiServerRequest<Community>(
			event,
			`/v1/organizations/${event.params.organizationId}/communities/${event.params.communityId}`
		);
		return { community, apiError: null, params: event.params };
	} catch (error) {
		return {
			community: null,
			apiError: error instanceof ApiError ? error.message : 'Unable to load community details.',
			params: event.params
		};
	}
};
