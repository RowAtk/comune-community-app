import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	const token = event.url.searchParams.get('token')?.trim() ?? '';
	if (!token) {
		throw redirect(303, '/');
	}

	throw redirect(303, `/invites/residents/${encodeURIComponent(token)}`);
};
