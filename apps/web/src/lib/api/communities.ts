import { apiRequest } from '$lib/api/http';
import type { Community } from '$lib/api/types';

export function getCommunity(
	fetcher: typeof fetch,
	organizationId: string,
	communityId: string
): Promise<Community> {
	return apiRequest(fetcher, `/v1/organizations/${organizationId}/communities/${communityId}`);
}
