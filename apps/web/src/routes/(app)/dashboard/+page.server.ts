import type { DashboardPayload } from '$lib/api/auth';
import { apiServerRequest } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	const dashboard = await apiServerRequest<DashboardPayload>(event, '/v1/auth/dashboard');

	return { dashboard };
};
