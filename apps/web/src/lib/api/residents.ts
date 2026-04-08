import { apiRequest } from '$lib/api/http';
import type { Resident } from '$lib/api/types';

export function listResidents(
	fetcher: typeof fetch,
	organizationId: string,
	communityId: string
): Promise<Resident[]> {
	return apiRequest(
		fetcher,
		`/v1/organizations/${organizationId}/communities/${communityId}/residents`
	);
}
