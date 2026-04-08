import { apiRequest } from '$lib/api/http';
import type { Unit } from '$lib/api/types';

export type CreateUnitPayload = {
	unit_number: string;
	block_floor: string;
	unit_type: string;
	status: string;
};

export function listUnits(
	fetcher: typeof fetch,
	organizationId: string,
	communityId: string
): Promise<Unit[]> {
	return apiRequest(fetcher, `/v1/organizations/${organizationId}/communities/${communityId}/units`);
}

export function createUnit(
	fetcher: typeof fetch,
	organizationId: string,
	communityId: string,
	payload: CreateUnitPayload
): Promise<Unit> {
	return apiRequest(fetcher, `/v1/organizations/${organizationId}/communities/${communityId}/units`, {
		method: 'POST',
		body: JSON.stringify(payload)
	});
}
