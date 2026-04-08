import { apiRequest } from '$lib/api/http';
import type { Household } from '$lib/api/types';

export function listHouseholds(
	fetcher: typeof fetch,
	organizationId: string,
	communityId: string
): Promise<Household[]> {
	return apiRequest(
		fetcher,
		`/v1/organizations/${organizationId}/communities/${communityId}/households`
	);
}
