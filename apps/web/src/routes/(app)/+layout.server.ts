import { apiServerRequest } from '$lib/server/api';
import { ApiError } from '$lib/api/http';
import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import type { AuthPayload } from '$lib/api/auth';

export const load: LayoutServerLoad = async (event) => {
	try {
		const auth = await apiServerRequest<AuthPayload>(event, '/v1/auth/me');
		return { auth };
	} catch (error) {
		if (error instanceof ApiError && error.status === 401) {
			throw redirect(303, '/');
		}

		throw error;
	}
};
