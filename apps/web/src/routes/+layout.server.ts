import type { AuthPayload } from '$lib/api/auth';
import { ApiError } from '$lib/api/http';
import { apiServerRequest } from '$lib/server/api';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async (event) => {
	try {
		const auth = await apiServerRequest<AuthPayload>(event, '/v1/auth/me');
		return { auth };
	} catch (error) {
		if (error instanceof ApiError && error.status === 401) {
			return { auth: null };
		}

		throw error;
	}
};
